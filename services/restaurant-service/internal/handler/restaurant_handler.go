package handler

import (
	"net/http"
	"strconv"

	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/service"
	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	restaurantService service.RestaurantService
}

func NewRestaurantHandler(restaurantService service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{restaurantService: restaurantService}
}

func (h *RestaurantHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		restaurants := v1.Group("/restaurants")
		{
			restaurants.POST("", h.CreateRestaurant)
			restaurants.GET("", h.GetAllRestaurants)
			restaurants.GET("/:id", h.GetRestaurantByID)
			restaurants.PUT("/:id", h.UpdateRestaurant)
			restaurants.DELETE("/:id", h.DeleteRestaurant)

			// Menu item routes nested under restaurant
			restaurants.POST("/:id/menu", h.AddMenuItem)
			restaurants.PUT("/:id/menu/:itemId", h.UpdateMenuItem)
			restaurants.DELETE("/:id/menu/:itemId", h.DeleteMenuItem)
		}
	}
}

// -------------------- Restaurant Handlers --------------------

// POST /api/v1/restaurants
func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	var req service.CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	restaurant, err := h.restaurantService.CreateRestaurant(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    restaurant,
	})
}

// GET /api/v1/restaurants
func (h *RestaurantHandler) GetAllRestaurants(c *gin.Context) {
	restaurants, err := h.restaurantService.GetAllRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    restaurants,
	})
}

// GET /api/v1/restaurants/:id
func (h *RestaurantHandler) GetRestaurantByID(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    restaurant,
	})
}

// PUT /api/v1/restaurants/:id
func (h *RestaurantHandler) UpdateRestaurant(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	var req service.UpdateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	restaurant, err := h.restaurantService.UpdateRestaurant(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    restaurant,
	})
}

// DELETE /api/v1/restaurants/:id
func (h *RestaurantHandler) DeleteRestaurant(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	if err := h.restaurantService.DeleteRestaurant(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "restaurant deleted",
	})
}

// -------------------- Menu Item Handlers --------------------

// POST /api/v1/restaurants/:id/menu
func (h *RestaurantHandler) AddMenuItem(c *gin.Context) {
	restaurantID, err := parseID(c, "id")
	if err != nil {
		return
	}

	var req service.AddMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	item, err := h.restaurantService.AddMenuItem(restaurantID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    item,
	})
}

// PUT /api/v1/restaurants/:id/menu/:itemId
func (h *RestaurantHandler) UpdateMenuItem(c *gin.Context) {
	restaurantID, err := parseID(c, "id")
	if err != nil {
		return
	}

	itemID, err := parseID(c, "itemId")
	if err != nil {
		return
	}

	var req service.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	item, err := h.restaurantService.UpdateMenuItem(restaurantID, itemID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    item,
	})
}

// DELETE /api/v1/restaurants/:id/menu/:itemId
func (h *RestaurantHandler) DeleteMenuItem(c *gin.Context) {
	restaurantID, err := parseID(c, "id")
	if err != nil {
		return
	}

	itemID, err := parseID(c, "itemId")
	if err != nil {
		return
	}

	if err := h.restaurantService.DeleteMenuItem(restaurantID, itemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "menu item deleted",
	})
}

// -------------------- Helper --------------------

// parseID is a helper to parse uint path params and write error response if invalid
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
