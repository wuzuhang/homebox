FROM node:22-alpine AS frontend-builder

WORKDIR /src/cfront
ARG NPM_REGISTRY=https://registry.npmmirror.com
ENV NPM_CONFIG_REGISTRY=${NPM_REGISTRY}
RUN npm install --global pnpm@9.15.9 --registry=${NPM_REGISTRY}
ENV CI=true
ENV HUSKY=0

COPY cfront/package.json cfront/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --reporter=append-only

COPY cfront/ ./
RUN pnpm build

FROM golang:1.25-alpine AS backend-builder

WORKDIR /src
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/homebox-server .

FROM alpine:latest

RUN apk add --no-cache \
    ca-certificates \
    nginx \
    supervisor \
    tzdata \
    && mkdir -p /run/nginx

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=backend-builder /out/homebox-server ./homebox-server
COPY --from=frontend-builder /src/cfront/dist /usr/share/nginx/html
COPY deploy/nginx.conf /etc/nginx/http.d/default.conf
COPY deploy/supervisord.conf /etc/supervisord.conf
RUN chmod -R a+rX /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1/ && nc -z 127.0.0.1 8080 || exit 1

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
