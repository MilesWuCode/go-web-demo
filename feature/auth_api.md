# Auth API

## 使用

- **golang-jwt/jwt** 授權token
- **crypto/bcrypt** 密碼加密
- **markbates/goth** 第三方登入
- **go-playground/validator** 表單傳值的驗證套件

## 目標

- 製作註冊
  - /api/auth/register
  - 格式
    - method post
    - Content-Type: application/x-www-form-urlencoded 或 multipart/form-data
    - body (表單欄位)
      - name
      - email
      - password
  - 驗證傳入的表單
  - 資料庫中檢查是否曾經註冊
  - 在資料庫中建立資料

- 製作登入
  - /api/auth/login
  - 格式
    - method post
    - Content-Type: application/x-www-form-urlencoded 或 multipart/form-data
    - body (表單欄位)
      - email
      - password
  - 驗證傳入的表單
  - 資料庫中檢查是否曾經註冊
  - 比對身份
  - 授權

- 製作登出
- 格式
  - /api/auth/logout
    - method post
    - Content-Type: application/x-www-form-urlencoded 或 multipart/form-data
    - body (表單欄位)
      - token
- 驗證傳入的表單
- 清除授權

- 製作第三方登入
  - Google登入
