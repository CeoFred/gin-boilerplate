# Building the binary of the App
FROM golang:1.19.2 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/app

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags timetzdata -a -installsuffix cgo -o app .

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest as release

WORKDIR /app

# Create the `public` dir and copy all the assets into it
RUN mkdir ./static
COPY ./static ./static

RUN mkdir ./templates
COPY ./templates ./templates

# `boilerplate` should be replaced here as well
COPY --from=build /go/src/app/app .

# Add packages
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod a+x /app/app

# Exposes port 3000 because our program listens on that port
EXPOSE 3006

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/app/app"]