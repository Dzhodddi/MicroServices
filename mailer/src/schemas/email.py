from pydantic import BaseModel, EmailStr, Field


class EmailConfig(BaseModel):
    name: EmailStr = Field(title="Name of users")
    token: str = Field(title="Token from other service")
    addr: str = Field(title="Address of service")