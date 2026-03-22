package middleware

import (
	"time"

	"github.com/kataras/iris/v12"
	irsjwt "github.com/kataras/iris/v12/middleware/jwt"
)

// Claims là payload được nhúng trong JWT token.
type Claims struct {
	UserID int32  `json:"user_id"`
	Email  string `json:"email"`
}

// JWTVerify trả về middleware xác thực Bearer token.
// Đặt middleware này trước các route cần bảo vệ.
func JWTVerify(secret string) iris.Handler {
	verifier := irsjwt.NewVerifier(irsjwt.HS256, []byte(secret))
	return verifier.Verify(func() any {
		return new(Claims)
	})
}

// GetClaims lấy Claims từ context (phải gọi sau JWTVerify).
func GetClaims(ctx iris.Context) *Claims {
	c, ok := irsjwt.Get(ctx).(*Claims)
	if !ok {
		return nil
	}
	return c
}

// GenerateToken tạo JWT token từ claims cho trước.
func GenerateToken(secret string, expiry time.Duration, claims Claims) (string, error) {
	signer := irsjwt.NewSigner(irsjwt.HS256, []byte(secret), expiry)
	token, err := signer.Sign(claims)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
