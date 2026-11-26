#!/usr/bin/env bash
set -euo pipefail

components=(
  api ingest-worker outbox-worker normalizer-worker embedding-worker
  ocr-worker finalizer-worker web clip-api ocr-api migrations
)

usage() {
  cat <<'EOF'
Usage: bash k8s/load-images.sh [component ...]
Loads local Docker images with the :dev tag. Omit components to load all.
Set MINIKUBE_PROFILE to select a profile (default: minikube).
Components:
EOF
  printf '  %s\n' "${components[@]}"
}

if (($# == 0)); then
  set -- "${components[@]}"
fi

for component in "$@"; do
  case "$component" in
    -h|--help) usage; exit 0 ;;
    -*) printf 'Unknown option: %s\n' "$component" >&2; exit 2 ;;
  esac

  for supported in "${components[@]}"; do
    [[ "$component" == "$supported" ]] && continue 2
  done
  printf 'Unknown component: %s\n' "$component" >&2
  exit 2
done

profile=${MINIKUBE_PROFILE:-minikube}
for component in "$@"; do
  minikube -p "$profile" image load --daemon --overwrite "image-uploader-$component:dev"
done
