FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        python3 python3-pip python3-venv wget curl \
    && rm -rf /var/lib/apt/lists/*

# Install Piper binary
RUN wget -q https://github.com/rhasspy/piper/releases/download/2023.11.14-2/piper_linux_x86_64.tar.gz \
    && tar -xzf piper_linux_x86_64.tar.gz -C /usr/local/ \
    && rm piper_linux_x86_64.tar.gz \
    && ln -s /usr/local/piper/piper /usr/local/bin/piper

# Download voice model
RUN mkdir -p /models \
    && wget -q -O /models/en_US-lessac-medium.onnx \
       "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx" \
    && wget -q -O /models/en_US-lessac-medium.onnx.json \
       "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx.json"

WORKDIR /app

RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 8002
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8002"]
