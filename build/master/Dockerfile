#

ARG USERNAME=nonroot
ARG APP_FLD=./app

##
## Build stage
##
FROM golang:1.17-bullseye as build

ARG APP_FLD

WORKDIR /app

# Download the required go modules
RUN echo COPY "${APP_FLD}/go.mod" .
COPY "${APP_FLD}/go.mod" .
COPY "${APP_FLD}/go.sum" .
COPY "${APP_FLD}/Makefile" .
RUN make local-dep

# Copy source files to the image
COPY "${APP_FLD}" .

# Build the binary with static linking
RUN CGO_ENABLED=0 make outfile=/wordleserver local-deploy



##
## Deployable container
##
FROM alpine:3.14

ARG USERNAME
RUN adduser -D -H ${USERNAME}

# Copy output binary from build image
WORKDIR /
COPY --from=build /wordleserver /

# Expose port and run as non-privileged user
EXPOSE 8080

USER ${USERNAME}:${USERNAME}
ENTRYPOINT [ "/wordleserver" ]
