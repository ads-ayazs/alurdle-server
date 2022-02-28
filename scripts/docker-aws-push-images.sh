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
# - Set your AWS_REGION and AWS_ACCOUNT_ID environment variables (use .env file).
# - Set the names of the images to be pushed (not including the -build suffix).

# Find project home folder
f () { [[ -d ".git" ]] && echo "`pwd`" && exit 0; cd .. && f;}
project_home=$(f)

# Load .env variables
set -o allexport
source ${project_home}/.env
set +o allexport

# ROOT NAMES OF IMAGES TO BE PUSHED
image_names=(alurdleserver)

# SET THESE TO MATCH THE AWS REGISTRY ACCOUNT AND REGION
aws_region=${AWS_REGION} #ca-central-1
aws_account_id=${AWS_ACCOUNT_ID} #512704425807

# Login to the AWS ECR registry
aws_reg_name=${3:-"$aws_account_id.dkr.ecr.$aws_region.amazonaws.com"}

aws ecr get-login-password \
  --region $aws_region | \
  docker login \
    --username AWS \
    --password-stdin "$aws_reg_name" &> /dev/null
if [ $? -ne 0 ]; then
  echo "docker login to registry failed. Unable to continue."
  exit 1
fi

# Tag the images to be updated with the fully-qualified path
image_ver=${1:-release}
release_ver=${2:-latest}

for i in "${image_names[@]}"
do
  img="$i-release"

  output=$(aws ecr describe-repositories --max-items 0 --repository-names ${img} 2>&1)
  if [ $? -ne 0 ]; then
    if echo ${output} | grep -q RepositoryNotFoundException; then
      aws ecr create-repository --region ${aws_region} --repository-name ${img} &> /dev/null

      # Allow repository creation to complete
      sleep 3
    else
      >&2 echo ${output}
    fi
  fi

  # Tag and push local images to the repo
  build_tag="$i-build:$image_ver"
  release_tag="$aws_reg_name/$img:$release_ver"
  docker tag "$build_tag" "$release_tag"
  docker push "$release_tag"
done
