package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreatePersonsHandler(t *testing.T) {
	mongoDbStorage := NewTestMongoDbStorage()
	handler := NewHandler(mongoDbStorage)
	r := gin.New()
	r.POST("/api/v1/persons/", handler.CreatePerson)

	testBody := `{"name":"testName", "address":"testAddress", "work":"testWork", "age":18}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/persons/", bytes.NewBufferString(testBody))

	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Unexpected status code, expected %d, got %d instead", 201, w.Code)
	}
}

func TestGetPersonsHandler(t *testing.T) {
	mongoDbStorage := NewTestMongoDbStorage()
	handler := NewHandler(mongoDbStorage)
	r := gin.New()
	r.GET("/api/v1/persons/", handler.GetPersons)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/persons/", nil)

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Unexpected status code, expected %d, got %d instead", 200, w.Code)
	}
}

func TestUpdatePersonsHandler(t *testing.T) {
	mongoDbStorage := NewTestMongoDbStorage()
	handler := NewHandler(mongoDbStorage)
	r := gin.New()
	r.PATCH("/api/v1/persons/:id", handler.UpdatePerson)

	testBody := `{"name":"testUpdatedName", "address":"testUpdatedAddress", "work":"testUpdatedWork", "age":19}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/api/v1/persons/0", bytes.NewBufferString(testBody))

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Unexpected status code, expected %d, got %d instead", 200, w.Code)
	}
}

func TestGetPersonHandler(t *testing.T) {
	mongoDbStorage := NewTestMongoDbStorage()
	handler := NewHandler(mongoDbStorage)
	r := gin.New()
	r.GET("/api/v1/persons/:id", handler.GetPerson)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/persons/0", nil)

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Unexpected status code, expected %d, got %d instead", 200, w.Code)
	}
}

func TestDeletePersonHandler(t *testing.T) {
	mongoDbStorage := NewTestMongoDbStorage()
	handler := NewHandler(mongoDbStorage)
	r := gin.New()
	r.DELETE("/api/v1/persons/:id", handler.DeletePerson)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/v1/persons/0", nil)

	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("Unexpected status code, expected %d, got %d instead", 204, w.Code)
	}
}
