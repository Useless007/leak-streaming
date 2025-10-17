import { expect, test } from '@playwright/test';

const baseURL = process.env.PLAYWRIGHT_BASE_URL;

test.skip(!baseURL, 'ต้องตั้งค่า PLAYWRIGHT_BASE_URL เพื่อชี้ไปยัง frontend ที่กำลังรันอยู่');

test.describe('admin movie management', () => {
	test('creates a new movie entry via admin form', async ({ page }) => {
		if (!baseURL) {
			test.skip();
		}

		const uniqueId = Date.now();
		const title = `Integration Premiere ${uniqueId}`;
		const streamUrl = 'https://main.24playerhd.com/m3u8/0378b65549cda348e910faf0/0378b65549cda348e910faf0168.m3u8';

		await page.goto(`${baseURL}/admin/movies/new`);

		await page.fill('input[name="title"]', title);
		await page.fill('textarea[name="synopsis"]', 'ภาพยนตร์สำหรับการทดสอบการสร้างผ่านหน้าแอดมิน');
		await page.fill('input[name="posterUrl"]', 'https://images.unsplash.com/photo-1524985069026-dd778a71c7b4?w=800');
		await page.fill('input[name="streamUrl"]', streamUrl);
		await page.fill('textarea[name="allowedHosts"]', ['main.24playerhd.com', 'm42.winplay4.com'].join('\n'));

		await page.getByRole('button', { name: 'เพิ่มคำบรรยาย' }).click();
		await page.fill('input[name="captions.0.languageCode"]', 'th');
		await page.fill('input[name="captions.0.label"]', 'Thai');
		await page.fill('input[name="captions.0.captionUrl"]', '/captions/sample-en.vtt');

		await page.getByRole('button', { name: 'บันทึกภาพยนตร์' }).click();

		await page.waitForURL(/\/movies\//, { timeout: 30_000 });
		await expect(page.getByRole('heading', { name: title })).toBeVisible();
	});
});
