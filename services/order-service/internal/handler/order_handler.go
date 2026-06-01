package handler

import (
	"net/http"
	"strconv"

	"github.com/448267450/food-delivery-platform/services/order-service/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", h.CreateOrder)
			orders.GET("/:id", h.GetOrderByID)
			orders.GET("/user/:userId", h.GetOrdersByUser)
			orders.GET("/restaurant/:restaurantId", h.GetOrdersByRestaurant)
			orders.PUT("/:id/status", h.UpdateOrderStatus)
			orders.DELETE("/:id", h.CancelOrder)
		}
	}
}

// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    order,
	})
}

// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// GET /api/v1/orders/user/:userId
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userID, err := parseID(c, "userId")
	if err != nil {
		return
	}

	orders, err := h.orderService.GetOrdersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
	})
}

// GET /api/v1/orders/restaurant/:restaurantId
func (h *OrderHandler) GetOrdersByRestaurant(c *gin.Context) {
	restaurantID, err := parseID(c, "restaurantId")
	if err != nil {
		return
	}

	orders, err := h.orderService.GetOrdersByRestaurant(restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
	})
}

// PUT /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	order, err := h.orderService.UpdateOrderStatus(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// DELETE /api/v1/orders/:id
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	if err := h.orderService.CancelOrder(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "order cancelled",
	})
}

// parseID is a helper to parse uint path params
func parseID(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid " + param,
		})
		return 0, err
	}
	return uint(val), nil
}
