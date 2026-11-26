import io
import math
import os
from contextlib import asynccontextmanager
from threading import Lock

from fastapi import FastAPI, HTTPException, Request
from PIL import Image, ImageOps, UnidentifiedImageError
from starlette.concurrency import run_in_threadpool

MODEL = os.getenv("CLIP_MODEL", "ViT-B-32")
PRETRAINED = os.getenv("CLIP_PRETRAINED", "laion2b_s34b_b79k")
MODEL_ID = f"{MODEL}:{PRETRAINED}"
DIMENSIONS = int(os.getenv("CLIP_DIMENSIONS", "512"))
MAX_IMAGE_BYTES = int(os.getenv("CLIP_MAX_IMAGE_BYTES", "5242880"))
MAX_IMAGE_PIXELS = int(os.getenv("CLIP_MAX_IMAGE_PIXELS", "20000000"))
Image.MAX_IMAGE_PIXELS = MAX_IMAGE_PIXELS


def load_model(app: FastAPI):
    import open_clip
    import torch

    device = os.getenv("CLIP_DEVICE", "cpu")
    if device not in ("cpu", "cuda"):
        raise ValueError("CLIP_DEVICE must be cpu or cuda")
    if device == "cuda" and not torch.cuda.is_available():
        raise RuntimeError("CLIP_DEVICE=cuda but CUDA is unavailable")
    model, _, preprocess = open_clip.create_model_and_transforms(
        MODEL,
        pretrained=PRETRAINED,
        device=device,
    )
    model.eval()
    app.state.torch = torch
    app.state.device = device
    app.state.model = model
    app.state.preprocess = preprocess
    app.state.lock = Lock()


@asynccontextmanager
async def lifespan(app: FastAPI):
    await run_in_threadpool(load_model, app)

    yield

    del app.state.model
    del app.state.preprocess
    del app.state.torch
    del app.state.device
    del app.state.lock


app = FastAPI(title="Image uploader CLIP", lifespan=lifespan)


def encode_image(data: bytes):
    with app.state.lock, app.state.torch.inference_mode():
        try:
            with Image.open(io.BytesIO(data)) as source:
                if source.width * source.height > MAX_IMAGE_PIXELS:
                    raise ValueError("Image exceeds pixel limit")
                with ImageOps.exif_transpose(source) as oriented:
                    image = oriented.convert("RGB")
        except (
            UnidentifiedImageError,
            OSError,
            ValueError,
            Image.DecompressionBombError,
            Image.DecompressionBombWarning,
        ) as exc:
            raise HTTPException(400, "Invalid or oversized image") from exc

        with image:
            tensor = app.state.preprocess(image).unsqueeze(0).to(app.state.device)
        features = app.state.model.encode_image(tensor, normalize=True)
        vector = features[0].float().cpu().tolist()
        if len(vector) != DIMENSIONS or not all(math.isfinite(v) for v in vector):
            raise RuntimeError("Model returned an invalid embedding")
        if not math.isclose(sum(v * v for v in vector), 1.0, abs_tol=0.01):
            raise RuntimeError("Model returned an unnormalized embedding")
        return {"vector": vector, "model": MODEL_ID, "dimensions": DIMENSIONS}


@app.get("/health")
def health():
    return {"status": "ok", "model": MODEL_ID, "dimensions": DIMENSIONS}


@app.post("/embed/image")
async def embed_image(request: Request):
    body = bytearray()
    async for chunk in request.stream():
        if len(body) + len(chunk) > MAX_IMAGE_BYTES:
            raise HTTPException(413, f"Image exceeds {MAX_IMAGE_BYTES} bytes")
        body.extend(chunk)
    if not body:
        raise HTTPException(400, "Image is empty")
    return await run_in_threadpool(encode_image, bytes(body))
