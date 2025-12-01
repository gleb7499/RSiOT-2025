# --------- builder ----------
FROM node:22-alpine AS builder
WORKDIR /app

# метаданные прокинем позже через --label в финальном образе
COPY package*.json ./
RUN npm ci --omit=dev

COPY src ./src

FROM node:22-alpine AS runner
WORKDIR /app

RUN addgroup -S app && adduser -S -u 10001 -G app app

ENV NODE_ENV=production
ENV STU_ID=STUDENT_ID \
    STU_GROUP=GROUP \
    STU_VARIANT=v16 \
    PORT=8064

COPY --from=builder /app/node_modules ./node_modules
COPY src ./src
COPY package.json ./

LABEL org.bstu.student.fullname="Никифоров Александр" \
      org.bstu.student.id="220020" \
      org.bstu.group="АС-63" \
      org.bstu.variant="16" \
      org.bstu.course="RSIOT"

EXPOSE 8064

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8064/health >/dev/null 2>&1 || exit 1

USER 10001

CMD ["node", "src/server.js"]
