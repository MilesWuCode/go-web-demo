package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-demo/config"
	"web-demo/model"
	"web-demo/server"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestApp 是一個輔助函式，用於設定測試環境，包括一個記憶體中的 SQLite 資料庫。
func newTestApp() (*server.Application, error) {
	// 使用記憶體中的 SQLite 資料庫進行測試，以避免依賴外部資料庫。
	// "cache=shared" 是必要的，以確保在測試執行期間資料庫連接保持活動狀態。
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自動遷移所有相關的資料庫模型。
	err = db.AutoMigrate(&model.User{}, &model.SocialProvider{}, &model.RefreshToken{}, &model.Media{}, &model.Post{})
	if err != nil {
		return nil, err
	}

	// 建立一個帶有測試配置和資料庫連接的 Application 實例。
	app := &server.Application{
		Config: &config.AppConfig{
			// 在此處可以為測試添加特定的配置值。
		},
		DB: db,
	}
	return app, nil
}

// TestGetAllUsers 包含了對 GetAllUsers 處理器的多個測試案例。
func TestGetAllUsers(t *testing.T) {
	// --- 設定 ---
	app, err := newTestApp()
	if err != nil {
		t.Fatalf("無法建立測試應用程式: %v", err)
	}
	userHandler := NewUserHandler(app)

	// 清理資料庫以確保測試的獨立性。
	app.DB.Exec("DELETE FROM users")

	// --- 填充資料庫 ---
	// 建立一些測試使用者。
	users := []model.User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
		{Name: "David", Email: "david@example.com"},
		{Name: "Eve", Email: "eve@example.com"},
	}
	if err := app.DB.Create(&users).Error; err != nil {
		t.Fatalf("無法填充測試使用者: %v", err)
	}

	// --- 測試案例 ---

	t.Run("預設分頁應回傳第一頁的結果", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users", nil)
		rr := httptest.NewRecorder()

		userHandler.GetAllUsers(rr, req)

		// 驗證狀態碼
		if rr.Code != http.StatusOK {
			t.Errorf("預期狀態碼 %d; 得到 %d", http.StatusOK, rr.Code)
		}

		var resp PaginatedUserResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("無法解碼回應: %v", err)
		}

		// 驗證資料
		// 預設頁面大小為 10，所以應該回傳所有 5 位使用者。
		if len(resp.Data) != 5 {
			t.Errorf("預期 5 位使用者; 得到 %d", len(resp.Data))
		}
		if resp.Pagination.TotalRecords != 5 {
			t.Errorf("預期總記錄數為 5; 得到 %d", resp.Pagination.TotalRecords)
		}
		if resp.Data[0].Name != "Alice" {
			t.Errorf("預期第一位使用者為 Alice; 得到 %s", resp.Data[0].Name)
		}
	})

	t.Run("指定頁面大小和頁碼", func(t *testing.T) {
		// 請求第二頁，每頁 2 筆資料。
		req := httptest.NewRequest("GET", "/api/users?page=2&pageSize=2", nil)
		rr := httptest.NewRecorder()

		userHandler.GetAllUsers(rr, req)

		// 驗證狀態碼
		if rr.Code != http.StatusOK {
			t.Errorf("預期狀態碼 %d; 得到 %d", http.StatusOK, rr.Code)
		}

		var resp PaginatedUserResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("無法解碼回應: %v", err)
		}

		// 驗證資料
		// 第二頁應該有 2 筆資料 (Charlie, David)。
		if len(resp.Data) != 2 {
			t.Errorf("預期 2 位使用者; 得到 %d", len(resp.Data))
		}
		if resp.Pagination.TotalRecords != 5 {
			t.Errorf("預期總記錄數為 5; 得到 %d", resp.Pagination.TotalRecords)
		}
		if resp.Pagination.CurrentPage != 2 {
			t.Errorf("預期目前頁碼為 2; 得到 %d", resp.Pagination.CurrentPage)
		}
		if resp.Data[0].Name != "Charlie" {
			t.Errorf("預期第二頁的第一位使用者為 Charlie; 得到 %s", resp.Data[0].Name)
		}
	})

	t.Run("當資料庫為空時應回傳空列表", func(t *testing.T) {
		// 清空使用者資料表
		app.DB.Exec("DELETE FROM users")

		req := httptest.NewRequest("GET", "/api/users", nil)
		rr := httptest.NewRecorder()

		userHandler.GetAllUsers(rr, req)

		// 驗證狀態碼
		if rr.Code != http.StatusOK {
			t.Errorf("預期狀態碼 %d; 得到 %d", http.StatusOK, rr.Code)
		}

		var resp PaginatedUserResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("無法解碼回應: %v", err)
		}

		// 驗證資料
		if len(resp.Data) != 0 {
			t.Errorf("預期 0 位使用者; 得到 %d", len(resp.Data))
		}
		if resp.Pagination.TotalRecords != 0 {
			t.Errorf("預期總記錄數為 0; 得到 %d", resp.Pagination.TotalRecords)
		}
	})
}
