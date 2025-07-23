package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Product represents a product entity
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductRequest represents the request payload for creating/updating products
type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Category    string  `json:"category" binding:"required"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	ImageURL    string  `json:"image_url"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// In-memory database (for demo purposes)
var products = make(map[string]*Product)

// Initialize mock data
func initMockData() {
	mockProducts := []*Product{
		{
			ID:          uuid.New().String(),
			Name:        "iPhone 15 Pro",
			Description: "สมาร์ทโฟนรุ่นล่าสุดจาก Apple พร้อมชิป A17 Pro",
			Price:       39900,
			Category:    "Electronics",
			Stock:       25,
			ImageURL:    "https://example.com/iphone15pro.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Name:        "MacBook Air M2",
			Description: "แล็ปท็อปที่บางและเบาพร้อมชิป M2 ประสิทธิภาพสูง",
			Price:       35900,
			Category:    "Electronics",
			Stock:       15,
			ImageURL:    "https://example.com/macbook-air-m2.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Name:        "AirPods Pro 2",
			Description: "หูฟังไร้สายพร้อมระบบตัดเสียงรบกวน",
			Price:       8900,
			Category:    "Electronics",
			Stock:       50,
			ImageURL:    "https://example.com/airpods-pro-2.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Name:        "Nike Air Max 270",
			Description: "รองเท้าผ้าใบสำหรับการออกกำลังกายและใส่ประจำวัน",
			Price:       4500,
			Category:    "Fashion",
			Stock:       30,
			ImageURL:    "https://example.com/nike-air-max-270.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Name:        "เสื้อยืดผ้าฝ้าย 100%",
			Description: "เสื้อยืดคุณภาพดีผ้าฝ้าย 100% สีพื้นหลากหลาย",
			Price:       299,
			Category:    "Fashion",
			Stock:       100,
			ImageURL:    "https://example.com/cotton-tshirt.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Name:        "กาแฟอาราบิก้าคั่วกลาง",
			Description: "เมล็ดกาแฟอาราบิก้าคุณภาพพรีเมียม คั่วกลาง 250g",
			Price:       350,
			Category:    "Food & Beverage",
			Stock:       80,
			ImageURL:    "https://example.com/arabica-coffee.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, product := range mockProducts {
		products[product.ID] = product
	}
}

// GET /api/products - Get all products with optional filtering
func getProducts(c *gin.Context) {
	category := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	
	var filteredProducts []*Product
	
	for _, product := range products {
		// Filter by category if specified
		if category != "" && product.Category != category {
			continue
		}
		
		// Filter by price range if specified
		if minPrice != "" {
			min, err := strconv.ParseFloat(minPrice, 64)
			if err == nil && product.Price < min {
				continue
			}
		}
		
		if maxPrice != "" {
			max, err := strconv.ParseFloat(maxPrice, 64)
			if err == nil && product.Price > max {
				continue
			}
		}
		
		filteredProducts = append(filteredProducts, product)
	}
	
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    filteredProducts,
	})
}

// GET /api/products/:id - Get product by ID
func getProduct(c *gin.Context) {
	id := c.Param("id")
	
	product, exists := products[id]
	if !exists {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Product not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// POST /api/products - Create a new product
func createProduct(c *gin.Context) {
	var req ProductRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}
	
	product := &Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	products[product.ID] = product
	
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// PUT /api/products/:id - Update product by ID
func updateProduct(c *gin.Context) {
	id := c.Param("id")
	
	product, exists := products[id]
	if !exists {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Product not found",
		})
		return
	}
	
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}
	
	// Update product fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Category = req.Category
	product.Stock = req.Stock
	product.ImageURL = req.ImageURL
	product.UpdatedAt = time.Now()
	
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DELETE /api/products/:id - Delete product by ID
func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	
	_, exists := products[id]
	if !exists {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Product not found",
		})
		return
	}
	
	delete(products, id)
	
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Product deleted successfully",
	})
}

// GET /api/categories - Get all unique categories
func getCategories(c *gin.Context) {
	categoryMap := make(map[string]bool)
	
	for _, product := range products {
		categoryMap[product.Category] = true
	}
	
	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}
	
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

func main() {
	// Initialize mock data
	initMockData()
	
	// Create Gin router
	r := gin.Default()
	
	// Middleware for CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", healthCheck)
		api.GET("/products", getProducts)
		api.GET("/products/:id", getProduct)
		api.POST("/products", createProduct)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
		api.GET("/categories", getCategories)
	}
	
	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Product API Server",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":                "GET /api/health",
				"get_all_products":      "GET /api/products",
				"get_product_by_id":     "GET /api/products/:id",
				"create_product":        "POST /api/products",
				"update_product":        "PUT /api/products/:id",
				"delete_product":        "DELETE /api/products/:id",
				"get_categories":        "GET /api/categories",
			},
		})
	})
	
	// Start server
	port := ":8080"
	println("🚀 Server starting on port", port)
	println("📝 API Documentation available at http://localhost:8080")
	r.Run(port)
}
