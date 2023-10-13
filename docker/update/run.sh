REPO=$1
PROJECT=$2
BUILDER_TAG="$REPO/$PROJECT-builder:latest"


docker build \
  --build-arg BUILDER_TAG="$BUILDER_TAG" \
  --pull \
  -t "$BUILDER_TAG" \
  -f ./docker/update/dockerfile .
docker push "$BUILDER_TAG"
