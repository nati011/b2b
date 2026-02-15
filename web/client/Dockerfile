FROM node:23.4.0-alpine
WORKDIR /app

# Install DNS utilities and ensure network access
RUN apk add --no-cache ca-certificates

# Copy package files first for better caching
COPY package*.json ./

# Install dependencies with retry logic
RUN for i in 1 2 3; do npm ci --legacy-peer-deps && break || (npm install --legacy-peer-deps && break) || sleep 10; done

# Copy the rest of the application
COPY . .

# Build args for Next.js public env (inlined at build time)
ARG NEXT_PUBLIC_BASE_URL=/api
ENV NEXT_PUBLIC_BASE_URL=$NEXT_PUBLIC_BASE_URL

# Build the application
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

ENV NODE_ENV production

RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001 && \
    chown -R nextjs:nodejs /app

USER nextjs

EXPOSE 3000

# Start in production mode
# Use exec to ensure proper signal forwarding
CMD ["sh", "-c", "exec npm start"] 