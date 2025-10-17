import { test } from '@playwright/test';

test.describe.skip('viewer streaming journey', () => {
  test('plays movie with captions toggle', async ({ page }) => {
    // TODO: implement Playwright scenario once backend endpoints are finalized
    await page.goto('/movies/sample-movie');
  });
});
