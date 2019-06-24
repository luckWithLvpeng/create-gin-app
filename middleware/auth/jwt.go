package auth

import (
	"eme/pkg/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	jwtSecret []byte
	//EffectiveDuration token 有效期
	EffectiveDuration int
	issuer            string
	//TokenBlackMap  token 黑名单 , cmap 多线程安全map的一种实现
	TokenBlackMap = cmap.New()
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
	secret := config.DefaultConfig.Section("auth").Key("secret").MustString("eme.alleyes.com")
	EffectiveDuration = config.DefaultConfig.Section("auth").Key("effective_duration").MustInt(4)
	issuer = config.DefaultConfig.Section("auth").Key("effective_duration").MustString("baimu")
	jwtSecret = []byte(secret)

	deleteTokenExpirse()
}

// deleteTokenExpirse 定期删除过期的token
func deleteTokenExpirse() {

	now := time.Now()
	// Insert items to temporary map.
	for item := range TokenBlackMap.IterBuffered() {
		time := item.Val.(time.Time)
		if now.After(time) {
			TokenBlackMap.Remove(item.Key)
		}

	}
	// 每 2 个小时清理一下过期的token
	time.AfterFunc(time.Hour*2, deleteTokenExpirse)
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
