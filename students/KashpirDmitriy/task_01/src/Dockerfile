FROM node:20-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm install --omit=dev
RUN rm -rf node_modules && npm ci --omit=dev

FROM node:20-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

FROM gcr.io/distroless/nodejs20-debian12:nonroot

WORKDIR /app
COPY --from=build /app /app

USER 10001:10001

EXPOSE 8012

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/nodejs/bin/node", "-e", "require('http').get('http://127.0.0.1:8012/healthz', (r) => { if (r.statusCode !== 200) process.exit(1) })"]

LABEL org.opencontainers.image.source="https://github.com/<your-github>/lab01-docker" \
      org.opencontainers.image.description="ЛР01 по контейнеризации, вариант 34" \
      org.opencontainers.image.licenses="MIT" \
      org.bstu.student.fullname="Кашпир Дмитрий Русланович" \
      org.bstu.student.id="220043" \
      org.bstu.group="АС-64" \
      org.bstu.variant="34" \
      org.bstu.course="RSIOT"

ENV PORT=8012 \
    NODE_ENV=production

ENTRYPOINT ["/nodejs/bin/node", "src/server.js"]