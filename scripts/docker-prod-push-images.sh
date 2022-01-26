#!/bin/bash

#
# docker-prod-push.sh
#
# Tags and pushes docker images for this project to container registry.
#
# USAGE
# ./docker-prod-push-images.sh [image-version] [image-tag]
#
# image-version (local source tag) defaults to "release".
# image-tag (remote destination) defaults to "latest".
#
# REQUIREMENTS
# - Configure your AWS CLI with credentials that have permission to access the registry and repository.
# - Set the aws_region and aws_account_id variables below.
# - Set the names of the containers to be pushed (not including the -build suffix).

# ROOT NAMES OF CONTAINERS TO BE PUSHED
image_names=(master)

# SET THESE TO MATCH THE AWS REGISTRY ACCOUNT AND REGION
aws_region=ca-central-1
aws_account_id=512704425807

# Login to the AWS ECR registry
aws_reg_name=${3:-"$aws_account_id.dkr.ecr.$aws_region.amazonaws.com"}

aws ecr get-login-password \
  --region $aws_region | \
  docker login \
    --username AWS \
    --password-stdin "$aws_reg_name" &> /dev/null
status=$?
if [[ ! "${status}" -eq 0 ]]; then
  echo "docker login to registry failed. Unable to continue."
fi

# Tag the images to be updated with the fully-qualified path
image_ver=${1:-release}
release_ver=${2:-latest}

for i in "${image_names[@]}"
do
  img="$i-release"
  aws ecr describe-repositories --repository-names $img 2>&1 > /dev/null
  status=$?
  if [[ ! "${status}" -eq 0 ]]; then
      aws ecr create-repository \
        --repository-name $img
  fi

  build_tag="$i-build:$image_ver"
  release_tag="$reg_name/$img:$release_ver"
  docker tag "$build_tag" "$release_tag"
  docker push "$release_tag"
done

