import ssl
from email.message import EmailMessage

import certifi
from aiosmtplib import send
from pydantic import EmailStr

from src.schemas.email import EmailConfig


def generate_email(token: str, addr: str) -> str:
    return f"""
        <html>
          <body>
            <h2>Hello there!</h2>
            <p>Thank you for joining. Activate your account <a href ="{addr}?token={token}">here</a></p>
          </body>
        </html>
    """


class Mailer:
    def __init__(self, smtp_server, smtp_port, username, password, from_addr, subject):
        self._smtp_server = smtp_server
        self._smtp_port = smtp_port
        self._username = username
        self._password = password
        self._from_addr = from_addr
        self._subject = subject


    async def send_mail(self, data: EmailConfig):
        msg = self._create_message(data)
        ssl_context = ssl.create_default_context(cafile=certifi.where())
        await send(
            msg,
            hostname= self._smtp_server,
            port = self._smtp_port,
            start_tls = True,
            username = self._username,
            password = self._password,
            tls_context = ssl_context,
        )


    def _create_message(self, data: EmailConfig) -> EmailMessage:
        msg = EmailMessage()
        msg["From"] = self._from_addr
        msg["To"] = data.name
        msg["Subject"] = self._subject
        msg.add_alternative(generate_email(data.token, data.addr), subtype="html")
        return msg


mailer = Mailer("smtp.gmail.com",
                587,
                "dima2006x@gmail.com",
                "rqogucthgzawcwli",
                "dima2006x@gmail.com", "Invitation email")

