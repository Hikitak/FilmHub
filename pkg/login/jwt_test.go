package login

import "testing"

func TestGenerateAndParseToken(t *testing.T) {
    Init("testsecret")
    tokenStr, err := GenerateToken(42, "admin")
    if err != nil {
        t.Fatalf("generate token: %v", err)
    }
    claims, err := ParseToken(tokenStr)
    if err != nil {
        t.Fatalf("parse token: %v", err)
    }
    if claims.UserID != 42 || claims.Role != "admin" {
        t.Errorf("unexpected claims: %+v", claims)
    }
} 