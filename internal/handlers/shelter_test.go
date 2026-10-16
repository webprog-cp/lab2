package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cp_lab1/internal/models"
	"cp_lab1/internal/repositories"

	"github.com/gin-gonic/gin"
)

func TestShelterHandlerReturnsErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewShelterHandler(repositories.NewShelterRepoMemory(), repositories.NewCatRepoMemory())
	r.GET("/shelters/:id", h.GetByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/shelters/nope", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body["error"] == "" {
		t.Fatalf("missing error response: %s", w.Body.String())
	}
}

func TestShelterUpdateRejectsCapacityBelowOccupancy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	shelterRepo := repositories.NewShelterRepoMemory()
	catRepo := repositories.NewCatRepoMemory()
	if _, err := shelterRepo.Create(models.Shelter{Name: "Test", Capacity: 2}); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := catRepo.Create(models.Cat{Name: "Cat", ShelterID: 1}); err != nil {
			t.Fatal(err)
		}
	}

	r := gin.New()
	h := NewShelterHandler(shelterRepo, catRepo)
	r.PUT("/shelters/:id", h.Update)
	body := `{"name":"Test","city":"Kyiv","address":"Street","capacity":1}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/shelters/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusConflict, w.Body.String())
	}
	shelter, err := shelterRepo.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}
	if shelter.Capacity != 2 {
		t.Fatalf("capacity = %d, want 2", shelter.Capacity)
	}
}
