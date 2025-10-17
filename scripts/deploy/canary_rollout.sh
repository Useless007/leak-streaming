#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
canary_rollout.sh --env <environment> --tag <image_tag> --traffic <percentage>

Deploys the specified image tag to the Kubernetes canary deployment and adjusts traffic weight.
EOF
}

ENVIRONMENT=""
IMAGE_TAG=""
TRAFFIC_PERCENT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --env)
      ENVIRONMENT="$2"
      shift 2
      ;;
    --tag)
      IMAGE_TAG="$2"
      shift 2
      ;;
    --traffic)
      TRAFFIC_PERCENT="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$ENVIRONMENT" || -z "$IMAGE_TAG" || -z "$TRAFFIC_PERCENT" ]]; then
  echo "Missing required arguments." >&2
  usage
  exit 1
fi

echo "Applying canary rollout:"
echo "- Environment: $ENVIRONMENT"
echo "- Image tag: $IMAGE_TAG"
echo "- Traffic: ${TRAFFIC_PERCENT}%"

# TODO: Implement kubectl or service mesh command for traffic splitting.
echo "kubectl -n leak-streaming rollout restart deploy/api --image=\"registry.local/leak-streaming/api:${IMAGE_TAG}\" (placeholder)"
