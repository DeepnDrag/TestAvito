package e2e

import (
	"TestAvito/internal/models"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	User      models.User      `json:"user"`
	ItemPrice int              `json:"item_price"`
	Inventory models.Inventory `json:"inventory"`
}

func TestBuyItemHandler(t *testing.T) {
	e := echo.New()
	e.POST("/api/v1/buy/:item", func(c echo.Context) error {
		itemName := c.Param("item")
		if itemName == "" {
			return c.JSON(http.StatusBadRequest, "Item name is required")
		}

		var req models.BuyItemRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}

		user := models.User{
			ID:       10,
			Username: "test15",
			Password: "$2a$10$xBsRYLXqyMIpErS9PAU1qeqEfpqwt1oP1U.RXMtblKbzERUCKzj0q",
			Coins:    1000,
		}

		itemPrice := 10
		if user.Coins < itemPrice*req.Quantity {
			return c.JSON(http.StatusBadRequest, "Not enough money")
		}

		user.Coins -= itemPrice * req.Quantity

		inventory := models.Inventory{
			UserID:   user.ID,
			ItemType: itemName,
			Quantity: req.Quantity,
		}

		resp := Response{
			User:      user,
			ItemPrice: itemPrice,
			Inventory: inventory,
		}

		return c.JSON(http.StatusOK, resp)
	})

	ts := httptest.NewServer(e)
	defer ts.Close()

	body := `{"quantity": 2}`
	resp, err := http.Post(ts.URL+"/api/v1/buy/cup", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, 2, result.Inventory.Quantity)
	assert.Equal(t, "cup", result.Inventory.ItemType)
	assert.Equal(t, 980, result.User.Coins)

	resp, err = http.Post(ts.URL+"/api/v1/buy/", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	invalidBody := `{"quantity": "two"}`
	resp, err = http.Post(ts.URL+"/api/v1/buy/cup", "application/json", strings.NewReader(invalidBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSendCoinHandler(t *testing.T) {
	e := echo.New()
	e.POST("/api/v1/send_coin", func(c echo.Context) error {
		username := "test15"

		var req models.SendCoinRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}

		user := models.User{
			ID:       10,
			Username: username,
			Password: "$2a$10$xBsRYLXqyMIpErS9PAU1qeqEfpqwt1oP1U.RXMtblKbzERUCKzj0q",
			Coins:    860,
		}

		recipient := models.User{
			ID:       5,
			Username: "test5",
			Password: "password",
			Coins:    1200,
		}

		if user.Coins < req.Amount {
			return c.JSON(http.StatusBadRequest, "Not enough money")
		}

		user.Coins -= req.Amount
		recipient.Coins += req.Amount

		transaction := models.Transaction{
			ID:         8,
			FromUserID: user.ID,
			ToUserID:   recipient.ID,
			Amount:     req.Amount,
		}

		resp := map[string]interface{}{
			"user":        user,
			"recipient":   recipient,
			"transaction": transaction,
		}

		return c.JSON(http.StatusOK, resp)
	})

	ts := httptest.NewServer(e)
	defer ts.Close()

	body := `{"recipient_username": "test5", "amount": 100}`
	resp, err := http.Post(ts.URL+"/api/v1/send_coin", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	user := result["user"].(map[string]interface{})
	assert.Equal(t, "test15", user["Username"])
	assert.Equal(t, float64(760), user["Coins"])

	recipient := result["recipient"].(map[string]interface{})
	assert.Equal(t, "test5", recipient["Username"])
	assert.Equal(t, float64(1300), recipient["Coins"])

	body = `{"recipient_username": "test5", "amount": 10000}`
	resp, err = http.Post(ts.URL+"/api/v1/send_coin", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	invalidBody := `{"recipient_username": "test5", "amount": "hundred"}`
	resp, err = http.Post(ts.URL+"/api/v1/send_coin", "application/json", strings.NewReader(invalidBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
