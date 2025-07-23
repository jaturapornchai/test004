# Product API Testing Script for PowerShell
$API_URL = "http://localhost:8080/api"

Write-Host "=== Product API Testing ===" -ForegroundColor Green
Write-Host ""

# 1. Health Check
Write-Host "1. Testing Health Check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/health" -Method GET
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# 2. Get All Products
Write-Host "2. Getting all products..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/products" -Method GET
    Write-Host "Found $($response.data.Count) products" -ForegroundColor Cyan
    $response.data | ForEach-Object { Write-Host "- $($_.name): $($_.price) THB" }
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# 3. Get Categories
Write-Host "3. Getting all categories..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/categories" -Method GET
    Write-Host "Available categories:" -ForegroundColor Cyan
    $response.data | ForEach-Object { Write-Host "- $_" }
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# 4. Filter by category
Write-Host "4. Getting Electronics products..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/products?category=Electronics" -Method GET
    Write-Host "Found $($response.data.Count) electronics products" -ForegroundColor Cyan
    $response.data | ForEach-Object { Write-Host "- $($_.name): $($_.price) THB" }
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# 5. Filter by price range
Write-Host "5. Getting products between 1000-10000 THB..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$API_URL/products?min_price=1000&max_price=10000" -Method GET
    Write-Host "Found $($response.data.Count) products in price range" -ForegroundColor Cyan
    $response.data | ForEach-Object { Write-Host "- $($_.name): $($_.price) THB" }
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# 6. Create new product
Write-Host "6. Creating new product..." -ForegroundColor Yellow
try {
    $newProduct = @{
        name = "Samsung Galaxy S24"
        description = "Latest smartphone from Samsung"
        price = 25900
        category = "Electronics"
        stock = 20
        image_url = "https://example.com/galaxy-s24.jpg"
    }
    
    $body = $newProduct | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$API_URL/products" -Method POST -Body $body -ContentType "application/json"
    $productId = $response.data.id
    Write-Host "Created product: $($response.data.name) with ID: $productId" -ForegroundColor Cyan
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
Write-Host ""

# 7. Get specific product
if ($productId) {
    Write-Host "7. Getting specific product..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "$API_URL/products/$productId" -Method GET
        Write-Host "Product details:" -ForegroundColor Cyan
        Write-Host "- Name: $($response.data.name)"
        Write-Host "- Price: $($response.data.price) THB"
        Write-Host "- Stock: $($response.data.stock)"
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
    }
    Write-Host ""

    # 8. Update product
    Write-Host "8. Updating product..." -ForegroundColor Yellow
    try {
        $updateProduct = @{
            name = "Samsung Galaxy S24 Ultra"
            description = "Premium flagship smartphone from Samsung"
            price = 35900
            category = "Electronics"
            stock = 15
            image_url = "https://example.com/galaxy-s24-ultra.jpg"
        }
        
        $body = $updateProduct | ConvertTo-Json
        $response = Invoke-RestMethod -Uri "$API_URL/products/$productId" -Method PUT -Body $body -ContentType "application/json"
        Write-Host "Updated product: $($response.data.name)" -ForegroundColor Cyan
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
    }
    Write-Host ""

    # 9. Delete product
    Write-Host "9. Deleting product..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "$API_URL/products/$productId" -Method DELETE
        Write-Host "Product deleted successfully" -ForegroundColor Cyan
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
    }
    Write-Host ""
}

Write-Host "=== Testing Complete ===" -ForegroundColor Green
