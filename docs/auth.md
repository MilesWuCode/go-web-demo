# 授權功能

## 使用

- **golang-jwt/jwt** 授權token
- **crypto/bcrypt** 密碼加密
- **go-playground/validator** 表單傳值的驗證套件

## 目標

- 規格
  - 資料庫 users 資料表
    - 帳號：email
    - 密碼：password
    - 一對多關聯表 refresh_tokens
  - 資料庫 refresh_tokens 資料表
    - user_id
    - token
    - expired_at

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
  - 生成JWT授權
    - access_token 有效期限1小時
    - refresh_token 有效期限7天

- 製作取得個人資料
  - /api/auth/me
  - 格式
    - method get
    - header Authorization: Bearer access_token
  - 中間層驗證JWT功能

- 製作RefreshToken
  - /api/auth/refresh-token
  - 格式
    - method post
    - header Authorization: Bearer refresh_token
  - 查詢資料庫 refresh_tokens 資料表 是否存在該refresh_token
  - 生成JWT授權
    - access_token 有效期限1小時
    - refresh_token 有效期限7天

- 製作登出
- 格式
  - /api/auth/logout
    - method post
    - Content-Type: application/x-www-form-urlencoded 或 multipart/form-data
    - body (表單欄位)
      - token
  - 中間層驗證JWT功能
  - 驗證傳入的表單
  - 清除授權

- 製作第三方登入
  - Google登入，不需要用戶頭像
  - 建立SocialProvider資料表來存GoogleID，未來有其他登入也可使用
  - 當前端登入完成後，取得Google提供的身份憑證
    - 驗證身份憑證
    - 若無帳號建立帳號
    - 生成JWT授權
      - access_token 有效期限1小時
      - refresh_token 有效期限7天
