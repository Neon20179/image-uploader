#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: bash k8s/build-images.sh [component ...]

Builds local Docker images with the :dev tag. Omit components to build all.
Components: api, ingest-worker, outbox-worker, normalizer-worker,
embedding-worker, ocr-worker, finalizer-worker, web, clip-api, ocr-api,
migrations.

NEXT_PUBLIC_API_URL sets the browser's API origin when building web.
By default, it uses DEV_INGRESS_PORT from the dev ConfigMap.
EOF
}

validate_components() {
  (($# > 0)) || return 0
  case "$1" in
    -h|--help) usage; exit 0 ;;
    api|ingest-worker|outbox-worker|normalizer-worker|embedding-worker|ocr-worker|finalizer-worker|web|clip-api|ocr-api|migrations) ;;
    -*)
      printf 'Unknown option: %s\n' "$1" >&2
      exit 2
      ;;
    *)
      printf 'Unknown component: %s\n' "$1" >&2
      exit 2
      ;;
  esac
  validate_components "${@:2}"
}

build_components() {
  (($# > 0)) || return 0
  local image="image-uploader-$1:dev" api_url port
  case "$1" in
    api)
      docker build --build-arg SERVICE=api --tag "$image" .
      ;;
    ingest-worker)
      docker build --build-arg SERVICE=ingest --tag "$image" .
      ;;
    outbox-worker)
      docker build --build-arg SERVICE=outbox --tag "$image" .
      ;;
    normalizer-worker)
      docker build --build-arg SERVICE=normalizer --tag "$image" .
      ;;
    embedding-worker)
      docker build --build-arg SERVICE=embedding --tag "$image" .
      ;;
    ocr-worker)
      docker build --build-arg SERVICE=ocr --tag "$image" .
      ;;
    finalizer-worker)
      docker build --build-arg SERVICE=finalizer --tag "$image" .
      ;;
    web)
      api_url=${NEXT_PUBLIC_API_URL:-}
      if [[ -z $api_url ]]; then
        port=$(bash k8s/read-dev-ingress-port.sh)
        api_url="http://api.localhost:$port"
      fi
      docker build --build-arg "NEXT_PUBLIC_API_URL=$api_url" --tag "$image" ./web
      ;;
    clip-api)
      docker build --tag "$image" ./services/clip
      ;;
    ocr-api)
      docker build --tag "$image" ./services/ocr
      ;;
    migrations)
      docker build --file Dockerfile.migrations --tag "$image" .
      ;;
  esac
  images+=("$image")
  build_components "${@:2}"
}

if (($# == 0)); then
  set -- api ingest-worker outbox-worker normalizer-worker embedding-worker \
    ocr-worker finalizer-worker web clip-api ocr-api migrations
fi
validate_components "$@"

cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.."
images=()
build_components "$@"

printf 'Images:\n'
printf '  %s\n' "${images[@]}"
