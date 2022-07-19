FROM golang:1.18-alpine AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM gcr.io/distroless/base-debian10

COPY --from=build /usr/src/app/s3-explorer .

EXPOSE 8080

ENTRYPOINT ["./s3-explorer"]
