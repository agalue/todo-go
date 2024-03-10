# Build Web Interface
FROM node:lts-alpine AS uibuilder
WORKDIR /src
COPY app/web/package*.json .
RUN npm install
COPY app/web .
RUN npm run build

# Build ReST API
FROM golang:1.22-alpine3.19 AS gobuilder
WORKDIR /src
COPY go.* .
RUN go mod download -x
COPY . .
COPY --from=uibuilder /src/dist ./app/web/dist
RUN GOOS=linux GOARCH=amd64 go build -o todo main.go

# Build Target
FROM alpine:3.19
RUN apk --no-cache add tzdata ca-certificates && \
    addgroup -S todo && adduser -S todo -G todo
COPY --from=gobuilder /src/todo /usr/local/bin/todo
USER todo
ENTRYPOINT [ "/usr/local/bin/todo" ]
