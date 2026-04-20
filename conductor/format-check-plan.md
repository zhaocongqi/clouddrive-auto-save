# 代码与文档格式检查修复计划 (Format Check & Fix Plan)

**Goal:** 执行 `make lint` 和 `make lint-md` 检查代码和 Markdown
文件的格式规范，并根据检查结果自动修复不合规的地方。

---

## Task 1: 运行格式检查

**Actions:**

- 运行 `make lint` 检查 Go 代码格式。
- 运行 `make lint-md` 检查 Markdown 文档格式。

### Task 2: 修复格式问题

**Actions:**

- 针对 `make lint` 报告的 Go 代码问题，使用 `gofmt -w .` 自动修复，或手动调整无法自动修复的代码。
- 针对 `make lint-md` 报告的 Markdown 问题，使用 `npx markdownlint-cli "**/*.md" --fix
  --ignore "node_modules" --ignore "web/node_modules"` 自动修复，并根据需要手动调整。

### Task 3: 提交更改

- 将修复后的文件提交到本地仓库。
