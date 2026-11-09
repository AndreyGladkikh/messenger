import path from 'node:path'
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: { '@': path.resolve(import.meta.dirname, './src') },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    env: {
      VITE_API_URL: 'https://api.test',
      VITE_WS_URL: 'wss://api.test/ws',
    },
  },
})
