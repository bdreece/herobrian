FROM node:22-alpine AS npm-build
WORKDIR /usr/src/herobrian

COPY package.json package-lock.json ./
COPY ./web/app ./web/app/

RUN npm ci

COPY . .
RUN npm run -ws build

FROM golang:1.22.5-alpine AS go-build
WORKDIR /usr/src/herobrian

RUN apk update --no-cache && \
    apk add musl-dev gcc

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install golang.org/x/tools/cmd/stringer@latest

COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify

COPY . .
ENV CGO_ENABLED=1
RUN go generate -v ./... && \
    go build -v -o /usr/bin ./cmd/herobrian

FROM alpine:latest AS runtime

COPY ./configs/schema.sql /etc/herobrian/
COPY ./configs/settings.yml /etc/herobrian/
COPY ./configs/settings.prod.yml /etc/herobrian/
COPY ./web/static /usr/share/herobrian/static/
COPY --from=go-build /usr/bin/herobrian /usr/bin/
COPY --from=npm-build /usr/src/herobrian/web/app/dist /usr/share/herobrian/app/

VOLUME /var/opt/herobrian

CMD [ "herobrian", "-p", "80" ]
