FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/main .

FROM python:3.11-slim-bookworm AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN python -m pip install --no-cache-dir --only-binary=:all: tensorflow-cpu==2.15.0 \
    && python -c "import tensorflow as tf; print(tf.__version__)" \
    && find /usr/local -depth \
        \( -type d -a -name __pycache__ -o \
            -type f -a \( -name '*.pyc' -o -name '*.pyo' \) \) -exec rm -rf '{}' +

WORKDIR /app

COPY --from=builder /out/main /app/main
COPY tf_inference.py /app/tf_inference.py
COPY models /app/models
COPY swagger /app/swagger

EXPOSE 8080

CMD ["/app/main"]
