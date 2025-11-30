# Build stage
FROM node:22-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci || npm install
COPY . .
RUN npm run build

# Run stage (vite preview)
FROM node:22-alpine
WORKDIR /app
# Copy built assets and config so vite preview can read preview/server settings
COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/vite.config.js /app/vite.config.js
COPY --from=builder /app/package.json /app/package.json
COPY --from=builder /app/node_modules /app/node_modules
ENV VITE_ALLOWED_HOSTS=inventory.justnansuri.com
RUN npm install -g vite
EXPOSE 5173
CMD ["vite", "preview", "--host", "0.0.0.0", "--port", "5173"]
