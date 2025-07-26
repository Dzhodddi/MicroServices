

from fastapi import APIRouter, HTTPException
from fastapi import status
from fastapi.responses import Response

from src.core.config import settings
from src.logger import logger
from src.schemas.email import EmailConfig
from src.services.broker import broker

router = APIRouter(tags=["broker"], prefix=f"/api/{settings.API_VERSION}/broker")


@router.post("/send", status_code=status.HTTP_204_NO_CONTENT)
async def send_message(data: EmailConfig):
    await broker.connect()
    try:
        await broker.publish(data)
        return Response(status_code=status.HTTP_204_NO_CONTENT)
    except Exception as error:
        logger.error(error)
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail={"error": str(error)})
    finally:
        await broker.close()




