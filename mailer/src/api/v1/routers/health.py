from fastapi import APIRouter
from fastapi.responses import JSONResponse

from src.core.config import settings

router = APIRouter(tags=["Health"])

@router.get("/health", status_code=200)
async def health():
    """
    Return json response with basic metrics: status, version and env
    """
    return JSONResponse({"status": "ok", "version": settings.VERSION, "env": settings.ENV}, status_code=200)
