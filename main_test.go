package main

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks_models "github.com/vinhut/posted/mocks_models"
	mocks_services "github.com/vinhut/posted/mocks_services"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_post := mocks_models.NewMockPostDatabase(ctrl)
	mock_auth := mocks_services.NewMockAuthService(ctrl)

	router := setupRouter(mock_post, mock_auth)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", SERVICE_NAME+"/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}
