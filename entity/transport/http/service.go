package http

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
	"go_scafold/entity/transport"
	"net/http"
	"strings"
)

func NewHTTPServer(ctx context.Context, endpoints transport.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(httpMiddleware)
	r.Use(jwtMiddleware)

	r.Methods("GET").Path("/entity").Handler(kithttp.NewServer(
		endpoints.GetEntity,
		decodeGetEntityRequest,
		encodeResponse))

	return r
}

func httpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Authorization header is empty"))
			if err != nil {
				return
			}
			return
		}
		authHeader := strings.Split(auth, "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Malformed Token"))
			if err != nil {
				return
			}
		} else {
			jwtToken := authHeader[1]

			//Get the latest Public JWK keys from Azure
			keySet, err := jwk.Fetch(r.Context(), "https://login.microsoftonline.com/common/discovery/keys")

			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				//Check if kid exists in the token
				kid, ok := token.Header["kid"].(string)
				if !ok {
					return nil, fmt.Errorf("kid header not found")
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
					return nil, fmt.Errorf("could not parse pubkey")
				}

				return publicKey, nil
			})
			//Check if token is valid
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "props", claims)
				//TO-DO: Add claims validation and adding user to database if not exists
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("Unauthorized"))
				if err != nil {
					return
				}
			}
		}
	})
}

func encodeResponse(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
	return json.NewEncoder(writer).Encode(i)
}

func decodeGetEntityRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req transport.GetEntityByIDRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
