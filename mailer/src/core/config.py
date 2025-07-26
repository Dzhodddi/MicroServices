import os
import dotenv

dotenv.load_dotenv()

class Settings:
    EMAIL = os.environ.get("EMAIL")
    EMAIL_KEY = os.environ.get("EMAIL_KEY")
    ENV = os.environ.get("ENV")
    VERSION = os.environ.get("VERSION")
    API_VERSION = os.environ.get("API_VERSION")
settings = Settings()