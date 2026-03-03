# ---------- Build Stage ----------
FROM node:20-alpine AS builder

WORKDIR /app

# Enable corepack (for pnpm)
RUN corepack enable

ARG VITE_WS_BASE_URL
ENV VITE_WS_BASE_URL=$VITE_WS_BASE_URL

# Copy only dependency files first (better layer caching)
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy rest of project
COPY . .

# Build production bundle
RUN pnpm build


# ---------- Runtime Stage ----------
FROM node:20-alpine

WORKDIR /app

# Install lightweight static server
RUN npm install -g serve

# Copy built output
COPY --from=builder /app/dist ./dist

EXPOSE 3000

CMD ["serve", "-s", "dist", "-l", "3000"]