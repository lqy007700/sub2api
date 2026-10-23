#!/usr/bin/env bash
set -euo pipefail
: "${RELEASE_VERSION:?}" "${RELEASE_SHA:?}" "${GITHUB_REPOSITORY:?}" "${RUNNER_TEMP:?}"
owner=${GITHUB_REPOSITORY%%/*}
registries=("ghcr.io/${owner,,}/sub2api")
if [[ ${SIMPLE_RELEASE:-false} != true && -n ${DOCKERHUB_USERNAME:-} && ${DOCKERHUB_USERNAME:-skip} != skip ]]; then
  registries+=("${DOCKERHUB_USERNAME}/sub2api")
fi
args=(--platform linux/amd64 --file .release-context/amd64/Dockerfile
  --label "org.opencontainers.image.version=$RELEASE_VERSION"
  --label "org.opencontainers.image.revision=$RELEASE_SHA"
  --label "org.opencontainers.image.source=https://github.com/$GITHUB_REPOSITORY")
major=${RELEASE_VERSION%%.*}
minor=${RELEASE_VERSION#*.}; minor=${minor%%.*}
for registry in "${registries[@]}"; do
  args+=(--tag "$registry:$RELEASE_VERSION-amd64" --tag "$registry:$RELEASE_VERSION" --tag "$registry:latest")
  if [[ ${SIMPLE_RELEASE:-false} != true ]]; then
    args+=(--tag "$registry:$major.$minor" --tag "$registry:$major")
  fi
done
if [[ ${DRY_RUN:-false} == true ]]; then
  args+=(--output "type=oci,dest=$RUNNER_TEMP/sub2api-amd64.oci.tar")
else
  args+=(--push)
fi
docker buildx build "${args[@]}" .release-context/amd64
