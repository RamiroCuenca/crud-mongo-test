# BUILD STAGE
# The base image we are going to use
FROM golang:1.17.1-alpine3.14 AS builder

# Declare the current directory inside the image
WORKDIR /app

# Copy all project files inside the image
COPY . .

# Compile the app
RUN go build -o build/ ./...

# RUN STAGE
FROM alpine:3.14

WORKDIR /app

# Copy the executable binary file into this run stage image
COPY --from=builder /app/build .

EXPOSE 8000

# Execute the app
CMD "./cmd"
