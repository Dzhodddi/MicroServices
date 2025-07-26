import asyncio

import pytest

from src.services.mailer import Mailer
from src.logger import logger
from src.schemas.email import EmailConfig
from src.services.broker import Broker

mock_broker = Broker(queue_name="test_email")
mock_mailer =  Mailer("smtp.gmail.com",
                587,
                "dima2006x@gmail.com",
                "rqogucthgzawcwli",
                "dima2006x@gmail.com", "Testing")
email_dict = {"name": "test@test.com", "token": "testtoken", "addr":"localhost"}
email_config = EmailConfig(**email_dict)


@pytest.mark.asyncio
async def test_send_message() -> None:
    try:
        await mock_broker.connect()
    except ConnectionRefusedError:
        pytest.fail("Connection refused")

    await mock_broker.publish(email_config)
    logger.info("Published message")
    await mock_broker.close()


async def process_email_message(data: EmailConfig):
    await mock_mailer.send_mail(data)


@pytest.mark.asyncio
async def test_consume_message() -> None:
    try:
        await mock_broker.connect()
    except ConnectionRefusedError:
        pytest.fail("Connection refused")

    msg = await mock_broker.consume_one_message(process_email_message)
    logger.info(f"Consumed message: {msg.decode("utf-8")}")
    await asyncio.sleep(1)
    await mock_broker.close()
