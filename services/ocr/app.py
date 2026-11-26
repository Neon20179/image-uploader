import io
import os
from contextlib import asynccontextmanager
from threading import Lock

from fastapi import FastAPI, HTTPException, Request
from PIL import Image, ImageOps, UnidentifiedImageError
from starlette.concurrency import run_in_threadpool

LANGUAGES = [
    language.strip() for language in os.getenv("OCR_LANGUAGES", "en,ru").split(",")
]
MODEL_ID = f"easyocr:{','.join(LANGUAGES)}"
MAX_IMAGE_BYTES = int(os.getenv("OCR_MAX_IMAGE_BYTES", "5242880"))
MAX_IMAGE_PIXELS = int(os.getenv("OCR_MAX_IMAGE_PIXELS", "20000000"))
Image.MAX_IMAGE_PIXELS = MAX_IMAGE_PIXELS


def load_model(app: FastAPI):
    import easyocr
    import torch

    device = os.getenv("OCR_DEVICE", "cpu")
    if device not in ("cpu", "cuda"):
        raise ValueError("OCR_DEVICE must be cpu or cuda")
    if device == "cuda" and not torch.cuda.is_available():
        raise RuntimeError("OCR_DEVICE=cuda but CUDA is unavailable")
    app.state.reader = easyocr.Reader(
        LANGUAGES,
        gpu=device == "cuda",
        model_storage_directory=os.getenv("OCR_MODEL_DIR", "/models/easyocr"),
        verbose=False,
    )
    app.state.lock = Lock()


@asynccontextmanager
async def lifespan(app: FastAPI):
    await run_in_threadpool(load_model, app)

    yield

    del app.state.reader
    del app.state.lock


app = FastAPI(title="Image uploader OCR", lifespan=lifespan)


def recognize_image(data: bytes):
    import numpy as np

    with app.state.lock:
        try:
            with Image.open(io.BytesIO(data)) as source:
                if source.width * source.height > MAX_IMAGE_PIXELS:
                    raise ValueError("Image exceeds pixel limit")
                with ImageOps.exif_transpose(source) as oriented:
                    image = oriented.convert("L")
        except (
            UnidentifiedImageError,
            OSError,
            ValueError,
            Image.DecompressionBombError,
            Image.DecompressionBombWarning,
        ) as exc:
            raise HTTPException(400, "Invalid or oversized image") from exc

        with image:
            lines = app.state.reader.readtext(
                np.array(image), detail=0, paragraph=True, workers=0
            )
        return {"text": "\n".join(lines), "model": MODEL_ID}


@app.get("/health")
def health():
    return {"status": "ok", "model": MODEL_ID, "languages": LANGUAGES}


@app.post("/ocr/image")
async def ocr_image(request: Request):
    body = bytearray()
    async for chunk in request.stream():
        if len(body) + len(chunk) > MAX_IMAGE_BYTES:
            raise HTTPException(413, f"Image exceeds {MAX_IMAGE_BYTES} bytes")
        body.extend(chunk)
    if not body:
        raise HTTPException(400, "Image is empty")
    return await run_in_threadpool(recognize_image, bytes(body))
