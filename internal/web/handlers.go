package web

import (
	"TestAvito/internal/models"
	"TestAvito/internal/utils"
	"errors"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

func (s *Server) RegisterHandlers(m *Middleware) {
	app := s.app

	apiGroup := app.Group("/api")
	apiGroup.POST("/auth", s.Authorize)
	apiGroup.POST("/sendCoin", s.SendCoin, m.AccessLog())
	apiGroup.GET("/buy/:item", s.BuyItem, m.AccessLog())
	apiGroup.GET("/info", s.GetUserInfo, m.AccessLog())
}

func (s *Server) Authorize(c echo.Context) error {
	var req models.AuthorizeUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := s.Storage.GetUserByUsername(req.Username)

	if err != nil {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Password hashing error"})
		}
		createdUser, createErr := s.Storage.CreateUser(req.Username, hashedPassword)
		if createErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}
		user = createdUser
	} else {
		if !utils.CheckPassword(req.Password, user.Password) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
	}

	token, err := utils.GenerateToken(user.Username, s.JWT.SecretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (s *Server) SendCoin(c echo.Context) error {
	username, ok := c.Get("user_name").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, "missing token")
	}

	var req models.SendCoinRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	user, err := s.Storage.GetUserByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	recipient, err := s.Storage.GetUserByUsername(req.RecipientUsername)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	if user.Coins < req.Amount {
		return c.JSON(http.StatusBadRequest, "Not enough money")
	}

	user.Coins -= req.Amount
	recipient.Coins += req.Amount

	user, recipient, err = s.Storage.UpdateTwoUsers(user, recipient)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	transaction, err := s.Storage.CreateTransaction(user.ID, recipient.ID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":        user,
		"recipient":   recipient,
		"transaction": transaction,
	})
}

func (s *Server) BuyItem(c echo.Context) error {
	username, ok := c.Get("user_name").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, "missing token")
	}

	var req models.BuyItemRequest

	itemName := c.Param("item")
	if itemName == "" {
		return c.JSON(http.StatusBadRequest, "Item name is required")
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	user, err := s.Storage.GetUserByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	price, err := s.Storage.GetItemPrice(itemName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	if user.Coins < price*req.Quantity {
		return c.JSON(http.StatusBadRequest, "Not enough money")
	}

	user.Coins -= price * req.Quantity

	user, err = s.Storage.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	inventory, err := s.Storage.UpdateInventory(user.ID, itemName, req.Quantity)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		inventory, err = s.Storage.CreateInventory(user.ID, itemName, req.Quantity)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":       user,
		"item_price": price,
		"inventory":  inventory,
	})
}

func (s *Server) GetUserInfo(c echo.Context) error {
	username, ok := c.Get("user_name").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, "Missing token")
	}

	user, err := s.Storage.GetUserByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	inventory, err := s.Storage.GetPurchasedItems(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	transactionsFromUser, err := s.Storage.GetGiftsGivenByUser(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	transactionsToUser, err := s.Storage.GetGiftsGivenToUser(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":                   user,
		"inventory":              inventory,
		"transactions_from_user": transactionsFromUser,
		"transactions_to_user":   transactionsToUser,
	})
}
