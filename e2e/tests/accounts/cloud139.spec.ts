import { test, expect } from '@playwright/test';
import { add139Account } from '../../fixtures/account.fixture';

test.describe('139 移动云盘账号管理', () => {
  test('成功绑定并展示 139 黄金会员账号', async ({ page }) => {
    await add139Account(page);
    await expect(page.getByText('E2E移动云盘用户')).toBeVisible();
    await expect(page.getByText('黄金会员').last()).toBeVisible();
    await expect(page.getByText('512').last()).toBeVisible(); // 512GB (MB 转换 GB 显示由于页面逻辑)
  });

  test('成功绑定并展示 139 普通用户小容量账号', async ({ page }) => {
    await add139Account(page, 'mock_normal', 'E2E139普通用户');
    await expect(page.getByText('E2E139普通用户')).toBeVisible();
    await expect(page.getByText('普通用户').last()).toBeVisible();
    await expect(page.getByText('10').last()).toBeVisible(); // 10GB used (20 - 10)
  });

  test('成功绑定并展示 139 超容状态账号', async ({ page }) => {
    await add139Account(page, 'mock_overcap', 'E2E139超容用户');
    await expect(page.getByText('E2E139超容用户')).toBeVisible();
    await expect(page.getByText('钻石会员').last()).toBeVisible();
    await expect(page.getByText('2 TB').last()).toBeVisible(); // 2TB used
    await expect(page.getByText('已超额 1 TB').last()).toBeVisible();
  });

  test('成功绑定并展示 139 白银会员账号', async ({ page }) => {
    await add139Account(page, 'mock_silver', 'E2E139白银会员');
    await expect(page.getByText('E2E139白银会员')).toBeVisible();
    await expect(page.getByText('白银会员').last()).toBeVisible();
  });
});
