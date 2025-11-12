/// <reference types="vitest" />

import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

import { resolve } from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: {
      provider: 'v8', // デフォルトは 'v8'
      reporter: ['text', 'html', 'lcov'], // 出力形式
      reportsDirectory: './coverage' // 保存先
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  build: {
    rollupOptions: {
      input: {
        "index": resolve(__dirname, 'htmls', 'index.html'),
        "apply": resolve(__dirname, 'htmls', 'apply.html'),
        "role": resolve(__dirname, 'htmls', 'role.html'),
        "manage": resolve(__dirname, 'htmls', 'manage.html'),
        "sign-up": resolve(__dirname, 'htmls', 'sign-up.html'),
        "sign-in": resolve(__dirname, 'htmls', 'sign-in.html'),
      },
      output: {
       entryFileNames: 'assets/[name].js',
       chunkFileNames: 'assets/[name].js',
       assetFileNames: 'assets/[name].[ext]',
     }
    }
  }
})
