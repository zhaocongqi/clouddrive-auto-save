package api

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/zcq/clouddrive-auto-save/internal/core"
	"github.com/zcq/clouddrive-auto-save/internal/core/renamer"
	"github.com/zcq/clouddrive-auto-save/internal/db"
)

type PreviewResult struct {
	OriginalName string `json:"original_name"`
	NewName      string `json:"new_name"`
	Matched      bool   `json:"matched"`
}

func previewTask(c *gin.Context) {
	var req struct {
		AccountID   uint   `json:"account_id" binding:"required"`
		ShareURL    string `json:"share_url" binding:"required"`
		ExtractCode string `json:"extract_code"`
		Pattern     string `json:"pattern"`
		Replacement string `json:"replacement"`
		Name        string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account db.Account
	if err := db.DB.First(&account, req.AccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	driver := core.GetDriver(&account)
	if driver == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "driver not found"})
		return
	}

	ctx := c.Request.Context()
	files, err := driver.ParseShare(ctx, req.ShareURL, req.ExtractCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	renameProc := renamer.NewProcessor()
	var results []PreviewResult

	for _, f := range files {
		res := PreviewResult{
			OriginalName: f.Name,
			NewName:      f.Name,
			Matched:      false,
		}

		// 重命名逻辑
		if req.Pattern != "" || req.Replacement != "" {
			opts := renamer.RenameOptions{
				TaskName:    req.Name,
				Pattern:     req.Pattern,
				Replacement: req.Replacement,
				FileName:    f.Name,
			}
			newName, err := renameProc.Process(opts)
			if err == nil {
				res.NewName = newName
				// 简单的匹配状态判断：如果新名不等于旧名，或者正则能匹配上
				if req.Pattern != "" {
					re, _ := regexp.Compile(req.Pattern)
					res.Matched = re != nil && re.MatchString(f.Name)
				} else {
					res.Matched = newName != f.Name
				}
			}
		}

		results = append(results, res)
	}

	log.Printf("[API] 预览计算完成: 处理了 %d 个文件", len(results))
	c.JSON(http.StatusOK, results)
}

func parseShareLinkInfo(c *gin.Context) {
	var req struct {
		AccountID   uint   `json:"account_id" binding:"required"`
		ShareURL    string `json:"share_url" binding:"required"`
		ExtractCode string `json:"extract_code"`
		SavePath    string `json:"save_path"`
		Pattern     string `json:"pattern"`
		Replacement string `json:"replacement"`
		Name        string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[API] 正在解析分享链接: %s (AccountID: %d, SavePath: %s)", req.ShareURL, req.AccountID, req.SavePath)

	var account db.Account
	if err := db.DB.First(&account, req.AccountID).Error; err != nil {
		log.Printf("[API] 解析失败: 账号未找到 (ID: %d)", req.AccountID)
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	driver := core.GetDriver(&account)
	if driver == nil {
		log.Printf("[API] 解析失败: 驱动加载失败 (Platform: %s)", account.Platform)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "driver not found"})
		return
	}

	ctx := c.Request.Context()
	files, err := driver.ParseShare(ctx, req.ShareURL, req.ExtractCode)
	if err != nil {
		log.Printf("[API] 解析失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 1. 如果有重命名规则，先执行预览重命名
	if req.Pattern != "" || req.Replacement != "" {
		renameProc := renamer.NewProcessor()
		for i := range files {
			opts := renamer.RenameOptions{
				TaskName:    req.Name,
				Pattern:     req.Pattern,
				Replacement: req.Replacement,
				FileName:    files[i].Name,
			}
			newName, err := renameProc.Process(opts)
			if err == nil {
				files[i].NewName = newName
			} else {
				files[i].NewName = files[i].Name
			}
		}
	} else {
		// 默认预览名即原名
		for i := range files {
			files[i].NewName = files[i].Name
		}
	}

	// 2. 执行同名检查逻辑 (基于预览名进行比对)
	if req.SavePath != "" && len(files) > 0 {
		targetID, err := driver.PrepareTargetPath(ctx, req.SavePath)
		if err == nil {
			existingMap := make(map[string]bool)
			existingFiles, _ := driver.ListFiles(ctx, targetID)
			for _, ef := range existingFiles {
				existingMap[ef.Name] = true
			}

			for i := range files {
				// 检查重命名后的名字是否已在网盘中
				if existingMap[files[i].NewName] {
					files[i].IsExisted = true
				}
			}
		}
	}

	log.Printf("[API] 解析完成: 发现 %d 个文件/文件夹", len(files))
	c.JSON(http.StatusOK, files)
}
