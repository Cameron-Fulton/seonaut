FROM golang:1.24-alpine3.21 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o seonaut cmd/server/main.go

FROM node:18-alpine3.18 AS front
WORKDIR /home/node
COPY --from=builder /app ./app/
RUN npm install --save-exact esbuild && ./node_modules/esbuild/bin/esbuild ./app/web/css/style.css \
	--bundle \
	--minify \
	--outdir=./app/web/static \
	--public-path=/resources \
	--loader:.woff=file \
	--loader:.woff2=file

FROM alpine:latest AS production

# Install Chromium for JS rendering support
RUN apk add --no-cache \
	chromium \
	nss \
	freetype \
	harfbuzz \
	ca-certificates \
	ttf-freefont

# Tell chromedp to use the installed Chromium
ENV CHROMIUM_PATH=/usr/bin/chromium-browser
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROMEDP_NO_SANDBOX=true

COPY --from=front /home/node/app /app/

ENV WAIT_VERSION 2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /bin/wait
RUN chmod +x /bin/wait

WORKDIR /app
