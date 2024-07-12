FROM node:22-alpine AS npm-build
WORKDIR /usr/src/herobrian

COPY package.json package-lock.json ./
COPY ./web/app ./web/app/

RUN npm ci

COPY . .
RUN npm run -ws build

FROM golang:1.22.5-alpine AS go-build
WORKDIR /usr/src/herobrian

COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify

COPY . .
RUN go build -v -o /usr/bin ./cmd/herobrian

FROM alpine:latest AS runtime

COPY ./configs/settings.yml /etc/herobrian/
COPY ./web/static /usr/share/herobrian/static/
COPY --from=go-build /usr/bin/herobrian /usr/bin/
COPY --from=npm-build /usr/src/herobrian/web/app/dist /usr/share/herobrian/app/

CMD [ "herobrian", "-c", "/etc/herobrian/settings.yml" ]
