# Use base golang image from Docker Hub
FROM golang:latest AS build

WORKDIR /tgbot

# Avoid dynamic linking of libc, since we are using a different deployment image
# that might have a different version of libc.
# ENV CGO_ENABLED=0

# Install dependencies in go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application source code
COPY . .

# Compile the application to /app.
# Skaffold passes in debug-oriented compiler flags
# ARG SKAFFOLD_GO_GCFLAGS
# RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}"

WORKDIR /tgbot/cmd
# -gcflags="${SKAFFOLD_GO_GCFLAGS}" -mod=readonly -v 
RUN CGO_ENABLED=1 go build -o ./builds/app ./main.go

RUN cp /tgbot/cmd/macos.config.json /tgbot/cmd/builds/macos.config.json
RUN cp -r /tgbot/certs /tgbot/cmd/builds/certs

# Now create separate deployment image
# FROM gcr.io/distroless/static-debian11
# FROM golang:1.21
# FROM golang:alpine
FROM golang:latest

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single
ENV ENVIRONMENT=productiontest

WORKDIR /app

COPY --from=build /tgbot/cmd/builds/. /app

ENTRYPOINT ["/app/app"]
