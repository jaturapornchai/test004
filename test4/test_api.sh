#!/bin/bash

# API Testing Script
API_URL="http://localhost:8080/api"

echo "=== Product API Testing ==="
echo

# 1. Health Check
echo "1. Testing Health Check..."
curl -s $API_URL/health | jq .
echo -e "\n"

# 2. Get All Products
echo "2. Getting all products..."
curl -s $API_URL/products | jq .
echo -e "\n"

# 3. Get Categories
echo "3. Getting all categories..."
curl -s $API_URL/categories | jq .
echo -e "\n"

# 4. Filter by category
echo "4. Getting Electronics products..."
curl -s "$API_URL/products?category=Electronics" | jq .
echo -e "\n"

# 5. Filter by price range
echo "5. Getting products between 1000-10000 baht..."
curl -s "$API_URL/products?min_price=1000&max_price=10000" | jq .
echo -e "\n"

# 6. Create new product
echo "6. Creating new product..."
NEW_PRODUCT='{
  "name": "Samsung Galaxy S24",
  "description": "สมาร์ทโฟนรุ่นใหม่จาก Samsung",
  "price": 25900,
  "category": "Electronics",
  "stock": 20,
  "image_url": "https://example.com/galaxy-s24.jpg"
}'

RESPONSE=$(curl -s -X POST $API_URL/products \
  -H "Content-Type: application/json" \
  -d "$NEW_PRODUCT")

echo $RESPONSE | jq .
PRODUCT_ID=$(echo $RESPONSE | jq -r '.data.id')
echo -e "\nCreated product ID: $PRODUCT_ID\n"

# 7. Get specific product
echo "7. Getting specific product..."
curl -s $API_URL/products/$PRODUCT_ID | jq .
echo -e "\n"

# 8. Update product
echo "8. Updating product..."
UPDATE_PRODUCT='{
  "name": "Samsung Galaxy S24 Ultra",
  "description": "สมาร์ทโฟนระดับท็อปจาก Samsung",
  "price": 35900,
  "category": "Electronics",
  "stock": 15,
  "image_url": "https://example.com/galaxy-s24-ultra.jpg"
}'

curl -s -X PUT $API_URL/products/$PRODUCT_ID \
  -H "Content-Type: application/json" \
  -d "$UPDATE_PRODUCT" | jq .
echo -e "\n"

# 9. Delete product
echo "9. Deleting product..."
curl -s -X DELETE $API_URL/products/$PRODUCT_ID | jq .
echo -e "\n"

echo "=== Testing Complete ==="
