package auth

import (
	"create-gin-app/pkg/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
)

var (
	jwtSecret []byte
	//EffectiveDuration token 有效期
	EffectiveDuration int
	issuer            string
	// BlackList  权限黑名单, 把token或者用户加入黑名单
	BlackList *cache.Cache
)

// Claims jwt包含信息
type Claims struct {
	// 用户名
	Username string `json:"username"`
	// Role 表示用户的权限
	RoleName string `json:"rolename"`
	jwt.StandardClaims
}

func init() {
	secret := config.DefaultConfig.Section("auth").Key("secret").MustString("create.gin.app.com")
	EffectiveDuration = config.DefaultConfig.Section("auth").Key("effective_duration").MustInt(4)
	issuer = config.DefaultConfig.Section("auth").Key("effective_duration").MustString("baimu")
	jwtSecret = []byte(secret)

	BlackList = cache.New(time.Duration(EffectiveDuration)*time.Hour, 10*time.Minute)
	// 从黑名单库,加载黑名单
	BlackList.LoadFile("blackList.db")
}

//GenerateToken  生成jwt token
func GenerateToken(username, rolename string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(EffectiveDuration) * time.Hour)

	claims := Claims{
		username,
		rolename,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析jwt token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
