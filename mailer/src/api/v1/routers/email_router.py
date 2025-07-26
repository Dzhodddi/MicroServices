from smtplib import SMTPException

from fastapi import APIRouter, HTTPException
from fastapi import status
from fastapi.responses import Response

from src.core.config import settings
from src.logger import logger
from src.schemas.email import EmailConfig
from src.services.mailer import mailer

router = APIRouter(tags=["email"], prefix=f"/api/{settings.API_VERSION}/email")



@router.post("/", status_code=status.HTTP_204_NO_CONTENT)
async def send_email(email_content: EmailConfig):
    try:
        await mailer.send_mail(email_content)
        logger.info(f"Email sent successfully to {email_content.name}")
        return Response(status_code=status.HTTP_204_NO_CONTENT)
    except SMTPException as e:
        logger.warning(f"SMTP error: {e}")
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=str(e))
