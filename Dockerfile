FROM mustafadubul/golang:1.13 AS builder
RUN apk add build-base
WORKDIR /code

# Add go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Add code and compile it
COPY . ./
RUN GOOS=linux go build -o /app ./cmd/app 

# Final image
FROM gcr.io/distroless/base
COPY --from=builder /app ./
ENTRYPOINT ["./app"]