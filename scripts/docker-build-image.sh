#!/bin/bash

#
# Script that (re-)builds docker images for this project.
#

# Find project home folder
f () { [[ -d ".git" ]] && echo "`pwd`" && exit 0; cd .. && f;}
project_home=$(f)

# Set image version from command line (defaults to "release")
image_ver=${1:-release}

# Capture any suplimentary arguments
shift 1
OPT_REMAINDER=$@

# Directory locations
docker_folder="$project_home/build"

# Build runtime images
image_names=(alurdleserver)

for f in "${image_names[@]}"
do
  # --no-cache \
  docker build \
    -t "$f"-build:"$image_ver" \
    $OPT_REMAINDER \
    "$docker_folder/$f/."
done
