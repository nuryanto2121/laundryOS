# =============================================================================
# Stage 1: Dependencies
# =============================================================================
FROM node:22-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json* ./

# Install production dependencies only
RUN npm ci --omit=dev

# =============================================================================
# Stage 2: Runner
# =============================================================================
FROM node:22-alpine AS runner

# Install wget for Docker health check
RUN apk add --no-cache wget

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy production dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy application source
COPY src ./src
COPY package.json ./

USER appuser

EXPOSE 3000

CMD ["node", "src/index.js"]
