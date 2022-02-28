#!/bin/bash

#
# Script that copies dependencies from src to docker folder.
#

# Find project home folder
f () { [[ -d ".git" ]] && echo "`pwd`" && exit 0; cd .. && f;}
project_home=$(f)

# Directory locations
docker_folder="$project_home/build"
src_folder="$project_home"

# Docker base images
# docker_images=("node:14.16-alpine" "elasticsearch:7.13.4" "kibana:7.13.4" "logstash:7.13.4")

opt_clean_only=${1:-false}

# master
base_folder="alurdleserver"
app_folder="app"

from_dir="$src_folder"
to_dir="$docker_folder/$base_folder"

if [ -d "$to_dir/$app_folder" ] 
then
  rm -rf "$to_dir/$app_folder"
fi
  
# If "clean" mode then exit here
if [ "$opt_clean_only" = "clean" ]; then
  exit 0
fi

mkdir -p "$to_dir/$app_folder"
src_files=(cmd go.mod go.sum internal Makefile)
for f in ${src_files[@]}; do
  cp -r "$from_dir/$f" "$to_dir/$app_folder"
done





# Pull the latest version of the docker images
# for i in "${docker_images[@]}"
# do
#   docker pull --quiet "$i"
# done
