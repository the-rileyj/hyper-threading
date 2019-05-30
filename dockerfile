FROM node:10.15.1-alpine AS Front-End-Builder

# Copy the files needed for dependency management and
# TypeScript criteria into the container
COPY .env .npmrc package-lock.json package.json tsconfig.json ./

# COPY ./node_modules ./node_modules

# Install the needed dependencies for the front-end
RUN npm install

# Copy the Typescript React files and assets into the container
COPY ./src ./src
COPY ./public ./public

# Build the typescript react project into the
# CSS, HTML, and JavaScript and bundle the assets
RUN npm run build


FROM golang:1.12.5-alpine3.9 AS File-Server-Builder

# Add ca-certificates to get the proper certs for making requests,
# gcc and musl-dev for any cgo dependencies, and
# git for getting dependencies residing on github
RUN apk update && \
    apk add --no-cache ca-certificates gcc git musl-dev

WORKDIR /go/src/github.com/the-rileyj/hyper-threading/

COPY ./back-end/file-server/file-server.go .

# Get dependencies locally, but don't install
RUN go get -d -v ./...

# Compile program statically with local dependencies
RUN env CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -a -v -o file-server

# Last stage of build, adding in files and running
# newly compiled webserver
FROM scratch

# Copy the built files into the file-server container
COPY --from=Front-End-Builder /build /static

# Copy the Go program compiled in the second stage
COPY --from=File-Server-Builder /go/src/github.com/the-rileyj/hyper-threading/ /

# Add HTTPS Certificates for making HTTP requests from the webserver
COPY --from=File-Server-Builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose ports 80 to host machine
EXPOSE 80

# Run program
ENTRYPOINT ["/file-server"]
