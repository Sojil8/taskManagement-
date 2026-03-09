import { resolve } from 'path';
import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        login: resolve(__dirname, 'login.html'),
        signup: resolve(__dirname, 'signup.html'),
        otp: resolve(__dirname, 'otp.html'),
        dashboard: resolve(__dirname, 'dashboard.html'),
        pomodoro: resolve(__dirname, 'pomodoro.html'),
      },
    },
  },
});
