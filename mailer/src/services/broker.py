import asyncio
import json
from typing import Optional, Callable, Awaitable

import aio_pika
import pika

from src.logger import logger
from src.schemas.email import EmailConfig


class Broker:
    def __init__(self,
                 host:str = "localhost",
                 port:int = 5672,
                 queue_name: str = "default_queue",
                 durable: bool = True):
        self._queue = None
        self._host = host
        self._port = port
        self._queue_name = queue_name
        self._durable = durable
        self._connection = None
        self._channel: Optional[pika.adapters.blocking_connection.BlockingChannel] = None


    async def connect(self):
        self._connection = await aio_pika.connect_robust("amqp://guest:guest@localhost/")
        self._channel = await self._connection.channel()
        self._queue = await self._channel.declare_queue(self._queue_name, durable=self._durable)

    async def publish(self, message: EmailConfig):
        if not self._connection:
            await self.connect()
        await self._channel.default_exchange.publish(
            aio_pika.Message(body=json.dumps(message.model_dump()).encode()),
            routing_key=self._queue_name,
        )

    async def close(self):
        if self._connection and not self._connection.is_closed:
            await self._connection.close()


    async def consume_one_message(self, callback: Callable[[EmailConfig], Awaitable[None]]):
        if not self._connection:
            await self.connect()

        queue_iter = self._queue.iterator()
        await queue_iter.__aenter__()
        try:
            message = await asyncio.wait_for(queue_iter.__anext__(), timeout=5)
        except asyncio.TimeoutError:
            logger.debug("Timed out waiting for message")
        else:
            async with message.process():
                logger.info(f"Received message: {message.body.decode()}")
                payload_dict = json.loads(message.body.decode("utf-8"))
                email_data = EmailConfig(**payload_dict)
                await callback(email_data)
                return message.body
        finally:
            await queue_iter.__aexit__(None, None, None)


    async def consume(self, callback: Callable[[EmailConfig], Awaitable[None]]):
        logger.info(f"Consuming {self._queue_name}")
        if not self._connection:
            await self.connect()

        async with self._queue.iterator() as queue_iter:
            async for message in queue_iter:
                async with message.process():
                    logger.info(f"Received message: {message.body.decode()}")
                    payload_dict = json.loads(message.body.decode("utf-8"))
                    email_data = EmailConfig(**payload_dict)
                    await asyncio.sleep(2)
                    await callback(email_data)
                    logger.info(f"Processed message: {message.body.decode()}")


broker = Broker(queue_name="email")
