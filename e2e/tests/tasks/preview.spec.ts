import { test, expect } from '@playwright/test';
import { add139Account, addQuarkAccount } from '../../fixtures/account.fixture';

test.describe('任务管理：重命名预览测试', () => {
  test.beforeEach(async ({ page }) => {
    await add139Account(page);
    await addQuarkAccount(page);
  });

  test('验证 139 移动云盘分享链接解析与重命名预览', async ({ page }) => {
    await page.goto('/tasks');
    await page.getByRole('button', { name: '创建任务' }).last().click();
    
    await page.locator('.el-select').first().click();
    await page.getByRole('option', { name: 'E2E移动云盘用户' }).first().click();

    await page.getByLabel('任务名称').fill('139预览测试');
    await page.getByLabel('分享链接').fill('https://yun.139.com/w/#/share/link/mock_link_id');
    await page.getByPlaceholder('匹配文件名的正则表达式').fill('.*\\.mp4$');
    await page.getByPlaceholder('支持 {TASKNAME}, {YEAR} 等变量').fill('[{DATE}] {TASKNAME}.{EXT}');
    
    await page.getByRole('button', { name: '全量重命名预览' }).click();

    const previewDialog = page.getByRole('dialog', { name: '重命名预览' });
    await expect(previewDialog).toBeVisible({ timeout: 15000 });
    // 验证能够抓取到文件，并且正则和变量替换正常工作
    await expect(previewDialog.getByText('[20240420] 139预览测试.mp4').first()).toBeVisible();
    
    // 按 Esc 关闭预览
    await page.keyboard.press('Escape');
    await expect(previewDialog).not.toBeVisible();
  });

  test('验证夸克网盘分享链接解析与重命名预览', async ({ page }) => {
    await page.goto('/tasks');
    await page.getByRole('button', { name: '创建任务' }).last().click();
    
    await page.locator('.el-select').first().click();
    await page.getByRole('option', { name: 'E2E夸克用户' }).first().click();

    await page.getByLabel('任务名称').fill('夸克预览测试');
    await page.getByLabel('分享链接').fill('https://pan.quark.cn/s/mock_link_id');
    await page.getByPlaceholder('匹配文件名的正则表达式').fill('.*\\.txt$');
    await page.getByPlaceholder('支持 {TASKNAME}, {YEAR} 等变量').fill('{TASKNAME}_已修改.{EXT}');
    
    await page.getByRole('button', { name: '全量重命名预览' }).click();

    const previewDialog = page.getByRole('dialog', { name: '重命名预览' });
    await expect(previewDialog).toBeVisible({ timeout: 15000 });
    // 验证能够抓取到文件，并且正则和变量替换正常工作
    await expect(previewDialog.getByText('夸克预览测试_已修改.txt').first()).toBeVisible();
    
    await page.keyboard.press('Escape');
    await expect(previewDialog).not.toBeVisible();
  });
});
