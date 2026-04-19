package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zcq/clouddrive-auto-save/internal/core/worker"
	"github.com/zcq/clouddrive-auto-save/internal/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	testDB.AutoMigrate(&db.Account{}, &db.Task{})
	db.DB = testDB // 设置全局 DB 供 handler 使用

	wm := worker.NewManager(1, testDB)
	r := InitRouter(wm)
	return r, testDB
}

func TestListAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, testDB := setupTestRouter(t)

	// 插入测试数据
	testDB.Create(&db.Account{Platform: "139", Nickname: "User1"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/accounts", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var accounts []db.Account
	err := json.Unmarshal(w.Body.Bytes(), &accounts)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expected 1 account, got %d", len(accounts))
	}
	if accounts[0].Nickname != "User1" {
		t.Errorf("expected nickname User1, got %s", accounts[0].Nickname)
	}
}
