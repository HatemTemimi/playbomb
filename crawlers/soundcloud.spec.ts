import { test, expect } from '@playwright/test';
import UserAgent from 'user-agents';
//import proxies from './proxylist';
const proxies = require('../proxy.json')

const timezones = [
  "Europe/Berlin",
  "Europe/Paris",
  "Europe/Bucharest",
  "Europe/Dublin",
  "Europe/Helsinki",
  "Europe/Kiev",
  "Europe/London",
  "Europe/Warsaw",
  "Europe/Vienna",
]

const local = [
  "de-DE",
  "fr-FR",
  "en-EN",
  "de-CH",
  "de-BE",
  "da-DK",
  "en-AT",
  "en-AU",
  "en-BZ",
  "fr-CA",
  "fr-FR",
  "fr-GP",
]

for (let i=0;i<200;i++){

  console.log(proxies)

  //generate random desktop user agent
  const userAgent = new UserAgent({ deviceCategory: 'desktop' });

  //pick random timezoneID / local / proxy
  const tz = timezones[Math.floor(Math.random()*timezones.length)];
  const loc = local[Math.floor(Math.random()*local.length)];
  const proxy = proxies[Math.floor(Math.random()*proxies.length)];

  console.log('using proxy: ', proxy, 'for test: ', i+1 )

  test.use({
    locale: loc,
    timezoneId: tz,
    userAgent: userAgent.toString(),
    proxy: {
      server: proxy,
      bypass: 'localhost',
    },
    ignoreHTTPSErrors: true
  });

  test.setTimeout(60000)

  test(`play sc ${i+1}`, async ({ page, context }) => {

    //disable web driver to avoid detection
    await context.addInitScript("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})")

    await page.goto('https://soundcloud.com/rxb-flac/break-side-i', {
      timeout: 35000
    });

    //optionally accepts the initial cookie request !fires only if cookie not set
    /*
     *
    const accept =  page.locator('#onetrust-accept-btn-handler')

    await expect.soft(accept).toBeVisible({ timeout: 20000 })

    await expect.soft(accept.click()).toBeTruthy()

    */

    const play = page.locator('.sc-button-play').first()

    await expect(play).toBeVisible({ timeout: 20000 })

    await expect(play).toBeEnabled({ timeout: 20000 })

    await new Promise(r => setTimeout(r, 5000));

    console.log('finished play: ', i+1)

  });

}

