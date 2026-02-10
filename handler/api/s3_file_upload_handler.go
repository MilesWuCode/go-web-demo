package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"web-demo/server"
	"web-demo/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Println("無法載入 AWS 設定:", err) // 使用 app.Logger
		return nil
	}

	// 從 AppConfig 中獲取 S3 儲存桶名稱
	s3BucketName := app.Config.S3BucketName
	if s3BucketName == "" {
		log.Println("AppConfig.S3BucketName 未設定") // 使用 app.Logger
		return nil
	}

	return &S3FileUploadHandler{
		App:          app,
		S3Client:     s3.NewFromConfig(cfg),
		S3BucketName: s3BucketName,
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

	// 將檔案上傳到 S3
	_, err = h.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(h.S3BucketName),
		Key:    aws.String(newFileName),
		Body:   file,
		ACL:    "public-read", // 設置為公開讀取，根據需求調整
	})
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法上傳檔案到 S3: %w", err), http.StatusInternalServerError)
		return
	}

	var s3ObjectURL string
	if h.App.Config.S3Endpoint != "" {
		// 如果有自訂 S3 Endpoint，則建構自訂 URL
		s3ObjectURL = fmt.Sprintf("%s/%s/%s",
			h.App.Config.S3Endpoint,
			h.S3BucketName,
			newFileName)
	} else {
		// 否則，建構標準 AWS S3 URL
		s3ObjectURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
			h.S3BucketName,
			os.Getenv("AWS_REGION"), // 假設 AWS_REGION 環境變數已設定
			newFileName)
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]string{
		"message":       "檔案上傳到 S3 成功",
		"file_name":     newFileName,
		"s3_object_url": s3ObjectURL,
	})
}
