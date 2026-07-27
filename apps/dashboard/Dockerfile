# =============================================================================
# Stage 1: Dependencies & Build
# =============================================================================
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./

RUN npm ci

COPY . .

# =============================================================================
# Stage 2: Runner (Vite dev server)
# =============================================================================
FROM node:22-alpine AS runner

# Install wget for Docker health check
RUN apk add --no-cache wget

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy dependencies and source from builder
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app ./

# Give appuser ownership so Vite can write .vite/deps_temp at runtime
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 5173

# Vite dev server — host:0.0.0.0 is configured in vite.config.js
CMD ["npm", "run", "dev"]
