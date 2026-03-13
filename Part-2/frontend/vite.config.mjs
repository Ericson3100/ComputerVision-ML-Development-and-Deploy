import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import http from 'node:http';

export default defineConfig({
  server: {
    port: 3000,
    host: '0.0.0.0',
    watch: {
      usePolling: true,
      interval: 1000
    },
    proxy: {
      '/api': {
        target: process.env.VUE_APP_API_URL || 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  plugins: [
    vue(),
    {
      name: 'sse-classify-stream-proxy',
      configureServer(server) {
        server.middlewares.use('/api/classify-stream', (req, res) => {
          if (req.method !== 'POST') {
            return;
          }

          const backendUrl = process.env.VUE_APP_API_URL || 'http://localhost:8080';
          const url = new URL(backendUrl);

          const options = {
            hostname: url.hostname,
            port: url.port || 8080,
            path: '/api/classify-stream',
            method: 'POST',
            headers: {
              ...req.headers,
              host: url.hostname
            }
          };

          const proxyReq = http.request(options, (proxyRes) => {
            res.writeHead(proxyRes.statusCode || 500, {
              'Content-Type': 'text/event-stream',
              'Cache-Control': 'no-cache, no-transform',
              Connection: 'keep-alive',
              'X-Accel-Buffering': 'no',
              'Access-Control-Allow-Origin': '*'
            });

            if (res.socket) {
              res.socket.setNoDelay(true);
              res.socket.setKeepAlive(true);
            }

            proxyRes.pipe(res);
          });

          proxyReq.on('error', (err) => {
            console.error('Error in SSE proxy:', err);
            if (!res.headersSent) {
              res.statusCode = 500;
            }
            res.end();
          });

          req.pipe(proxyReq);
        });
      }
    }
  ]
});
