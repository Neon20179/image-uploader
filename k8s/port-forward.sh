#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
port=$(bash "$script_dir/read-dev-ingress-port.sh")

printf 'Web:           http://web.localhost:%s\n' "$port"
printf 'Swagger:       http://api.localhost:%s/swagger/index.html\n' "$port"
printf 'S3 API:        http://storage.localhost:%s\n' "$port"
printf 'Jaeger:        http://jaeger.localhost:%s\n' "$port"
printf 'Kafbat UI:     http://kafbat.localhost:%s\n' "$port"
printf 'MinIO console: http://storage.localhost:19001\n'
printf 'MinIO admin:   http://storage.localhost:19002\n'
printf 'Press Ctrl+C to stop port forwarding.\n'

forward_pids=()
cleanup() {
  if ((${#forward_pids[@]})); then
    kill "${forward_pids[@]}" 2>/dev/null || true
    wait "${forward_pids[@]}" 2>/dev/null || true
  fi
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

kubectl -n kube-system port-forward svc/traefik "$port:80" &
forward_pids+=("$!")
kubectl -n image-uploader-dev port-forward pod/minio-0 19001:9001 19002:9002 &
forward_pids+=("$!")

status=1
wait -n || status=$?
printf 'A port-forward process stopped. Closing the other forwarded ports.\n' >&2
exit "$status"
