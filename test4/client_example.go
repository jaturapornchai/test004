package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const API_BASE = "http://localhost:8080/api"

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	fmt.Println("🛍️  Product API Client Example")
	fmt.Println("================================")

	// 1. Get all products
	fmt.Println("\n1. Getting all products...")
	products, err := getAllProducts()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Found %d products:\n", len(products))
	for _, p := range products {
		fmt.Printf("- %s: %.2f THB\n", p.Name, p.Price)
	}

	// 2. Create new product
	fmt.Println("\n2. Creating new product...")
	newProduct := Product{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       999.99,
		Category:    "Test",
		Stock:       10,
		ImageURL:    "https://example.com/test.jpg",
	}
	
	created, err := createProduct(newProduct)
	if err != nil {
		fmt.Printf("Error creating product: %v\n", err)
		return
	}
	
	fmt.Printf("Created product: %s (ID: %s)\n", created.Name, created.ID)

	// 3. Get specific product
	fmt.Println("\n3. Getting specific product...")
	product, err := getProduct(created.ID)
	if err != nil {
		fmt.Printf("Error getting product: %v\n", err)
		return
	}
	
	fmt.Printf("Retrieved: %s - %.2f THB\n", product.Name, product.Price)

	// 4. Update product
	fmt.Println("\n4. Updating product...")
	product.Name = "Updated Test Product"
	product.Price = 1299.99
	
	updated, err := updateProduct(created.ID, *product)
	if err != nil {
		fmt.Printf("Error updating product: %v\n", err)
		return
	}
	
	fmt.Printf("Updated: %s - %.2f THB\n", updated.Name, updated.Price)

	// 5. Delete product
	fmt.Println("\n5. Deleting product...")
	err = deleteProduct(created.ID)
	if err != nil {
		fmt.Printf("Error deleting product: %v\n", err)
		return
	}
	
	fmt.Println("Product deleted successfully!")

	// 6. Get filtered products
	fmt.Println("\n6. Getting Electronics products...")
	electronics, err := getProductsByCategory("Electronics")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Found %d electronics:\n", len(electronics))
	for _, p := range electronics {
		fmt.Printf("- %s: %.2f THB\n", p.Name, p.Price)
	}
}

func getAllProducts() ([]Product, error) {
	resp, err := http.Get(API_BASE + "/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	// Convert interface{} to []Product
	var products []Product
	jsonData, _ := json.Marshal(apiResp.Data)
	json.Unmarshal(jsonData, &products)

	return products, nil
}

func getProduct(id string) (*Product, error) {
	resp, err := http.Get(API_BASE + "/products/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var product Product
	jsonData, _ := json.Marshal(apiResp.Data)
	json.Unmarshal(jsonData, &product)

	return &product, nil
}

func createProduct(product Product) (*Product, error) {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(API_BASE+"/products", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var created Product
	jsonData, _ = json.Marshal(apiResp.Data)
	json.Unmarshal(jsonData, &created)

	return &created, nil
}

func updateProduct(id string, product Product) (*Product, error) {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", API_BASE+"/products/"+id, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var updated Product
	jsonData, _ = json.Marshal(apiResp.Data)
	json.Unmarshal(jsonData, &updated)

	return &updated, nil
}

func deleteProduct(id string) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", API_BASE+"/products/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getProductsByCategory(category string) ([]Product, error) {
	resp, err := http.Get(API_BASE + "/products?category=" + category)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	var products []Product
	jsonData, _ := json.Marshal(apiResp.Data)
	json.Unmarshal(jsonData, &products)

	return products, nil
}
