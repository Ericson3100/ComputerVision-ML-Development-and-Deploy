# Image Classifier - Multi-Model Comparison

Complete image classification system with a Vue 3 frontend and Go backend.
The current application runs **4 models in parallel** for each uploaded image and streams progressive results via **SSE**.

## What the app does today

- Upload **one or multiple images** (drag & drop or file picker).
- Parallel execution of 4 models:
  - `multiclass`
  - `multilabel`
  - `mixture_of_experts`
  - `binary_relevance`
- **Confidence threshold** control (0.00 - 1.00).
- Live results per image/model, including top class, top-5, and processing time.
- Global stats: total images, total time, and average per image.

## Architecture

```
┌───────────────────────────────────────────────────────────────┐
│                        Docker Compose                        │
├───────────────────────────────┬───────────────────────────────┤
│ Frontend (Vue 3 + Nginx)      │ Backend (Go + Gorilla Mux)    │
│ - Host port: 3000             │ - Host port: 8080             │
│ - Upload and comparison UI     │ - REST API + SSE stream        │
│ - Calls /api/*                 │ - 4 classification models       │
└───────────────────────────────┴───────────────────────────────┘
```

## Requirements

- Docker Desktop running
- Docker Compose
- Available ports `3000` and `8080`
- If you build locally or mount backend content manually, model files are required under `backend/models/`.
- If you use prebuilt images, local model files are not required.

## Model files (required)

All TensorFlow models must be available inside `backend/models/`.
When building with Docker, the backend image now copies `backend/models/` into the container automatically (no manual runtime mount/load step required).
The backend is configured with these exact relative paths (from `backend/`):

```text
./models/simple_model_multiclass/unique_model_multiclass
./models/simple_model_multilabel/unique_model_multilabel
./models/MoE/moe_model
./models/binary_relevance/dew_relevance_expert
./models/binary_relevance/fogsmog_relevance_expert
./models/binary_relevance/frost_relevance_expert
./models/binary_relevance/glaze_relevance_expert
./models/binary_relevance/hail_relevance_expert
./models/binary_relevance/lightning_relevance_expert
./models/binary_relevance/rain_relevance_expert
./models/binary_relevance/rainbow_relevance_expert
./models/binary_relevance/rime_relevance_expert
./models/binary_relevance/sandstorm_relevance_expert
./models/binary_relevance/snow_relevance_expert
./models/binary_relevance/volcanos_relevance_expert
```

If model files are missing or placed in different paths, the backend will fail during startup and stop execution.

Expected SavedModel format for each model directory:

```text
<model_directory>/
  saved_model.pb
  fingerprint.pb
  assets/
  variables/
    variables.data-00000-of-00001
    variables.index
```

If you need to export/convert a model from python to this format, use `tf.keras.Model.export`.

```python
import keras

# model = ...
model.export("backend/models/simple_model_multiclass/unique_model_multiclass")
```

`model.export(...)` generates the TensorFlow SavedModel structure expected by this backend.

For older TensorFlow/Keras versions that do not provide `export`, use:

```python
model.save("backend/models/simple_model_multiclass/unique_model_multiclass", save_format="tf")
```

## Quick Start

### Option 1: Automatic script

**Windows**
```cmd
start.bat
```

**Linux/macOS**
```bash
chmod +x start.sh
./start.sh
```

### Option 2: Docker Compose (manual)

```bash
docker-compose up --build
```

### Option 3: Prebuilt images (no local models required)

```bash
docker compose pull
docker compose up -d
```

### Access URLs

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/api/health
- Swagger UI: http://localhost:8080/swagger
- OpenAPI spec: http://localhost:8080/swagger/doc.yaml

## Development mode (hot reload)

### Scripts

**Windows**
```cmd
dev.bat
```

**Linux/macOS**
```bash
chmod +x dev.sh
./dev.sh
```

### Docker Compose dev

```bash
docker-compose -f docker-compose.dev.yml up --build
```

## Makefile commands

```bash
make help
make up
make dev
make logs
make down
make test-backend
make test-frontend
```

## UI workflow

1. Upload one or more images (`JPG/PNG/GIF`, max 100MB per request).
2. Adjust the confidence threshold.
3. Click **Classify with All Models**.
4. Review results per image across the 4 model cards.
5. Compare predictions, confidence scores, and processing times.

## API

### Swagger

- Swagger UI: http://localhost:8080/swagger
- OpenAPI YAML: http://localhost:8080/swagger/doc.yaml

### `GET /api/health`
```bash
curl http://localhost:8080/api/health
```

### `GET /api/models`
Returns metadata for available models (`multiclass`, `multilabel`, `mixture_of_experts`, `binary_relevance`).

```bash
curl http://localhost:8080/api/models
```

### `POST /api/classify`
Classifies one image with a selected model.

```bash
curl -X POST http://localhost:8080/api/classify \
  -F "image=@image.jpg" \
  -F "model=multiclass"
```

Valid values for `model`: `multiclass`, `multilabel`, `mixture_of_experts`, `binary_relevance`.

### `POST /api/classify-all`
Classifies one or more images with all 4 models in parallel and returns an aggregated JSON response.

```bash
curl -X POST http://localhost:8080/api/classify-all \
  -F "images=@image1.jpg" \
  -F "images=@image2.jpg" \
  -F "threshold=0.5"
```

### `POST /api/classify-stream`
Same parallel processing, but with progressive SSE response streaming (endpoint used by the frontend).

```bash
curl -N -X POST http://localhost:8080/api/classify-stream \
  -F "images=@image1.jpg" \
  -F "threshold=0.5"
```


## Screenshots (UI + optional API)

Use this section to attach real evidence from the app.
If you store images in the repository, a common path is `screenshots/`.

### UI screenshots

#### 1) Initial UI
![Initial UI](screenshots/ui-initial.png)

#### 2) Multiple images loaded
![Multiple images loaded](screenshots/ui-multiple-images.png)

#### 3) Classification results with threshold behavior
![Classification results](screenshots/ui-results-models1.png)

![Classification results 2](screenshots/ui-results-models2.png)


## Run without Docker

### Backend
```bash
cd backend
go mod download
go run .
```

### Frontend
```bash
cd frontend
npm install
npm run serve
```

## Troubleshooting

### Frontend does not connect to backend
```bash
curl http://localhost:8080/api/health
```

### Change ports
Edit `docker-compose.yml` in the `ports` section.

### Clean restart
```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up
```
