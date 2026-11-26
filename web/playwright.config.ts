import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  fullyParallel: true,
  use: { baseURL: "http://127.0.0.1:3102", trace: "retain-on-failure" },
  projects: [
    { name: "edge", use: { ...devices["Desktop Chrome"], channel: "msedge" } }
  ],
  webServer: {
    command: "npm run start -- --hostname 127.0.0.1 --port 3102",
    url: "http://127.0.0.1:3102/images",
    reuseExistingServer: !process.env.CI,
    timeout: 30_000
  }
});
