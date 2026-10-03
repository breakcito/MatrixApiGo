package utils

import (
	"matrix-api-go/src/utils"
	"testing"
)

func TestJWT_UserTokenLifecycle(t *testing.T) {
	token, expiresIn, err := utils.GenerateUserToken("test-user")
	if err != nil {
		t.Fatalf("error generando token de usuario: %v", err)
	}
	if token == "" || expiresIn <= 0 {
		t.Fatalf("token vacío o expiresIn inválido: %s, %d", token, expiresIn)
	}

	claims, err := utils.ValidateUserToken(token)
	if err != nil {
		t.Fatalf("error validando token de usuario: %v", err)
	}
	if claims.Username != "test-user" {
		t.Fatalf("esperaba usuario test-user, obtuve %s", claims.Username)
	}
}

func TestJWT_InvalidToken(t *testing.T) {
	_, err := utils.ValidateUserToken("token-invalido-o-falso")
	if err == nil {
		t.Fatal("se esperaba error al validar token inválido")
	}
}

func TestJWT_ServiceToken(t *testing.T) {
	token, err := utils.GenerateServiceToken()
	if err != nil {
		t.Fatalf("error generando token de servicio: %v", err)
	}
	if token == "" {
		t.Fatal("token de servicio no debe ser vacío")
	}
}
