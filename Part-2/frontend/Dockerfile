# Multi-stage build for Vue frontend
FROM node:22-alpine AS builder

# Set working directory
WORKDIR /app

# Use latest npm
RUN npm install -g npm@latest

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Production build
RUN npm run build

# Production stage with nginx
FROM nginx:alpine

# Copy build files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Default command
CMD ["nginx", "-g", "daemon off;"]
