import { expect, test } from '@playwright/test';

const baseURL = process.env.PLAYWRIGHT_BASE_URL;

test.skip(!baseURL, 'ต้องตั้งค่า PLAYWRIGHT_BASE_URL เพื่อชี้ไปยัง frontend ที่กำลังรันอยู่');

test.describe('movie catalogue', () => {
  test('lists movies and opens detail page', async ({ page }) => {
    if (!baseURL) {
      test.skip();
    }

    await page.goto(`${baseURL}/movies`);

    await expect(page.locator('h1')).toContainText('เลือกภาพยนตร์ที่คุณอยากรับชม');

	const cards = page.getByTestId('movie-card');
	await expect(cards).not.toHaveCount(0);

    const secondMovieLink = page.locator('[data-movie-slug="demo-movie-2"] a').first();
    await secondMovieLink.click();

    await expect(page).toHaveURL(/\/movies\/demo-movie-2$/);
    const player = page.getByTestId('movie-player');
    await expect(player).toBeVisible();
  });
});
