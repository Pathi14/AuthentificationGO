package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pathi14/AuthentificationGO/internal" // Assurez-vous que le chemin est correct
)

func TestHealth(t *testing.T) {
	r := gin.Default()
	r.GET("/health", internal.Health)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Attendu : %d, Reçu : %d, Détails : %s", http.StatusOK, w.Code, w.Body.String())
	}

	expectedResponse := `{"message":"API is running","status":"OK"}`
	if w.Body.String() != expectedResponse {
		t.Errorf("Attendu : %s, Reçu : %s", expectedResponse, w.Body.String())
	}
}
