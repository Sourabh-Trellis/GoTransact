package utils

import (
	"GoTransact/apps/accounts/models"
	"context"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-redis/redis/v8"
)

var (
	secretKey = paseto.NewV4AsymmetricSecretKey() // don't share this!!!
	publicKey = secretKey.Public()                // DO share this one
	ctx       = context.Background()
	rdb       = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
)

func CreateToken(user models.User) (string, error) {

	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	token.Set("user", user)

	signedToken := token.V4Sign(secretKey, nil)

	return signedToken, nil
}

func VerifyToken(signedToken string) (any, error) {

	val, err := rdb.Get(ctx, signedToken).Result()
	if err == nil && val == "Blacklisted" {
		return nil, fmt.Errorf("token has been revoked")
	}

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	verifiedtoken, err := parser.ParseV4Public(publicKey, signedToken, nil)
	if err != nil {
		return "", err
	}
	var User models.User
	if err := verifiedtoken.Get("user", &User); err != nil {
		return "", err
	}
	return User, nil
}
