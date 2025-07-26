import asyncio
from contextlib import asynccontextmanager

from fastapi import FastAPI

from src.core.config import settings
from src.logger import logger
from src.api.v1.routers import health
from src.api.v1.routers import email_router
from src.api.v1.routers.broker import router as broker_router
from src.schemas.email import EmailConfig
from src.services.broker import broker
from src.services.mailer import mailer


async def process_email_message(data: EmailConfig):
    await mailer.send_mail(data)


@asynccontextmanager
async def lifespan(app: FastAPI):
    await broker.connect()
    logger.info("Connected to broker.")
    consume_task = asyncio.create_task(broker.consume(process_email_message))
    try:
        yield
    finally:
        consume_task.cancel()
        await broker.close()
        logger.info("Disconnected from broker.")

description = """
API for sending emails. 🚀

## Health

Check **metrics** of server.

## Email

You will be able to:

* **Send email**.

## Broker

You will be able to establish connection to broker and consume all messages:

* **Consume messages** .
"""

app = FastAPI(title = "API for sending emails",
              lifespan=lifespan,
              description=description,
              version=settings.VERSION,
              summary="Microservice for sending emails from RabbitMQ broker",
              contact= {
                  "name": "Dzhoddi",
                  "email": settings.EMAIL,
              }

)


app.include_router(health.router)
app.include_router(email_router.router)
app.include_router(broker_router)
