import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './frontend/tests/e2e',
  testMatch: /.*\.spec\.ts/,
  fullyParallel: true,
  timeout: 60_000,
  expect: {
    timeout: 10_000
  },
  use: {
    baseURL: process.env.PLAYWRIGHT_BASE_URL ?? 'http://localhost:3000',
    video: 'retain-on-failure',
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure'
  },
  reporter: [['list']]
});
