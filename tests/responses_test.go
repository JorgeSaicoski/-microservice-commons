package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JorgeSaicoski/microservice-commons/responses"
	"github.com/gin-gonic/gin"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"key": "value"}
	responses.Success(c, "Test message", testData)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response responses.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %s", response.Message)
	}

	if response.Data == nil {
		t.Error("Expected data to be present")
	}

	if response.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestErrorResponses(t *testing.T) {
	tests := []struct {
		name           string
		errorFunc      func(*gin.Context, string)
		expectedStatus int
		expectedCode   string
	}{
		{"BadRequest", responses.BadRequest, http.StatusBadRequest, responses.ErrCodeBadRequest},
		{"Unauthorized", responses.Unauthorized, http.StatusUnauthorized, responses.ErrCodeUnauthorized},
		{"Forbidden", responses.Forbidden, http.StatusForbidden, responses.ErrCodeForbidden},
		{"NotFound", responses.NotFound, http.StatusNotFound, responses.ErrCodeNotFound},
		{"InternalError", responses.InternalError, http.StatusInternalServerError, responses.ErrCodeInternalError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.errorFunc(c, "Test error")

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response responses.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Code != tt.expectedCode {
				t.Errorf("Expected code %s, got %s", tt.expectedCode, response.Code)
			}

			if response.Error != "Test error" {
				t.Errorf("Expected error 'Test error', got %s", response.Error)
			}
		})
	}
}

func TestPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := []string{"item1", "item2", "item3"}
	responses.Paginated(c, testData, 10, 1, 3)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response responses.PaginationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Total != 10 {
		t.Errorf("Expected total 10, got %d", response.Total)
	}

	if response.Page != 1 {
		t.Errorf("Expected page 1, got %d", response.Page)
	}

	if response.PageSize != 3 {
		t.Errorf("Expected page size 3, got %d", response.PageSize)
	}

	if response.TotalPages != 4 {
		t.Errorf("Expected total pages 4, got %d", response.TotalPages)
	}

	if !response.HasNext {
		t.Error("Expected HasNext to be true")
	}

	if response.HasPrev {
		t.Error("Expected HasPrev to be false")
	}
}
