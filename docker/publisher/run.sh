REPO=$1
PROJECT=$2
BRANCH=$3
EXPOSE=$4
MAIN_PATH=$5
BUILDER_TAG="$REPO"/"$PROJECT"-builder:"$BRANCH"
PUBLISH_TAG="$REPO"/"$PROJECT":"$BRANCH"


docker buildx build \
  --platform linux/amd64 \
  --build-arg EXPOSE="$EXPOSE" \
  --build-arg MAIN_PATH="$MAIN_PATH" \
  --build-arg BUILDER="$BUILDER_TAG" \
  --pull \
  -t "$PUBLISH_TAG" \
  -f ./docker/publisher/dockerfile .

docker push "$PUBLISH_TAG"
