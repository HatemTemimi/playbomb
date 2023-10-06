import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  // Look for test files in the "tests" directory, relative to this configuration file.
  testDir: 'crawlers',

  // Run all tests in parallel.
  fullyParallel: true,

  // Retry on CI only.
  retries: 2,

  // Opt out of parallel tests on CI.
  workers: 6,

  // Reporter to use
  reporter: 'html',

  use: {
    // Collect trace when retrying the failed test.
    trace: 'on-first-retry',
    video:'on',
  },
  // Configure projects for major browsers.
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
