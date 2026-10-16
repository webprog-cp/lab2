package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	application "cp_lab1/internal/app"
	"cp_lab1/internal/config"
	"cp_lab1/internal/handlers"
	"cp_lab1/internal/models"
	"cp_lab1/internal/repositories"

	"github.com/gin-gonic/gin"
)

func TestRequireAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/protected", requireAPIKey, func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/protected", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status without key = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-API-Key", apiKey)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status with key = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestRequestDTOValidation(t *testing.T) {
	catRepo := repositories.NewCatRepoMemory()
	shelterRepo := repositories.NewShelterRepoMemory()
	if _, err := shelterRepo.Create(models.Shelter{Name: "Test", City: "Kyiv", Address: "Test", Capacity: 1}); err != nil {
		t.Fatal(err)
	}
	a := &application.App{
		CatHandler:     handlers.NewCatHandler(catRepo, shelterRepo),
		ShelterHandler: handlers.NewShelterHandler(shelterRepo, catRepo),
	}
	r, err := New(&config.Config{GinMode: gin.TestMode, APIPrefix: "/api", APIVersion: "v1"}, a)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		body string
		want int
	}{
		{
			name: "cat accepts zero months",
			path: "/api/v1/cats",
			body: `{"name":"Kitten","ageMonths":0,"gender":"female","adopted":false,"shelterID":1}`,
			want: http.StatusCreated,
		},
		{
			name: "cat rejects client ID",
			path: "/api/v1/cats",
			body: `{"id":42,"name":"Kitten","ageMonths":0,"gender":"female","adopted":false,"shelterID":1}`,
			want: http.StatusBadRequest,
		},
		{
			name: "cat rejects excessive age",
			path: "/api/v1/cats",
			body: `{"name":"Old cat","ageMonths":361,"gender":"male","adopted":false,"shelterID":1}`,
			want: http.StatusBadRequest,
		},
		{
			name: "cat rejects invalid gender",
			path: "/api/v1/cats",
			body: `{"name":"Cat","ageMonths":12,"gender":"unknown","adopted":false,"shelterID":1}`,
			want: http.StatusBadRequest,
		},
		{
			name: "shelter rejects zero capacity",
			path: "/api/v1/shelters",
			body: `{"name":"Shelter","city":"Kyiv","address":"Main Street","capacity":0}`,
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-API-Key", apiKey)
			r.ServeHTTP(w, req)
			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d; body = %s", w.Code, tt.want, w.Body.String())
			}
			if tt.want == http.StatusBadRequest && !strings.Contains(w.Body.String(), "invalid request body") {
				t.Fatalf("response exposed validation details: %s", w.Body.String())
			}
		})
	}
}
