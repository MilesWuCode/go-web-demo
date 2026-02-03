## 技術

- **Go**
- **GORM**

## 規範

- 使用中文回答問題
- 使用中文撰寫註解與文件

## API格式

- 使用 RESTful 風格
- 使用 HTTP 狀態碼表示請求結果
- 使用 JSON 格式回傳資料
- 單一資料直接回傳資料

  ```json
  {
    "id": 1,
    "name": "apple",
    "email": "apple@email.com"
  }
  ```

- 多筆資料使用 data 陣列包覆，pagination 為分頁資訊

  ```json
  {
    "data": [
      {
        "id": 1,
        "name": "apple",
        "email": "apple@email.com"
      },
      {
        "id": 2,
        "name": "banana",
        "email": "banana@email.com"
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_records": 2,
      "total_pages": 1
    }
  }
  ```
