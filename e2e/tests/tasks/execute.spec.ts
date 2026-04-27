import { test, expect } from '@playwright/test';
import { add139Account, addQuarkAccount } from '../../fixtures/account.fixture';

test.describe('任务管理：状态机与执行测试', () => {
  test.beforeEach(async ({ page }) => {
    await add139Account(page);
    await addQuarkAccount(page);
  });

  test('手动执行 139 转存任务并验证成功状态', async ({ page }) => {
    // 1. 创建任务
    await page.goto('/tasks');
    await page.getByRole('button', { name: '创建任务' }).last().click();
    await page.locator('.el-select').first().click();
    await page.getByRole('option', { name: 'E2E移动云盘用户' }).first().click();
    await page.getByLabel('任务名称').fill('E2E_139_执行测试');
    await page.getByLabel('分享链接').fill('https://yun.139.com/w/#/share/link/mock_link_id');
    await page.getByLabel('保存路径').fill('/139_exec');
    await page.getByRole('button', { name: '确认并保存' }).click();
    await expect(page.getByText('E2E_139_执行测试')).toBeVisible();

    // 2. 点击执行按钮
    const taskRow = page.locator('tr').filter({ hasText: 'E2E_139_执行测试' });
    await taskRow.getByRole('button', { name: '运行' }).click();

    // 3. 验证状态变更为“SUCCESS” (前端显示大写)
    await expect(taskRow.locator('.el-tag').filter({ hasText: 'SUCCESS' })).toBeVisible({ timeout: 15000 });
  });

  test('手动执行夸克转存任务并验证成功状态', async ({ page }) => {
    // 1. 创建任务
    await page.goto('/tasks');
    await page.getByRole('button', { name: '创建任务' }).last().click();
    await page.locator('.el-select').first().click();
    await page.getByRole('option', { name: 'E2E夸克用户' }).first().click();
    await page.getByLabel('任务名称').fill('E2E_Quark_执行测试');
    await page.getByLabel('分享链接').fill('https://pan.quark.cn/s/mock_link_id');
    await page.getByLabel('保存路径').fill('/quark_exec');
    await page.getByRole('button', { name: '确认并保存' }).click();
    await expect(page.getByText('E2E_Quark_执行测试')).toBeVisible();

    // 2. 点击执行按钮
    const taskRow = page.locator('tr').filter({ hasText: 'E2E_Quark_执行测试' });
    await taskRow.getByRole('button', { name: '运行' }).click();

    // 3. 验证状态变更为“SUCCESS”
    await expect(taskRow.locator('.el-tag').filter({ hasText: 'SUCCESS' })).toBeVisible({ timeout: 15000 });
  });
});
