package ctxdata

import "github.com/golang-jwt/jwt"

const Identify = "mingri"

// GetJwtToken 生成JWT令牌
// 参数:
//
//	secretKey: 用于签名的密钥
//	iat: 令牌签发时间戳
//	seconds: 令牌有效期（秒）
//	uid: 用户唯一标识
//
// 返回值:
//
//	string: 生成的JWT令牌字符串
//	error: 生成过程中可能发生的错误
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	// 创建JWT声明集合
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identify] = uid

	// 创建新的JWT令牌并设置声明
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	// 使用密钥对令牌进行签名并返回
	return token.SignedString([]byte(secretKey))
}
