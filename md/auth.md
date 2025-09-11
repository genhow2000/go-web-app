# 安全認證系統

## 概述

基於 JWT (JSON Web Token) 的現代化認證系統，提供安全的用戶認證和會話管理功能。

## 🚀 功能特色

- **JWT 認證**：基於 JWT 的無狀態認證機制
- **密碼加密**：使用 bcrypt 進行密碼加密存儲
- **角色權限**：基於角色的訪問控制 (RBAC)
- **會話管理**：安全的會話管理和過期控制

## 🔐 認證流程

### 1. 用戶登入

```json
POST /login
{
    "email": "user@example.com",
    "password": "password123"
}
```

### 2. 密碼驗證

```go
// 使用 bcrypt 驗證密碼
err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
```

### 3. 生成 JWT Token

```go
// 創建 JWT claims
claims := jwt.MapClaims{
    "user_id": user.ID,
    "email": user.Email,
    "role": user.Role,
    "exp": time.Now().Add(time.Hour * 24).Unix(),
}

// 生成 token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(secretKey))
```

## 🛡️ 安全特性

- **密碼安全**：使用 bcrypt 算法加密存儲密碼
- **JWT 安全**：使用 HMAC-SHA256 簽名
- **中間件保護**：認證和權限中間件
- **會話管理**：安全的會話管理

## 👥 角色權限

- **admin**：系統管理員，擁有所有權限
- **merchant**：商戶用戶，擁有商戶相關權限
- **customer**：一般用戶，基本權限

## 🔧 配置設定

```go
type JWTConfig struct {
    SecretKey      string        `json:"secret_key"`
    ExpirationTime time.Duration `json:"expiration_time"`
    Issuer         string        `json:"issuer"`
}
```

## 🚀 未來擴展

- 支援 OAuth 2.0 登入
- 支援多因素認證 (MFA)
- 支援單點登入 (SSO)
- 支援 API 金鑰認證

---

**這個認證系統展現了現代 Web 應用的安全設計理念和 JWT 認證的最佳實踐！**
