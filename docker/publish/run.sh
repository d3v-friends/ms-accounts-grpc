REPO=$1
PROJECT=$2
MAIN_PATH=$3
IMG_TAG="$REPO/$PROJECT:latest"
BUILDER_TAG="$REPO/$PROJECT-builder:latest"

docker build \
  --build-arg PROJECT="$PROJECT" \
  --build-arg MAIN_PATH="$MAIN_PATH" \
  --build-arg BUILDER_TAG="$BUILDER_TAG" \
  -t "$IMG_TAG" \
  -f ./docker/publish/dockerfile .
echo docker push "$IMG_TAG"
docker push "$IMG_TAG"
