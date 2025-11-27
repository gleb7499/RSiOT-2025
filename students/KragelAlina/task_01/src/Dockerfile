# Stage 1: Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files first for caching
COPY src/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the source code
COPY src/ ./

# Stage 2: Runtime stage (minimal image)
FROM node:20-alpine

# Add labels
LABEL org.bstu.student.fullname="Крагель Алина Максимовна" \
      org.bstu.student.id="220046" \
      org.bstu.group="АС-63" \
      org.bstu.variant="10" \
      org.bstu.course="RSIOT"

# Create non-root user (UID 10001)
RUN adduser -u 10001 -D appuser

# Set working directory
WORKDIR /app

# Copy built artifacts from builder stage
COPY --from=builder /app /app

# Switch to non-root user
USER 10001

# Expose port
EXPOSE 8072

# Healthcheck endpoint (/ping for liveness/readiness)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8072/ping || exit 1

# Run the app (configurable via env vars)
CMD ["npm", "start"]