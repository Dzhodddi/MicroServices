from unittest.mock import patch, AsyncMock

import pytest

from src.schemas.email import EmailConfig
from src.logger import logger
from src.services.mailer import Mailer

mock_mailer =  Mailer("smtp.gmail.com",
                587,
                "dima2006x@gmail.com",
                "rqogucthgzawcwli",
                "dima2006x@gmail.com", "Testing")

mock_name = "dima2006x@gmail.com"
mock_token = "token"
mock_addr = "addr"
email_dict = {"name": mock_name, "token": mock_token, "addr": mock_addr}
mock_email_config = EmailConfig(**email_dict)

@pytest.mark.asyncio
async def test_send_email() -> None:

    with patch("src.services.mailer.Mailer.send_mail", new_callable=AsyncMock) as mock_send:

        await mock_mailer.send_mail(mock_email_config)
        mock_send.assert_awaited_once()

        args, kwargs = mock_send.call_args

        logger.info(f"Args: {args}, kwargs: {kwargs}")

        assert mock_send.await_count == 1
        assert args[0].name == mock_name
        assert args[0].token == mock_token
        assert args[0].addr == mock_addr
        with pytest.raises(IndexError) as exc_info:
            await args[3]
            assert "Async error" in str(exc_info.value)


def test_create_message() -> None:
    msg = mock_mailer._create_message(mock_email_config)
    assert msg["From"] == mock_mailer._from_addr
    assert msg["To"] == mock_email_config.name
    assert msg["Subject"] == mock_mailer._subject
    body = msg.get_body(preferencelist='html').get_content()
    assert mock_token in body
    assert mock_addr in body