import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import svgr from 'vite-plugin-svgr' 

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        react(),
        svgr(),
    ],
    define: {
        'process.env': process.env,
    },
    server: {
        host: true,
        port: 5173,
        watch: {
            usePolling: true
        }
    },
    base: './',
});
