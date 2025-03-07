package must

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"go-tour/internal/serializers"
	"time"
)

func ParseToken(tokenString string, publicKey string) (*serializers.UserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing")
		}

		key, err := base64.StdEncoding.DecodeString(publicKey)
		if err != nil {
			return nil, errors.Wrap(err, "DecodeString")
		}

		pk, err := jwt.ParseRSAPublicKeyFromPEM(key)
		if err != nil {
			return nil, errors.Wrap(err, "ParseRSAPrivateKeyFromPEM")
		}

		return pk, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if customerID, ok1 := claims["id"]; ok1 {
			if strID, ok2 := customerID.(string); ok2 && len(strID) > 0 {
				return &serializers.UserInfo{
					ID: strID,
				}, nil
			}
		}

		return nil, errors.New("User not found")
	}

	return nil, err
}

func CreateNewWithClaims(data *serializers.UserInfo, secretKey string, expire time.Time) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"id":    fmt.Sprintf("%d", data.ID),
		"email": data.Email,
		"exp":   expire.Unix(),
	})

	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", errors.Wrap(err, "DecodeString")
	}

	privateKey, _ := pem.Decode(key)
	k, err := x509.ParsePKCS1PrivateKey(privateKey.Bytes)
	if err != nil {
		return "", errors.Wrap(err, "ParsePKCS1PrivateKey")
	}

	if err := k.Validate(); err != nil {
		return "", errors.Wrap(err, "Validate")
	}

	token, err := t.SignedString(k)
	if err != nil {
		return "", err
	}

	return token, nil
}
