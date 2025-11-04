# syntax=docker/dockerfile:1

#########################
# Stage 1: dependencies
#########################
FROM node:20-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN --mount=type=cache,target=/root/.npm npm ci --omit=dev

#########################
# Stage 2: build
#########################
FROM node:20-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

#########################
# Stage 3: final
#########################
FROM node:20-alpine AS runtime
WORKDIR /app
COPY --from=build /app /app

# Удалим ненужные кэши
RUN rm -rf /app/node_modules/.cache && npm prune --omit=dev

USER 65532:65532
EXPOSE 8081

LABEL org.bstu.student.fullname="Vashchuk Anatoli Valerevich" \
    org.bstu.student.id="1" \
    org.bstu.group="feis" \
    org.bstu.variant="v01" \
    org.bstu.course="RSIOT"

ENV STU_ID=1 \
    STU_GROUP=feis \
    STU_VARIANT=v01 \
    PORT=8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=20s --retries=3 \
    CMD wget -q --spider http://localhost:8081/health || exit 1

ENTRYPOINT ["node", "src/server.js"]

