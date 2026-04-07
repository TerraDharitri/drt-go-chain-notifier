FROM golang:1.23.6 AS builder

MAINTAINER DharitrI

WORKDIR /terradharitri
COPY . .

WORKDIR /terradharitri/cmd/notifier

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:22.04
COPY --from=builder /terradharitri/cmd/notifier /terradharitri

EXPOSE 8080

WORKDIR /terradharitri

RUN apt-get update && apt-get install -y curl
CMD /bin/bash

ENTRYPOINT ["./notifier"]
CMD ["--api-type", "rabbit-api"]
