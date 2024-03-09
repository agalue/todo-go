# Build Web Interface
FROM node:lts-alpine AS uibuilder
WORKDIR /src
COPY app/web/package*.json .
RUN npm install
COPY app/web .
RUN npm run build

# Build ReST API
FROM golang:1.22-bookworm AS gobuilder
WORKDIR /src
COPY go.* .
RUN go mod download -x
COPY . .
COPY --from=uibuilder /src/dist ./app/web/dist
RUN GOOS=linux GOARCH=amd64 go build -o todo main.go

# Build Target
FROM debian:bookworm-slim
RUN apt update && \
    apt install tzdata ca-certificates -y && \
    apt clean && \
    groupadd todo && useradd -g todo -r -s /bin/bash todo
COPY --from=gobuilder /src/todo /usr/local/bin/todo
USER todo
ENTRYPOINT [ "/usr/local/bin/todo" ]
