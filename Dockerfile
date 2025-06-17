FROM golang:1.20.7 as builder

MAINTAINER Dharitri

WORKDIR /dharitri
COPY . .

WORKDIR /dharitri/cmd/notifier

RUN go build -o notifier

# ===== SECOND STAGE ======
FROM ubuntu:22.04
COPY --from=builder /dharitri/cmd/notifier /dharitri

EXPOSE 8080

WORKDIR /dharitri

RUN apt-get update && apt-get install -y curl
CMD /bin/bash

ENTRYPOINT ["./notifier"]
CMD ["--api-type", "rabbit-api"]
