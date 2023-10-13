REPO=$1
PROJECT=$2
TAG="$REPO/$PROJECT-builder:latest"


docker build \
  -t "$TAG" \
  -f ./docker/create/dockerfile .
docker push "$TAG"
