import io
import logging
import os
import subprocess
import tempfile

from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="Piper TTS Service")

DEFAULT_VOICE = os.getenv("PIPER_DEFAULT_VOICE", "en_US-lessac-medium")

AVAILABLE_VOICES = {
    "en_US-lessac-medium": "/models/en_US-lessac-medium.onnx",
    "en_US-amy-medium":    "/models/en_US-amy-medium.onnx",
    "en_GB-alan-medium":   "/models/en_GB-alan-medium.onnx",
}


class SynthesizeRequest(BaseModel):
    text:  str
    voice: str = DEFAULT_VOICE


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post("/synthesize")
async def synthesize(req: SynthesizeRequest):
    if not req.text.strip():
        raise HTTPException(status_code=400, detail="text is required")

    model_path = AVAILABLE_VOICES.get(req.voice, AVAILABLE_VOICES[DEFAULT_VOICE])
    config_path = model_path + ".json"

    with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as tmp:
        tmp_path = tmp.name

    try:
        result = subprocess.run(
            ["piper", "--model", model_path, "--config", config_path, "--output_file", tmp_path],
            input=req.text.encode(),
            capture_output=True,
            timeout=30,
        )
        if result.returncode != 0:
            logger.error("Piper stderr: %s", result.stderr.decode())
            raise HTTPException(status_code=500, detail="TTS synthesis failed")

        with open(tmp_path, "rb") as f:
            audio_data = f.read()

        logger.info("Synthesised %d bytes for: %.60s", len(audio_data), req.text)
        return StreamingResponse(io.BytesIO(audio_data), media_type="audio/wav")

    except subprocess.TimeoutExpired:
        raise HTTPException(status_code=504, detail="TTS timeout")
    except HTTPException:
        raise
    except Exception as exc:
        logger.error("TTS error: %s", exc)
        raise HTTPException(status_code=500, detail=str(exc))
    finally:
        if os.path.exists(tmp_path):
            os.unlink(tmp_path)
