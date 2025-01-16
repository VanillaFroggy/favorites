package integration

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"favorites/internal/db"
	"favorites/internal/handlers"
	"favorites/internal/models/favorite"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
)

var testDB *sqlx.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	ctx := context.Background()
	envFile, err := godotenv.Read("../../deploy/.env")
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     envFile["DATABASE_USER"],
			"POSTGRES_PASSWORD": envFile["DATABASE_PASSWORD"],
			"POSTGRES_DB":       envFile["DATABASE_NAME"],
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	host, _ := postgresContainer.Host(ctx)
	mappedPort, _ := postgresContainer.MappedPort(ctx, "5432/tcp")
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		envFile["DATABASE_USER"],
		envFile["DATABASE_PASSWORD"],
		host, mappedPort.Port(),
		envFile["DATABASE_NAME"],
	)
	testDB, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	err = db.RunMigrations(testDB, "file://../../internal/db/migrations")
	if err != nil {
		panic(err)
	}
	router = gin.Default()
	handlers.RegisterRoutes(testDB, router)
	code := m.Run()
	_ = testDB.Close()
	_ = postgresContainer.Terminate(ctx)
	os.Exit(code)
}

func clearDB() {
	_, err := testDB.Exec("TRUNCATE TABLE favorites RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}
}

func TestGetFavorites(t *testing.T) {
	clearDB()
	var id uuid.UUID
	var ownerID uuid.UUID
	err := testDB.QueryRowx(`
		INSERT INTO favorites (id, project_id, owner_type, owner_id, object_id, object_type, created_at)
		VALUES (gen_random_uuid(), gen_random_uuid(), 'USER', gen_random_uuid(), gen_random_uuid(), 'IMAGE', NOW())
		RETURNING id, owner_id;
	`).Scan(&id, &ownerID)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
		return
	}
	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	urlQuery := req.URL.Query()
	urlQuery.Add("owner_type", "USER")
	urlQuery.Add("owner_id", ownerID.String())
	urlQuery.Add("limit", "25")
	urlQuery.Add("cursor", "")
	req.URL.RawQuery = urlQuery.Encode()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		t.Errorf("Error message: %s", w.Body)
		return
	}
	cursorBase64 := w.Header()["X-Next-Cursor"][0]
	if cursorBase64 == "" {
		t.Errorf("Error message: %s", "cursor is nil")
		return
	}
	decodedCursor, err := base64.URLEncoding.DecodeString(cursorBase64)
	if err != nil {
		t.Errorf("Error message: %s", w.Body)
		return
	}
	cursorID, err := uuid.Parse(string(decodedCursor))
	if err != nil {
		t.Errorf("Error message: %s", w.Body)
		return
	} else if id != cursorID {
		t.Errorf("Error message: %s", "ids are not equal")
		return
	}
	var favorites []favorite.Favorite
	err = json.Unmarshal(w.Body.Bytes(), &favorites)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
		return
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
		return
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
