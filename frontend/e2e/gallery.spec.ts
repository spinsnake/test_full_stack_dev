import { expect, test } from '@playwright/test';
import { installMockApi } from './support/mockApi';

test('filters gallery items by tag and opens manage page', async ({ page }) => {
  await installMockApi(page);

  await page.goto('/');

  await expect(page.getByRole('heading', { name: 'Masonry image gallery.' })).toBeVisible();
  await expect(page.getByTestId('gallery-grid').locator('[data-testid^="gallery-card-"]')).toHaveCount(3);
  await expect(page.getByText('Forest Path')).toBeVisible();
  await expect(page.getByText('City Glow')).toBeVisible();

  await page.getByRole('button', { name: 'Tags Filter' }).click();

  const tagModal = page.getByTestId('gallery-tag-modal');
  await expect(tagModal).toBeVisible();
  await tagModal.getByRole('button', { name: /#nature/i }).click();

  await expect(page.getByTestId('gallery-grid').locator('[data-testid^="gallery-card-"]')).toHaveCount(1);
  await expect(page.getByText('Forest Path')).toBeVisible();
  await expect(page.getByText('City Glow')).toHaveCount(0);
  await expect(page.getByRole('banner').getByText('#nature')).toBeVisible();

  await page.getByRole('link', { name: 'Manage Gallery' }).click();
  await expect(page).toHaveURL(/\/manage$/);
  await expect(page.getByRole('heading', { name: 'Manage images, tags, and assignments.' })).toBeVisible();
});
