FROM node:23.4.0-alpine
WORKDIR /app
COPY . .

RUN npm ci -g npm@latest && npm run build || true

ENV NODE_ENV production

RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001

USER nextjs

EXPOSE 3000

CMD ["npm", "start"] 