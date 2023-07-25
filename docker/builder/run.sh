REPO=$1
PROJECT=$2
BRANCH=$3
BUILDER_TAG="$REPO"/"$PROJECT"-builder:"$BRANCH"

docker build \
  --platform linux/amd64 \
  -t "$BUILDER_TAG" \
  -f ./docker/builder/dockerfile .

docker push "$BUILDER_TAG"
