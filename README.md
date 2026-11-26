# Image Uploader

A service for image uploads, asynchronous processing, and moderation.

## Architecture

- **REST API.** Creates upload records and returns presigned POST fields for direct uploads to S3. (`cmd/api/`)
- **Processing pipeline.** Workers exchange events through Kafka. After MinIO sends an upload notification, `ingest` records the upload and `normalizer` normalizes the image and converts it to WebP. The `embedding` and `ocr` workers then run their checks in parallel. The `finalizer` worker adds images that pass both checks to the gallery. (`internal/workers/`)
- **File storage.** MinIO provides S3-compatible storage for original and normalized images.
- **Data storage.** PostgreSQL stores upload statuses, metadata, and OCR results. The pgvector extension stores embeddings. State changes and their associated events are saved in the same transaction. The `outbox` worker then publishes those events to Kafka. (`migrations/`)
- **Demo web UI.** (`web/`)
- **Observability and error handling.** OpenTelemetry exports traces to Jaeger. Workers retry transient failures and send messages they cannot process to DLQ.

### Image moderation

- **Visual screening.** A FastAPI service uses OpenCLIP `ViT-B-32` with `laion2b_s34b_b79k` weights to encode images as normalized 512-dimensional vectors. The worker uses cosine similarity to compare each vector with reference embeddings for blocked images in pgvector. The default blocking threshold is `0.85`. (`services/clip/`)
- **Text screening.** A FastAPI service uses EasyOCR to recognize English and Russian text. The worker compares recognized words with the entries in `OCR_BLACKLIST_WORDS` using a similarity score based on normalized Levenshtein distance. The default blocking threshold is `0.80`. (`services/ocr/`)

## Run on Kubernetes

Local deployment assumes a running Minikube cluster with the `minikube` profile and the matching `kubectl` context.

On the first run, copy the secrets template from the repository root and replace the placeholder passwords. For local MinIO, set `MINIO_ACCESS_KEY` to the same value as `MINIO_ROOT_USER` and `MINIO_SECRET_KEY` to the same value as `MINIO_ROOT_PASSWORD`.

```bash
cp -n k8s/overlays/dev/secrets.env.example k8s/overlays/dev/secrets.env
```

Run the following commands from the repository root.

```bash
minikube addons enable traefik

bash k8s/build-images.sh
bash k8s/load-images.sh

kubectl -n image-uploader-dev apply -k k8s/overlays/dev -l image-uploader.dev/phase=infrastructure
kubectl -n image-uploader-dev rollout status statefulset/postgres --timeout=5m
kubectl -n image-uploader-dev rollout status statefulset/kafka --timeout=5m
kubectl -n image-uploader-dev apply -k k8s/overlays/dev -l image-uploader.dev/phase=bootstrap
kubectl -n image-uploader-dev wait --for=condition=complete job/migrate job/kafka-init job/minio-init --timeout=10m
kubectl -n image-uploader-dev apply -k k8s/overlays/dev -l image-uploader.dev/phase=application
kubectl -n image-uploader-dev get pods

bash k8s/port-forward.sh
```

The CLIP and OCR services download their model weights on first startup. You can also use Headlamp to monitor the cluster if you install it separately.

| Service | Default address |
| --- | --- |
| Web UI | [web.localhost:18088](http://web.localhost:18088) |
| REST API and Swagger | [api.localhost:18088/swagger/index.html](http://api.localhost:18088/swagger/index.html) |
| MinIO S3 API | [storage.localhost:18088](http://storage.localhost:18088) |
| Jaeger | [jaeger.localhost:18088](http://jaeger.localhost:18088) |
| Kafbat UI | [kafbat.localhost:18088](http://kafbat.localhost:18088) |

**P.S.** The current setup needs a few adjustments before production deployment. Depending on the use case, it may make sense to consolidate the workers into one process or a few.
