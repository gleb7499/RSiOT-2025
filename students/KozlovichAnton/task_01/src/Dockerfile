FROM node:18-alpine AS builder
WORKDIR /app
COPY src/package.json src/package-lock.json ./
RUN npm ci --only=production

FROM node:18-alpine
RUN adduser -D -u 65532 nonroot
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY src/server.js .
RUN chown -R nonroot:nonroot /app
USER nonroot
ENV PORT=8093 \
    REDIS_HOST=redis \
    REDIS_PORT=6379
LABEL org.bstu.student.fullname="Козлович Антон Александрович" \
      org.bstu.student.id="220011" \
      org.bstu.group="as-63" \
      org.bstu.variant="07" \
      org.bstu.course="RSIOT"
EXPOSE $PORT
HEALTHCHECK --interval=15s --timeout=3s --start-period=3s --retries=2 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:$PORT/health || exit 1
CMD ["node", "server.js"]