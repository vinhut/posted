package main

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks_models "github.com/vinhut/posted/mocks_models"
	mocks_services "github.com/vinhut/posted/mocks_services"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestCheckUser(t *testing.T) {

	now := time.Now()
	token := "852a37a34b727c0e0b331806-7af4bdfdcc60990d427f383efecc8529289d040dd67e0753b9e2ee5a1e938402186f28324df23f6faa4e2bbf43f584ae228c55b00143866215d6e92805d470a1cc2a096dcca4d43527598122313be412e17fbefdcdab2fae02e06a405791d936862d4fba688b3c7fd784d4"
	user_data := "{\"uid\": \"1\", \"email\": \"test@email.com\", \"role\": \"standard\", \"created\": \"" + now.Format("2006-01-02T15:04:05") + "\"}"

	os.Setenv("KEY", "12345678901234567890123456789012")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_auth := mocks_services.NewMockAuthService(ctrl)
	mock_auth.EXPECT().Check(gomock.Any(), gomock.Any()).Return(user_data, nil)

	data, _ := checkUser(mock_auth, token)
	var test_data map[string]interface{}

	if err := json.Unmarshal([]byte(user_data), &test_data); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, test_data, data)
}

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

func TestGetPost(t *testing.T) {

	now := time.Now()
	token := "852a37a34b727c0e0b331806-7af4bdfdcc60990d427f383efecc8529289d040dd67e0753b9e2ee5a1e938402186f28324df23f6faa4e2bbf43f584ae228c55b00143866215d6e92805d470a1cc2a096dcca4d43527598122313be412e17fbefdcdab2fae02e06a405791d936862d4fba688b3c7fd784d4"
	user_data := "{\"uid\": \"1\", \"email\": \"test@email.com\", \"role\": \"standard\", \"created\": \"" + now.Format("2006-01-02T15:04:05") + "\"}"
	postid := "1"

	os.Setenv("KEY", "12345678901234567890123456789012")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_post := mocks_models.NewMockPostDatabase(ctrl)
	mock_auth := mocks_services.NewMockAuthService(ctrl)

	mock_auth.EXPECT().Check(gomock.Any(), gomock.Any()).Return(user_data, nil)
	mock_post.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := setupRouter(mock_post, mock_auth)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", SERVICE_NAME+"/post?postid="+postid, nil)
	req.Header.Set("Cookie", "token="+token+";")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}

func TestCreatePost(t *testing.T) {

	now := time.Now()
	token := "852a37a34b727c0e0b331806-7af4bdfdcc60990d427f383efecc8529289d040dd67e0753b9e2ee5a1e938402186f28324df23f6faa4e2bbf43f584ae228c55b00143866215d6e92805d470a1cc2a096dcca4d43527598122313be412e17fbefdcdab2fae02e06a405791d936862d4fba688b3c7fd784d4"
	user_data := "{\"uid\": \"1\", \"email\": \"test@email.com\", \"role\": \"standard\", \"created\": \"" + now.Format("2006-01-02T15:04:05") + "\"}"
	image_url := "http://localhost/img.png"
	caption := "test caption"

	os.Setenv("KEY", "12345678901234567890123456789012")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_post := mocks_models.NewMockPostDatabase(ctrl)
	mock_auth := mocks_services.NewMockAuthService(ctrl)

	mock_auth.EXPECT().Check(gomock.Any(), gomock.Any()).Return(user_data, nil)
	mock_post.EXPECT().Create(gomock.Any()).Return(true, nil)

	router := setupRouter(mock_post, mock_auth)

	var param = url.Values{}
	param.Set("img_url", image_url)
	param.Set("post_caption", caption)
	var payload = bytes.NewBufferString(param.Encode())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", SERVICE_NAME+"/post", payload)
	req.Header.Set("Cookie", "token="+token+";")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}

func TestDeletePost(t *testing.T) {

	now := time.Now()
	token := "852a37a34b727c0e0b331806-7af4bdfdcc60990d427f383efecc8529289d040dd67e0753b9e2ee5a1e938402186f28324df23f6faa4e2bbf43f584ae228c55b00143866215d6e92805d470a1cc2a096dcca4d43527598122313be412e17fbefdcdab2fae02e06a405791d936862d4fba688b3c7fd784d4"
	user_data := "{\"uid\": \"1\", \"email\": \"test@email.com\", \"role\": \"standard\", \"created\": \"" + now.Format("2006-01-02T15:04:05") + "\"}"
	postid := "1"

	os.Setenv("KEY", "12345678901234567890123456789012")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_post := mocks_models.NewMockPostDatabase(ctrl)
	mock_auth := mocks_services.NewMockAuthService(ctrl)

	mock_auth.EXPECT().Check(gomock.Any(), gomock.Any()).Return(user_data, nil)
	mock_post.EXPECT().Delete(gomock.Any()).Return(true, nil)

	router := setupRouter(mock_post, mock_auth)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", SERVICE_NAME+"/post?postid="+postid, nil)
	req.Header.Set("Cookie", "token="+token+";")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}

func TestGetAllPost(t *testing.T) {

	token := "852a37a34b727c0e0b331806-7af4bdfdcc60990d427f383efecc8529289d040dd67e0753b9e2ee5a1e938402186f28324df23f6faa4e2bbf43f584ae228c55b00143866215d6e92805d470a1cc2a096dcca4d43527598122313be412e17fbefdcdab2fae02e06a405791d936862d4fba688b3c7fd784d4"

	os.Setenv("KEY", "12345678901234567890123456789012")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_post := mocks_models.NewMockPostDatabase(ctrl)
	mock_auth := mocks_services.NewMockAuthService(ctrl)

	mock_post.EXPECT().FindAll().Return(make([]string, 1), nil)

	router := setupRouter(mock_post, mock_auth)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", SERVICE_NAME+"/allpost", nil)
	req.Header.Set("Cookie", "token="+token+";")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}
