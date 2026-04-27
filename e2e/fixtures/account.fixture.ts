import { Page, expect } from '@playwright/test';

export async function add139Account(page: Page, authStr: string = 'mock_auth', username: string = 'E2E移动云盘用户') {
  await page.goto('/accounts');
  // 如果已经存在，则不再重复添加（使用 first 解决冲突）
  if (await page.getByText(username).first().isVisible()) return;

  await page.getByRole('button', { name: /立即绑定账号|添加账号/ }).first().click();
  // 使用更精确的定位器避免与列表冲突
  await page.getByLabel('网盘平台').getByText('移动云盘').click();
  await page.getByLabel('Authorization').fill(authStr);
  await page.getByRole('button', { name: '确认添加' }).click();
  await expect(page.getByText(username).first()).toBeVisible({ timeout: 10000 });
}

export async function addQuarkAccount(page: Page, cookieStr: string = '__uid=mock; mock_cookie', username: string = 'E2E夸克用户') {
  await page.goto('/accounts');
  if (await page.getByText(username).first().isVisible()) return;

  await page.getByRole('button', { name: /立即绑定账号|添加账号/ }).first().click();
  // 使用更精确的定位器避免与列表冲突
  await page.getByLabel('网盘平台').getByText('夸克网盘').click();
  await page.getByLabel('Cookie 全量字符串').fill(cookieStr);
  await page.getByRole('button', { name: '确认添加' }).click();
  await expect(page.getByText(username).first()).toBeVisible({ timeout: 10000 });
}
