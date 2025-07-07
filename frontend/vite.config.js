import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    lib: {
      entry: 'src/index.ts',
      name: 'FernFS',
      formats: ['es', 'umd']
    },
    sourcemap: true,
    outDir: 'dist'
  },
  server: {
    port: 5173,
    host: true
  },
}) 