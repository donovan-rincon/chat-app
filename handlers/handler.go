package handlers

import (
	"log"
	"net/http"
	"time"

	"chat-app/db"
	"chat-app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var clients = make(map[uint]map[*websocket.Conn]bool)
var broadcast = make(map[uint]chan models.UserMessage)
var upgrader = websocket.Upgrader{}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("session", store))

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/chatroom/:name", authMiddleware(), chatroomHandler)
	r.GET("/ws/:name", authMiddleware(), handleConnections)

	r.Static("/public", "./public")
	r.Static("/css", "./public/css")
	r.LoadHTMLGlob("public/*.html")

	return r
}

func registerHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err := db.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "user registered"})
}

func loginHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbUser, err := db.GetUserByUsername(user.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"redirect": "/chatroom/home"}) // Redirect to default chatroom
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("username")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func chatroomHandler(c *gin.Context) {
	chatroomName := c.Param("name")
	if chatroomName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chatroom name is required"})
		return
	}
	log.Printf("Chatroom requested: %s", chatroomName)
	chatroom, err := db.GetOrCreateChatroom(chatroomName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "chatroom creation failed"})
		return
	}
	c.HTML(http.StatusOK, "chatroom.html", gin.H{"chatroom": chatroom.Name})
}

func handleConnections(c *gin.Context) {
	chatroomName := c.Param("name")
	chatroom, err := db.GetChatroomByName(chatroomName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "chatroom not found"})
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	if clients[chatroom.ID] == nil {
		clients[chatroom.ID] = make(map[*websocket.Conn]bool)
		broadcast[chatroom.ID] = make(chan models.UserMessage)
		go handleMessages(chatroom)
	}

	clients[chatroom.ID][ws] = true

	for {
		var wsMsg models.WSMessage
		err := ws.ReadJSON(&wsMsg)
		if err != nil {
			delete(clients[chatroom.ID], ws)
			break
		}

		user, err := db.GetUserByUsername(wsMsg.Username)
		if err != nil {
			log.Printf("Failed to retrieve user: %v", err)
			continue
		}

		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05")

		msg := models.UserMessage{
			ChatroomID: chatroom.ID,
			UserID:     user.ID,
			Username:   user.Username,
			Message:    wsMsg.Message,
			Timestamp:  formattedTime,
		}
		err = db.CreateUserMessage(&msg)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			continue
		}
		broadcast[chatroom.ID] <- msg
	}
}

func handleMessages(chatroom *models.Chatroom) {
	for {
		msg := <-broadcast[chatroom.ID]
		log.Printf("New message in chatroom '%s': %v", chatroom.Name, msg)

		messages, err := db.GetLastUserMessages(chatroom.ID, 50)
		if err != nil {
			log.Printf("Failed to retrieve messages: %v", err)
			continue
		}

		for client := range clients[chatroom.ID] {
			err := client.WriteJSON(messages)
			if err != nil {
				client.Close()
				delete(clients[chatroom.ID], client)
			}
		}
	}
}
