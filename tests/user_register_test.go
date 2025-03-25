package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pathi14/AuthentificationGO/internal/infrastructure/database"
)

func TestRegisterSucess(t *testing.T) {
	payload := map[string]string{
		"name":     "TestUser",
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Attendu : %d, Reçu : %d, Détails : %s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestRegisterInvalidInput(t *testing.T) {
	payload := map[string]string{
		"name":     "InvalidUser",
		"email":    "invalid-email",
		"password": "password123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Attendu : %d, Reçu : %d, Détails : %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	db, err := database.ConnectTestDB()
	if err != nil {
		t.Fatalf("Erreur lors de la connexion à la DB de test : %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE email = 'test@example.com'")
	if err != nil {
		t.Fatalf("Erreur lors de la suppression d'un utilisateur existant : %v", err)
	}

	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES ('ExistingUser', 'test@example.com', 'password123')")
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion d'un utilisateur dans la DB de test : %v", err)
	}

	payload := map[string]string{
		"name":     "DuplicateUser",
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Attendu : %d, Reçu : %d, Détails : %s", http.StatusConflict, w.Code, w.Body.String())
	}

	_, err = db.Exec("DELETE FROM users WHERE email = 'test@example.com'")
	if err != nil {
		t.Errorf("Erreur lors du nettoyage des utilisateurs de test : %v", err)
	}
}
