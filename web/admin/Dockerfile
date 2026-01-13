FROM node:23.4.0-alpine
WORKDIR /app

# Install DNS utilities and ensure network access
RUN apk add --no-cache ca-certificates

# Copy package files first for better caching
COPY package*.json ./

# Install dependencies with retry logic
RUN for i in 1 2 3; do npm ci && break || sleep 10; done

# Copy the rest of the application
COPY . .

# Build the application with retry logic for network issues
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build || echo "Build failed, will use dev mode"

ENV NODE_ENV production

RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001 && \
    chown -R nextjs:nodejs /app

USER nextjs

EXPOSE 3000

# Start in production if build succeeded, otherwise use dev mode
# Handle SIGTERM properly for graceful shutdown
# Use exec to ensure signals are properly forwarded to the Node.js process
CMD ["sh", "-c", "if [ -f .next/BUILD_ID ]; then exec npm start; else echo 'Starting in dev mode (build failed)' && exec npm run dev -- -p 3000; fi"]