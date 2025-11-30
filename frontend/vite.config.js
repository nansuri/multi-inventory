import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { VantResolver } from 'unplugin-vue-components/resolvers'
import AutoImport from 'unplugin-auto-import/vite'

// https://vitejs.dev/config/
// Allow adding extra hosts via environment variable VITE_ALLOWED_HOSTS (comma separated)
const extraHostsEnv = (process.env.VITE_ALLOWED_HOSTS || '').split(',').map(h => h.trim()).filter(Boolean)
const allowedHosts = ['inventory.justnansuri.com', ...extraHostsEnv]

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [VantResolver()],
    }),
    Components({
      resolvers: [VantResolver()],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  preview: {
    allowedHosts,
  }
})
