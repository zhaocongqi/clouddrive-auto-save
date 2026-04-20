# 更新系统名称和网页标题计划

## 1. 目标 (Objective)

将前端界面的网页标签名 (title) 修改为“统一云盘自动转存系统”，并将系统左上角的项目名称从“CloudSaver”修改为“UCAS”。

## 2. 涉及的关键文件 (Key Files & Context)

* **入口 HTML:** `web/index.html`
* **前端布局组件:** `web/src/layout/MainLayout.vue`

## 3. 实施步骤 (Implementation Steps)

### 3.1 修改网页标题

1. 打开 `web/index.html`。
2. 找到 `<title>web</title>` 标签。
3. 将其内容替换为 `<title>统一云盘自动转存系统</title>`。

### 3.2 修改项目显示名称

1. 打开 `web/src/layout/MainLayout.vue`。
2. 在模板的 `.logo` 区域，找到 `<span>CloudSaver</span>`。
3. 将其内容替换为 `<span>UCAS</span>`。

## 4. 验证与测试 (Verification & Testing)

1. 保存文件后，Vite 将自动重新编译并刷新浏览器。
2. 检查浏览器标签页的名称是否已变更为“统一云盘自动转存系统”。
3. 检查页面左上角的系统名称是否已更新为“UCAS”。
