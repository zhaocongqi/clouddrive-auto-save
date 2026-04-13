package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zcq/clouddrive-auto-save/internal/core"
	"github.com/zcq/clouddrive-auto-save/internal/db"
	"log"
	"net/http"
)

// FolderItem 为前端 TreeSelect 提供的结构
type FolderItem struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	Label  string `json:"label"`
	IsLeaf bool   `json:"isLeaf"`
}

func getAccountFolders(c *gin.Context) {
	id := c.Param("id")
	parentID := c.DefaultQuery("parent_id", "")
	parentPath := c.DefaultQuery("parent_path", "/")

	log.Printf("[API] 正在获取账号目录树: AccountID=%s, ParentID=%s, ParentPath=%s", id, parentID, parentPath)

	var account db.Account
	if err := db.DB.First(&account, id).Error; err != nil {
		log.Printf("[API] 获取目录失败: 账号未找到 (ID: %s)", id)
		c.PureJSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	driver := core.GetDriver(&account)
	if driver == nil {
		log.Printf("[API] 获取目录失败: 驱动加载失败 (Platform: %s)", account.Platform)
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": "driver not found"})
		return
	}

	ctx := c.Request.Context()
	files, err := driver.ListFiles(ctx, parentID)
	if err != nil {
		log.Printf("[API] 获取目录异常: %v", err)
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var folders []FolderItem
	for _, f := range files {
		if f.IsFolder {
			childPath := parentPath
			if childPath == "/" {
				childPath = "/" + f.Name
			} else {
				childPath = childPath + "/" + f.Name
			}
			folders = append(folders, FolderItem{
				ID:     f.ID,
				Path:   childPath,
				Label:  f.Name,
				IsLeaf: false,
			})
		}
	}

	log.Printf("[API] 获取目录完成: AccountID=%s, 发现 %d 个子文件夹", id, len(folders))
	c.PureJSON(http.StatusOK, folders)
}

func createAccountFolder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ParentID   string `json:"parent_id"`
		ParentPath string `json:"parent_path"`
		Name       string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[API] 正在创建文件夹: AccountID=%s, Name=%s, ParentPath=%s", id, req.Name, req.ParentPath)

	var account db.Account
	if err := db.DB.First(&account, id).Error; err != nil {
		log.Printf("[API] 创建文件夹失败: 账号未找到 (ID: %s)", id)
		c.PureJSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	driver := core.GetDriver(&account)
	if driver == nil {
		log.Printf("[API] 创建文件夹失败: 驱动加载失败 (Platform: %s)", account.Platform)
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": "driver not found"})
		return
	}

	ctx := c.Request.Context()
	newFolder, err := driver.CreateFolder(ctx, req.ParentID, req.Name)
	if err != nil {
		log.Printf("[API] 创建文件夹异常: %v", err)
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	childPath := req.ParentPath
	if childPath == "/" || childPath == "" {
		childPath = "/" + newFolder.Name
	} else {
		childPath = childPath + "/" + newFolder.Name
	}

	log.Printf("[API] 创建文件夹完成: Path=%s", childPath)
	c.PureJSON(http.StatusOK, FolderItem{
		ID:     newFolder.ID,
		Path:   childPath,
		Label:  newFolder.Name,
		IsLeaf: false,
	})
}
