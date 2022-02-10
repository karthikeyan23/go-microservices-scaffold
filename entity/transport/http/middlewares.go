package http

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"net/http"
	"strings"
)

var (
	ErrKIDNotFound            = errors.New("kid header not found")
	ErrUnableToParsePublicKey = errors.New("could not parse public key")
	ErrUnexpectedTokenVersion = errors.New("unexpected token version")
)

func genericMiddlewareToSetHTTPHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func jwtMiddlewareForMicrosoftIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := getJWTTokenFromHTTPHeader(w, r)
		if err != nil {
			return
		}

		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getJWKPublicKeyForMicrosoftIdentity(r.Context(), token)
		})
		//Check if token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "props", claims)
			//TO-DO: Add claims validation and adding user to database if not exists
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte(""))
			if err != nil {
				return
			}
		}
	})
}

func getJWTTokenFromHTTPHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(ErrEmptyAuthHeader.Error()))
		if err != nil {
			return "", ErrEmptyAuthHeader
		}
		return "", ErrEmptyAuthHeader
	}
	authHeader := strings.Split(auth, "Bearer ")
	if len(authHeader) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(ErrMalformedToken.Error()))
		if err != nil {
			return "", ErrMalformedToken
		}
		return "", ErrMalformedToken
	}
	jwtToken := authHeader[1]
	return jwtToken, nil
}

func getJWKPublicKeyForMicrosoftIdentity(ctx context.Context, token *jwt.Token) (interface{}, error) {
	//Get the latest Public JWK keys from Microsoft Identity Platform
	jwtKeyURl, err := getJwkKeyURL(token)
	if err != nil {
		return nil, err
	}
	keySet, err := jwk.Fetch(ctx, jwtKeyURl)
	if err != nil {
		return nil, err
	}
	//Check if kid exists in the token
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, ErrKIDNotFound
	}
	//Get the keys based on kid
	keys, ok := keySet.LookupKeyID(kid)
	if !ok {
		return nil, fmt.Errorf("key %v not found", kid)
	}
	//Get the public key from the key
	publicKey := &rsa.PublicKey{}
	err = keys.Raw(publicKey)
	if err != nil {
		return nil, ErrUnableToParsePublicKey
	}
	return publicKey, nil
}

func getJwkKeyURL(token *jwt.Token) (string, error) {
	tokenVersion := token.Claims.(jwt.MapClaims)["ver"].(string)
	if tokenVersion == "1.0" {
		return "https://login.microsoftonline.com/common/discovery/keys", nil
	} else if tokenVersion == "2.0" {
		return "https://login.microsoftonline.com/common/discovery/v2.0/keys", nil
	} else {
		return "", ErrUnexpectedTokenVersion
	}
}
