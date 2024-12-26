package integration

import (
	"bytes"
	"encoding/json"
	"favorites/internal/db"
	"favorites/internal/handlers"
	"favorites/internal/models/favorite"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
)

var testDB *sqlx.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	var err error
	testDB, err = sqlx.Connect("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		panic(err)
	}
	db.RunMigrations(testDB)
	router = gin.Default()
	handlers.RegisterRoutes(testDB, router)
	code := m.Run()
	_ = testDB.Close()
	os.Exit(code)
}

func clearDB() {
	_, err := testDB.Exec("TRUNCATE TABLE favorites RESTART IDENTITY CASCADE")
	if err != nil {
		return
	}
}

func TestGetFavorites(t *testing.T) {
	clearDB()
	_, err := testDB.Exec(`
		INSERT INTO favorites (id, project_id, owner_type, owner_id, object_id, object_type, created_at)
		VALUES
			(gen_random_uuid(), gen_random_uuid(), 'USER', gen_random_uuid(), gen_random_uuid(), 'IMAGE', NOW())
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	var favorites []favorite.Favorite
	err = json.Unmarshal(w.Body.Bytes(), &favorites)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if len(favorites) != 1 {
		t.Errorf("Expected 1 favorite, got %d", len(favorites))
	}
}

func TestCreateFavorite(t *testing.T) {
	clearDB()
	requestBody := map[string]any{
		"project_id":  uuid.New(),
		"owner_type":  "USER",
		"owner_id":    uuid.New(),
		"object_id":   uuid.New(),
		"object_type": "IMAGE",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		t.Errorf("Error message: %s", w.Body)
	}
	var count int
	err := testDB.Get(&count, "SELECT COUNT(*) FROM favorites")
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 favorite in database, got %d", count)
	}
}

func TestDeleteFavorite(t *testing.T) {
	clearDB()
	var favoriteID uuid.UUID
	err := testDB.QueryRowx(`
		INSERT INTO favorites (id, project_id, owner_type, owner_id, object_id, object_type, created_at)
		VALUES
			(gen_random_uuid(), gen_random_uuid(), 'USER', gen_random_uuid(), gen_random_uuid(), 'IMAGE', NOW())
		RETURNING id
	`).Scan(&favoriteID)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
	req := httptest.NewRequest(http.MethodDelete, "/favorites/"+favoriteID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
	var count int
	err = testDB.Get(&count, "SELECT COUNT(*) FROM favorites")
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 favorites in database, got %d", count)
	}
}
