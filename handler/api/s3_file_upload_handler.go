package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime" // 新增導入
	"net/http"
	"path/filepath" // 新增導入
	"web-demo/server"
	"web-demo/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3FileUploadHandler 結構用於處理 S3 檔案上傳相關請求
type S3FileUploadHandler struct {
	App          *server.Application
	S3Client     *s3.Client
	S3BucketName string
}

// NewS3FileUploadHandler 建立並回傳一個新的 S3FileUploadHandler 實例
func NewS3FileUploadHandler(app *server.Application) *S3FileUploadHandler {
	// 載入 AWS 預設配置
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"), // 硬編碼區域，如果 RustFS 需要特定區域，或應從 AppConfig 獲取
		config.WithBaseEndpoint(app.Config.S3Endpoint),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(app.Config.S3AccessKeyID, app.Config.S3SecretAccessKey, "")),
	)
	if err != nil {
		log.Println("無法載入 AWS 設定:", err)
		return nil
	}

	// 建立 S3 客戶端，並應用自訂端點解析器和路徑樣式訪問
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// 如果 S3 相容服務需要路徑樣式訪問，請設定為 true
		// 例如：MinIO, RustFS 等
		o.UsePathStyle = true
	})

	return &S3FileUploadHandler{
		App:          app,
		S3Client:     s3Client,
		S3BucketName: app.Config.S3BucketName,
	}
}

// UploadS3File 處理 S3 檔案上傳
func (h *S3FileUploadHandler) UploadS3File(w http.ResponseWriter, r *http.Request) {
	// 限制上傳檔案大小為 10MB
	r.ParseMultipartForm(10 << 20) // 10MB

	file, handler, err := r.FormFile("file") // "file" 是表單中檔案輸入欄位的名稱
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法獲取上傳檔案: %w", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 產生一個混合英文和數字的唯一檔案名稱
	newFileName := util.GenerateEnglishMixedFileName(handler.Filename)

	// 讀取檔案開頭部分以偵測 Content-Type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		h.App.ErrorJSON(w, fmt.Errorf("無法讀取檔案內容以偵測類型: %w", err), http.StatusInternalServerError)
		return
	}
	// 將讀取指針重置回檔案開頭
	_, err = file.Seek(0, 0)
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法重置檔案讀取指針: %w", err), http.StatusInternalServerError)
		return
	}

	// 優先使用 http.DetectContentType 偵測 MIME 類型
	contentType := http.DetectContentType(buffer)
	if contentType == "application/octet-stream" {
		// 如果 http.DetectContentType 仍是通用類型，則嘗試使用副檔名判斷
		extType := mime.TypeByExtension(filepath.Ext(handler.Filename))
		if extType != "" {
			contentType = extType
		}
	}

	// 將檔案上傳到 S3
	_, err = h.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.S3BucketName),
		Key:         aws.String(newFileName),
		Body:        file,
		ACL:         "public-read",           // 重新啟用 ACL: "public-read"
		ContentType: aws.String(contentType), // 設定 Content-Type
	})

	if err != nil {
		// 修正錯誤訊息和狀態碼，並在錯誤後返回
		h.App.ErrorJSON(w, fmt.Errorf("無法上傳檔案到 S3: %w", err), http.StatusInternalServerError)
		return
	}

	// 根據是否使用自訂 S3 Endpoint 建構 S3 物件的公開 URL
	var s3ObjectURL string
	if h.App.Config.S3Endpoint != "" {
		// 如果有自訂 S3 Endpoint，則建構自訂 URL
		s3ObjectURL = fmt.Sprintf("%s/%s/%s",
			h.App.Config.S3Endpoint,
			h.S3BucketName,
			newFileName)
	} else {
		// 否則，建構標準 AWS S3 URL
		// 注意：cfg.Region 來自 NewS3FileUploadHandler 中的 config.WithRegion("us-east-1")
		s3ObjectURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
			h.S3BucketName,
			"us-east-1", // 這裡使用硬編碼區域，因為 cfg.Region 在此作用域不可用
			newFileName)
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]string{
		"message":   "檔案上傳到 S3 成功",
		"file_name": newFileName,
		"file_url":  s3ObjectURL,
	})
}
