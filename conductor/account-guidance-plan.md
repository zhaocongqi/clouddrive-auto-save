# 账号添加文档指引计划 (Account Guidance Plan)

**Goal:** 在“添加新账号”对话框中，为 Authorization 和 Cookie
表单项的标题栏（label）右侧添加带有图标的“如何获取？”超链接，引导用户查看相关文档。

---

## Task 1: 引入图标并定义帮助方法

**Files:**

- Modify: `web/src/views/Accounts.vue` (Script 部分)

**Changes:**

1. 从 `lucide-vue-next` 中导入 `ExternalLink` 图标。
2. 创建帮助链接生成方法，根据不同平台或不同认证类型返回相应的文档 URL：
   - 移动云盘 Authorization ->
     `https://doc.oplist.org/guide/drivers/139#authorization-1`
   - 移动云盘 Cookie -> `https://doc.oplist.org/guide/drivers/139#cookie-1`
   - 夸克网盘 Cookie -> `https://doc.oplist.org/guide/drivers/quark#cookie-1`

## Task 2: 改造表单 Label 模板

**Files:**

- Modify: `web/src/views/Accounts.vue` (Template 部分)

**Changes:**

1. 将 `<el-form-item label="Authorization">` 改造为：

   ```html
   <el-form-item>
     <template #label>
       <div class="form-label-with-link">
         <span>Authorization</span>
         <el-link 
           type="primary" 
           :underline="false" 
           href="https://doc.oplist.org/guide/drivers/139#authorization-1" 
           target="_blank" 
           class="help-link"
         >
           <el-icon><ExternalLink /></el-icon>
           如何获取？
         </el-link>       </div>
     </template>
     ...
   </el-form-item>
   ```

2. 将 `<el-form-item label="Cookie 全量字符串">` 改造为：

   ```html
   <el-form-item>
     <template #label>
       <div class="form-label-with-link">
         <span>Cookie 全量字符串</span>
         <el-link 
           type="primary" 
           :underline="false" 
           :href="accountForm.platform === '139' ? 'https://doc.oplist.org/guide/drivers/139#cookie-1' : 'https://doc.oplist.org/guide/drivers/quark#cookie-1'" 
           target="_blank" 
           class="help-link"
         >
           <el-icon><ExternalLink /></el-icon>
           如何获取？
         </el-link>
       </div>
     </template>
     ...
   </el-form-item>
   ```

## Task 3: 补充对齐样式

**Files:**

- Modify: `web/src/views/Accounts.vue` (Style 部分)

**Changes:**

1. 添加 `.form-label-with-link` 样式：

   ```css
   .form-label-with-link {
     display: flex;
     justify-content: space-between;
     align-items: center;
     width: 100%;
   }
   .help-link {
     font-size: 12px;
     font-weight: normal;
   }
   .help-link .el-icon {
     margin-right: 4px;
   }
   ```
