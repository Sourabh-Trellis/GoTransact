package utils

import (
	"GoTransact/apps/accounts/models"
	"time"

	"aidanwoods.dev/go-paseto"
	// "github.com/aead/chacha20poly1305"
	// "github.com/o1egl/paseto"
)

var (
	secretKey = paseto.NewV4AsymmetricSecretKey() // don't share this!!!
	publicKey = secretKey.Public()                // DO share this one
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
