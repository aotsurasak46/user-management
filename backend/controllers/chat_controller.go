package controllers

import (
	"fmt"
	"log"
	"sync"

	"github.com/aotsurasak46/user-management/dto"
	"github.com/aotsurasak46/user-management/models"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var clients = make(map[uint][]*websocket.Conn)
var mutex = &sync.Mutex{}

// ChatSocketHandler godoc
// @Summary WebSocket chat connection
// @Description Upgrades to WebSocket for chat. After connection, let client send JSON messages.
// @Tags chat
// @Produce json
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /ws/chat [get]
func ChatSocketHandler(db *gorm.DB) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		userIDValue := c.Locals("userID")
		userID, ok := userIDValue.(uint)
		if !ok {
			fmt.Println("Invalid or missing userID in WebSocket connection")
			return
		}

		mutex.Lock()
		clients[userID] = append(clients[userID], c)
		mutex.Unlock()
		fmt.Printf("User %d connected\n", userID)

		defer func() {
			mutex.Lock()
			connections := clients[userID]
			for i, conn := range connections {
				if conn == c {
					clients[userID] = append(connections[:i], connections[i+1:]...)
					break
				}
			}
			if len(clients[userID]) == 0 {
				delete(clients, userID)
			}
			mutex.Unlock()
			c.Close()
		}()

		for {
			requestMessage := new(dto.MessageRequest)
			if err := c.ReadJSON(requestMessage); err != nil {
				fmt.Println("Error reading json:", err)
				break
			}

			if requestMessage.To == 0 || requestMessage.Content == "" {
				fmt.Println("Invalid message data.")
				continue
			}

			message := models.Message{
				ToID:    requestMessage.To,
				Content: requestMessage.Content,
				FromID:  userID,
			}

			if err := db.Create(&message).Error; err != nil {
				fmt.Printf("Failed to save message: %v", err)
				continue
			}

			if err := db.Preload("From").Preload("To").First(&message, message.ID).Error; err != nil {
				fmt.Printf("Failed to load associations: %v", err)
			}

			c.WriteJSON(map[string]any{
				"type": "sent",
				"data": fiber.Map{
					"ID":        message.ID,
					"content":   message.Content,
					"from_id":   message.FromID,
					"to_id":     message.ToID,
					"timestamp": message.Timestamp,
					"tempId":    requestMessage.TempID,
				},
			})

			mutex.Lock()
			if message.FromID != message.ToID {
				if conns, ok := clients[message.ToID]; ok {
					for _, conn := range conns {
						err := conn.WriteJSON(map[string]any{
							"type": "incoming",
							"data": message,
						})
						if err != nil {
							fmt.Println("Error sending to recipient:", err)
						} else {
							fmt.Printf("Message from user %d sent to user %d: %s\n", message.FromID, message.ToID, message.Content)
						}
					}
				}
			}
			mutex.Unlock()
		}
	})
}

// GetChatHistory godoc
// @Summary Get chat history of user
// @Description Get all message in chat of user with another user
// @Tags chat
// @Accept json
// @Produce json
// @param id path int true "User id"
// @Success 200 {array}  dto.MessageResponse
// @Failure 400 {object} object{error=string} "Bad request, User ID is missing in the request"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/messages/:userId [get]
func GetChatHistory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fromUserID := c.Locals("userID")
		fromID, ok := fromUserID.(uint)
		if !ok {
			log.Println("Invalid or missing fromUserID in request context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		toId := c.Params("userId")
		if toId == "" {
			fmt.Printf("User ID is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
		}

		messages := new([]models.Message)
		if err := db.Where(
			"(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
			fromID, toId, toId, fromID,
		).Order("timestamp ASC").Find(messages).Error; err != nil {
			fmt.Printf("Error finding chat history in database: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get chat history"})
		}

		return c.JSON(messages)
	}
}

// GetConversations godoc
// @Summary Get conversations of user
// @Description Get conversations of user that include the latest message exchanged with each participant
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {array}  dto.ConversationResponse
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/conversations [get]
func GetConversations(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDValue := c.Locals("userID")
		userId, ok := userIDValue.(uint)
		if !ok {
			log.Println("Invalid or missing UserID in request context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var messages []models.Message
		if err := db.Raw(`
			SELECT * FROM messages m1
			WHERE (m1.from_id = ? OR m1.to_id = ?)
			AND m1.timestamp = (
				SELECT MAX(m2.timestamp)
				FROM messages m2
				WHERE ((m2.from_id = m1.from_id AND m2.to_id = m1.to_id)
				OR (m2.from_id = m1.to_id AND m2.to_id = m1.from_id))
			)
			ORDER BY m1.timestamp DESC
		`, userId, userId).Scan(&messages).Error; err != nil {
			fmt.Printf("Failed to fetch conversations : %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var conversations []dto.ConversationResponse

		for _, msg := range messages {
			var otherUser models.User
			var otherUserID uint

			if msg.FromID == userId {
				otherUserID = msg.ToID
			} else {
				otherUserID = msg.FromID
			}

			if err := db.Where("id = ?", otherUserID).First(&otherUser).Error; err != nil {
				log.Printf("Failed to fetch user %d: %v", otherUserID, err)
				continue
			}

			conversations = append(conversations, dto.ConversationResponse{
				User: dto.UserResponse{
					ID:        otherUser.ID,
					Name:      otherUser.Name,
					Email:     otherUser.Email,
					CreatedAt: otherUser.CreatedAt,
					UpdatedAt: otherUser.UpdatedAt,
					Role:      otherUser.Role,
				},
				LastMessage: msg.Content,
				Timestamp:   msg.Timestamp,
			})
		}

		return c.JSON(conversations)
	}
}