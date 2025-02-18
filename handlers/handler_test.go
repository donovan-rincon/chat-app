package handlers

import (
	"bytes"
	"chat-app/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) GetOrCreateChatroom(name string) (*models.Chatroom, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Chatroom), args.Error(1)
}

func (m *MockDB) GetChatroomByName(name string) (*models.Chatroom, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Chatroom), args.Error(1)
}

func (m *MockDB) CreateUserMessage(msg *models.UserMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockDB) GetLastNUserMessages(chatroomID uint, n int) ([]models.UserMessage, error) {
	args := m.Called(chatroomID, n)
	return args.Get(0).([]models.UserMessage), args.Error(1)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestRegisterHandler(t *testing.T) {
	mockDB := new(MockDB)
	router := SetupRouter(mockDB)

	user := models.User{
		Username: "testuser",
		Password: "password",
	}

	mockDB.On("CreateUser", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockDB.AssertExpectations(t)
}

func TestLoginHandler(t *testing.T) {
	mockDB := new(MockDB)
	router := SetupRouter(mockDB)

	store := cookie.NewStore([]byte("super-secret-key"))
	router.Use(sessions.Sessions("session", store))

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	dbUser := &models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockDB.On("GetUserByUsername", "testuser").Return(dbUser, nil)

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.User{
		Username: "testuser",
		Password: "password",
	})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}
