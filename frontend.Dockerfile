# ── Build Stage ──────────────────────────────────────────────────────────────
FROM node:22-alpine AS builder

WORKDIR /app

# Vite needs this at BUILD time so the compiled JS knows the API base path
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

ARG VITE_CONTRACT_ADDRESS
ENV VITE_CONTRACT_ADDRESS=${VITE_CONTRACT_ADDRESS}

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

# ── Runtime Stage ────────────────────────────────────────────────────────────
FROM nginx:1.27-alpine

# Remove default Nginx site
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom Nginx config (reverse proxy + SPA)
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# Copy built React app
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
