import { expect, test } from '@playwright/test';

const baseURL = process.env.PLAYWRIGHT_BASE_URL;

test.skip(!baseURL, 'ต้องตั้งค่า PLAYWRIGHT_BASE_URL เพื่อชี้ไปยัง frontend ที่กำลังรันอยู่');

test.describe('viewer streaming journey', () => {
  test('plays movie with captions toggle', async ({ page }) => {
    if (!baseURL) {
      test.skip();
    }

    const tokenResponse = page.waitForResponse(
      (response) =>
        response.url().includes('/playback-token') && response.request().method() === 'POST' && response.status() === 200
    );
    const manifestResponse = page.waitForResponse(
      (response) => response.url().includes('/manifest.m3u8') && response.status() === 200
    );

    await page.goto(`${baseURL}/movies/sample-movie`);
    await Promise.all([tokenResponse, manifestResponse]);

    await expect(page.getByRole('heading', { name: 'ตัวอย่างภาพยนตร์' })).toBeVisible();

    const player = page.getByTestId('movie-player');
    await expect(player).toBeVisible();
    await expect(player).toHaveAttribute('data-status', 'ready', { timeout: 45_000 });

    const captionButton = page.getByRole('button', { name: 'English' });
    await expect(captionButton).toBeVisible();
    await captionButton.click();

    await expect(player).toHaveAttribute('data-caption', 'en', { timeout: 5_000 });
  });
});
