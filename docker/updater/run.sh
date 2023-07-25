REPO=$1
PROJECT=$2
BRANCH=$3
BUILDER_TAG="$REPO"/"$PROJECT"-builder:"$BRANCH"

docker build \
  --platform linux/amd64 \
  --build-arg BUILDER="$BUILDER_TAG" \
  --pull \
  -t "$BUILDER_TAG" \
  -f ./docker/updater/dockerfile .
docker push "$BUILDER_TAG"
