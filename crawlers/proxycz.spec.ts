import { test, expect } from '@playwright/test';

test('scrap', async ({ page, context }) => {
    await page.goto('http://free-proxy.cz/en/proxylist/country/all/https/uptime/all')
})
