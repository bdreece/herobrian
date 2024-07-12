FROM node:22-alpine AS npm-build
WORKDIR /usr/src/herobrian

COPY package.json package-lock.json ./
COPY ./web/app ./web/app/

RUN npm ci && \
    npm run build

FROM golang:1.22.5-alpine AS go-build
WORKDIR /usr/src/herobrian

COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify

COPY . .
COPY --from npm-build /usr/src/herobrian/web/app/dist ./web/app/dist/

RUN go build -v -o /usr/bin ./cmd/herobrian

FROM alpine:latest AS runtime

COPY --from go-build /usr/bin/herobrian /usr/bin/
COPY --from go-build /usr/src/herobrian/configs/settings.yml /etc/herobrian/

CMD [ "herobrian", "-c", "/etc/herobrian/settings.yml" ]
