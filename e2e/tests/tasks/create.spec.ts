import { test, expect } from '@playwright/test';
import { add139Account, addQuarkAccount } from '../../fixtures/account.fixture';

test.describe('任务管理：创建功能测试', () => {
  test.beforeEach(async ({ page }) => {
    await add139Account(page);
    await addQuarkAccount(page);
  });

  test('创建 139 移动云盘转存任务', async ({ page }) => {
    await page.goto('/tasks');
    await page.getByRole('button', { name: '创建任务' }).last().click();

    await page.locator('.el-select').first().click();
    await page.getByRole('option', { name: 'E2E移动云盘用户' }).first().click();
    
    await page.getByLabel('任务名称').fill('E2E_139_转存任务');
    await page.getByLabel('分享链接').fill('https://yun.139.com/w/#/share/link/mock_link_id');
    await page.getByLabel('保存路径').fill('/139_sync_folder');
    
    await page.getByRole('button', { name: '确认并保存' }).click();

    // 验证回到任务列表并展示该任务
    await expect(page.getByText('E2E_139_转存任务')).toBeVisible({ timeout: 10000 });
  });

  test('创建夸克网盘转存任务', async ({ page }) => {
    await page.goto('/tasks');
    await page.getByRole('button', { name: '创建任务' }).last().click();

    await page.locator('.el-select').first().click();
    await page.getByRole('option', { name: 'E2E夸克用户' }).first().click();
    
    await page.getByLabel('任务名称').fill('E2E_Quark_转存任务');
    await page.getByLabel('分享链接').fill('https://pan.quark.cn/s/mock_link_id');
    await page.getByLabel('保存路径').fill('/quark_sync_folder');
    
    await page.getByRole('button', { name: '确认并保存' }).click();

    // 验证回到任务列表并展示该任务
    await expect(page.getByText('E2E_Quark_转存任务')).toBeVisible({ timeout: 10000 });
  });
});
