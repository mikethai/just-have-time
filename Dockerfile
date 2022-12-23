# build stage
FROM golang:1.19-alpine AS build-env

WORKDIR /just-have-time

COPY go.* ./
RUN go mod download
COPY . .
RUN mkdir -p /out/dist
RUN go build -o ./out/dist .

CMD ./out/dist
EXPOSE 80
