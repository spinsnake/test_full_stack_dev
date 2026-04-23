import { expect, test } from '@playwright/test';
import { installMockApi } from './support/mockApi';

test('creates a new image and attaches tags during creation', async ({ page }) => {
  await installMockApi(page);

  await page.goto('/manage');

  await expect(page.getByRole('heading', { name: 'Manage images, tags, and assignments.' })).toBeVisible();
  await expect(page.getByTestId('image-catalog')).toContainText('3 active images');

  await page.getByLabel('Image URL').fill('https://example.com/playwright-sunset.jpg');
  await page.getByLabel('Thumbnail URL').fill('https://example.com/playwright-sunset-thumb.jpg');
  await page.getByLabel('Width').fill('1600');
  await page.getByLabel('Height').fill('900');
  await page.getByLabel('Alt Text').fill('Playwright Sunset');
  await page.getByLabel('Source').fill('playwright-e2e');

  const createTagPicker = page.getByTestId('create-image-tag-picker');
  await createTagPicker.getByRole('button', { name: /#nature/i }).click();
  await createTagPicker.getByRole('button', { name: /#city/i }).click();

  await page.getByRole('button', { name: 'Create Image' }).click();

  await expect(page.getByTestId('manage-status')).toContainText('Image created with 2 tag(s).');
  await expect(page.getByTestId('image-catalog')).toContainText('4 active images');

  const createdCard = page
    .getByTestId('image-catalog')
    .locator('article')
    .filter({ hasText: 'Playwright Sunset' })
    .first();

  await expect(createdCard).toContainText('#nature');
  await expect(createdCard).toContainText('#city');
  await expect(page.getByLabel('Alt Text')).toHaveValue('Playwright Sunset');
});
