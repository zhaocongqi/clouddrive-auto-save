import { test, expect } from '@playwright/test';

test.describe('仪表盘：数据概览测试', () => {
  test('场景一：空数据状态 (Empty State)', async ({ page }) => {
    // 拦截 /api/dashboard/stats 接口，返回全 0 数据
    await page.route('**/api/dashboard/stats', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          scheduled_tasks: 0,
          capacity_used: 0,
          today_completed: 0,
          active_accounts: 0,
          recent_activities: [],
          running_tasks_list: []
        }),
      });
    });

    await page.goto('/');
    
    // 验证基础标题
    await expect(page.locator('body')).toContainText('云端转存概览');
    
    // 验证四个统计卡片的值
    await expect(page.locator('.stat-card').nth(0).locator('.stat-value')).toHaveText('0');
    await expect(page.locator('.stat-card').nth(1).locator('.stat-value')).toHaveText('0 B');
    await expect(page.locator('.stat-card').nth(2).locator('.stat-value')).toHaveText('0');
    await expect(page.locator('.stat-card').nth(3).locator('.stat-value')).toHaveText('0');
    
    // 验证空状态提示
    await expect(page.getByText('暂无活动记录')).toBeVisible();
  });

  test('场景二：数据填充状态 (Populated Data State)', async ({ page }) => {
    await page.route('**/api/dashboard/stats', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          scheduled_tasks: 15,
          capacity_used: 1073741824, // 1GB
          today_completed: 5,
          active_accounts: 2,
          recent_activities: [],
          running_tasks_list: []
        }),
      });
    });

    await page.goto('/');

    await expect(page.locator('.stat-card').nth(0).locator('.stat-value')).toHaveText('15');
    // 验证前端 formatSize 的转换
    await expect(page.locator('.stat-card').nth(1).locator('.stat-value')).toHaveText('1.00 GB');
    await expect(page.locator('.stat-card').nth(2).locator('.stat-value')).toHaveText('5');
    await expect(page.locator('.stat-card').nth(3).locator('.stat-value')).toHaveText('2');
  });

  test('场景三：实时执行与近期动态状态', async ({ page }) => {
    await page.route('**/api/dashboard/stats', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          scheduled_tasks: 10,
          capacity_used: 1024,
          today_completed: 1,
          active_accounts: 1,
          running_tasks_list: [
            {
              id: 1,
              name: '正在下载电影合集',
              percent: 45,
              stage: 'Downloading',
              message: '正在下载 第3集...'
            }
          ],
          recent_activities: [
            {
              id: 2,
              name: '失效的转存任务',
              status: 'failed',
              last_run: new Date(Date.now() - 3600000).toISOString() // 1 hour ago
            }
          ]
        }),
      });
    });

    await page.goto('/');

    // 验证活跃任务 Tag
    await expect(page.getByText('1 活跃')).toBeVisible();
    
    // 验证活跃任务详情
    await expect(page.getByText('正在下载电影合集')).toBeVisible();
    await expect(page.getByText('Downloading')).toBeVisible();
    await expect(page.getByText('正在下载 第3集...')).toBeVisible();

    // 验证近期失败记录及重试按钮
    await expect(page.getByText('失效的转存任务')).toBeVisible();
    await expect(page.getByRole('button', { name: '重试' })).toBeVisible();
  });
});
