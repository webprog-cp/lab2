package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cp_lab1/internal/const_errors"
	"cp_lab1/internal/models"
	"cp_lab1/internal/repositories"

	"github.com/gin-gonic/gin"
)

type failingCatRepo struct {
	repositories.CatRepo
}

func (failingCatRepo) GetAll() ([]models.Cat, error) {
	return nil, errors.New("sensitive database details")
}

func TestCatHandlerReturnsErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCatHandler(repositories.NewCatRepoMemory(), repositories.NewShelterRepoMemory())
	r.GET("/cats/:id", h.GetByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cats/nope", nil)
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

func TestRepositoryErrorsBecomeSafeResponses(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
		body string
	}{
		{name: "not found", err: const_errors.NoEntityByID, want: http.StatusNotFound, body: "resource not found"},
		{name: "invalid reference", err: const_errors.InvalidReference, want: http.StatusBadRequest, body: "related resource not found"},
		{name: "has dependents", err: const_errors.EntityHasDependents, want: http.StatusConflict, body: "resource has dependent records"},
		{name: "shelter full", err: const_errors.ShelterAtCapacity, want: http.StatusConflict, body: "shelter is at full capacity"},
		{name: "capacity below usage", err: const_errors.CapacityBelowUsage, want: http.StatusConflict, body: "capacity is below current occupancy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/test", func(c *gin.Context) { respondRepositoryError(c, tt.err) })

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

			if w.Code != tt.want || !strings.Contains(w.Body.String(), tt.body) {
				t.Fatalf("status = %d, body = %s; want %d and %q", w.Code, w.Body.String(), tt.want, tt.body)
			}
		})
	}
}

func TestCatCreateRejectsFullShelter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	catRepo := repositories.NewCatRepoMemory()
	shelterRepo := repositories.NewShelterRepoMemory()
	if _, err := shelterRepo.Create(models.Shelter{Name: "Test", Capacity: 1}); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	h := NewCatHandler(catRepo, shelterRepo)
	r.POST("/cats", h.Create)
	body := `{"name":"Cat","ageMonths":12,"gender":"male","adopted":false,"shelterID":1}`

	for i, want := range []int{http.StatusCreated, http.StatusConflict} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/cats", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != want {
			t.Fatalf("request %d: status = %d, want %d; body = %s", i+1, w.Code, want, w.Body.String())
		}
	}

	cats, err := catRepo.GetCatsByShelterID(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 1 {
		t.Fatalf("cats in shelter = %d, want 1", len(cats))
	}
}

func TestCatUpdateRejectsMoveToFullShelter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	catRepo := repositories.NewCatRepoMemory()
	shelterRepo := repositories.NewShelterRepoMemory()
	for i := 0; i < 2; i++ {
		if _, err := shelterRepo.Create(models.Shelter{Name: "Test", Capacity: 1}); err != nil {
			t.Fatal(err)
		}
	}
	first, err := catRepo.Create(models.Cat{Name: "First", ShelterID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := catRepo.Create(models.Cat{Name: "Second", ShelterID: 2}); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	h := NewCatHandler(catRepo, shelterRepo)
	r.PUT("/cats/:id", h.Update)
	body := `{"name":"First","ageMonths":12,"gender":"male","adopted":false,"shelterID":2}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/cats/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusConflict, w.Body.String())
	}
	unchanged, err := catRepo.GetByID(first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if unchanged.ShelterID != 1 {
		t.Fatalf("shelter ID = %d, want 1", unchanged.ShelterID)
	}
}

func TestCatHandlerHidesRepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCatHandler(failingCatRepo{}, repositories.NewShelterRepoMemory())
	r.GET("/cats", h.GetAll)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/cats", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	if strings.Contains(w.Body.String(), "sensitive database details") {
		t.Fatalf("response exposed repository error: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "internal server error") {
		t.Fatalf("unexpected response: %s", w.Body.String())
	}
}
