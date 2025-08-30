package network

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenClaim interface {
	Valid() error
}

type Claim struct {
	UserId     int64  `json:"id,omitempty"`
	FirstName  string `json:"fname,omitempty"`
	OtherNames string `json:"onames,omitempty"`
	Role       string `json:"role,omitempty"`
	IssuedAt   int64  `json:"iat,omitempty"`
	ExpiresAt  int64  `json:"eat,omitempty"`
}

func (t *Claim) Valid() error {
	now := time.Now().Unix()
	dif := t.ExpiresAt - now
	if dif < 0 {
		return fmt.Errorf("expired token")
	}
	if t.Role == "" {
		return fmt.Errorf("invalid role")
	}
	if now < t.IssuedAt {
		return fmt.Errorf("token issued in the future")
	}

	return nil
}

func GenerateClaim(
	id int64,
	first_name string,
	other_names string,
	role string,
	expire_seconds int64,
) (TokenClaim, string, error) {
	now := time.Now().Unix()
	exp := now + expire_seconds

	c := &Claim{UserId: id, FirstName: first_name, OtherNames: other_names, Role: role, IssuedAt: now, ExpiresAt: int64(exp)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
  
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, "", err
	}
	return c, tokenString, nil
}
