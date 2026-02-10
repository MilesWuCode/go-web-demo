package api

import (
	"context"
	"fmt"
	"log"
	"net/http" // 新增導入
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
		config.WithRegion("us-east-1"),
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

	_, err = h.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket: aws.String(h.S3BucketName),
		Key:    aws.String(newFileName),
		Body:   file,
	})

	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法獲取上傳檔案: %w", err), http.StatusBadRequest)
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]string{
		"message":   "檔案上傳到 S3 成功",
		"file_name": newFileName,
	})
}
