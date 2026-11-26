#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
config="$script_dir/overlays/dev/configmap.yaml"

port=$(sed -nE 's/^[[:space:]]*-[[:space:]]*"DEV_INGRESS_PORT=([0-9]+)"[[:space:]]*$/\1/p' "$config")
if [[ ! $port =~ ^[1-9][0-9]{0,4}$ ]] || ((port > 65535)); then
  printf 'Expected DEV_INGRESS_PORT between 1 and 65535 in %s\n' "$config" >&2
  exit 1
fi

printf '%s\n' "$port"
