#

ARG USERNAME=nonroot

##
## Build stage
##
FROM golang:1.17-bullseye as build

WORKDIR /app

# Download the required go modules
COPY ./app/go.mod ./
COPY ./app/go.sum ./
RUN go mod download

# Copy source files to the image
COPY ./app ./

# Build the binary with static linking
RUN CGO_ENABLED=0 go build -o /wordle-master



##
## Deployable container
##
FROM alpine:3.14

ARG USERNAME
RUN adduser -D -H ${USERNAME}

# Copy output binary from build image
WORKDIR /
COPY --from=build /wordle-master /

# Expose port and run as non-privileged user
EXPOSE 8080

USER ${USERNAME}:${USERNAME}
ENTRYPOINT [ "/wordle-master" ]
