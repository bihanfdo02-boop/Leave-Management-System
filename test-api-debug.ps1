# Leave Management API Test Script with Debug Info

$baseUrl = "http://localhost:8080"

Write-Host "🧪 Starting API Tests..." -ForegroundColor Green

# Test 1: Health Check
Write-Host "`n✅ Test 1: Health Check"
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/health" -Method Get -UseBasicParsing
    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $response.Content | ConvertFrom-Json | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2: Register User
Write-Host "`n✅ Test 2: Register User"
try {
    $body = @{
        email = "john.doe@company.com"
        password = "password123"
        first_name = "John"
        last_name = "Doe"
        phone = "12345678901"
    } | ConvertTo-Json

    Write-Host "Sending request body:" -ForegroundColor Cyan
    Write-Host $body

    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/auth/register" `
        -Method Post `
        -ContentType "application/json" `
        -Body $body `
        -UseBasicParsing

    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $data = $response.Content | ConvertFrom-Json
    Write-Host "Response:" -ForegroundColor Cyan
    $data | ConvertTo-Json | Write-Host
    
    $script:token = $data.data.access_token
    Write-Host "Token saved: $($script:token.Substring(0, 20))..." -ForegroundColor Yellow
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Response Body:" -ForegroundColor Cyan
    try {
        $errorResponse = $_.ErrorDetails.Message | ConvertFrom-Json
        $errorResponse | ConvertTo-Json | Write-Host
    } catch {
        Write-Host $_.ErrorDetails.Message
    }
}

# Test 3: Create Department
Write-Host "`n✅ Test 3: Create Department"
try {
    if ($null -eq $script:token) {
        Write-Host "⚠️  Skipping - No token available" -ForegroundColor Yellow
    } else {
        $body = @{
            name = "Engineering"
            description = "Software Development Team"
        } | ConvertTo-Json

        $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/departments" `
            -Method Post `
            -ContentType "application/json" `
            -Headers @{"Authorization" = "Bearer $script:token"} `
            -Body $body `
            -UseBasicParsing

        Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
        $data = $response.Content | ConvertFrom-Json
        $script:deptId = $data.data.id
        Write-Host "Department ID: $script:deptId" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: List Departments
Write-Host "`n✅ Test 4: List Departments"
try {
    if ($null -eq $script:token) {
        Write-Host "⚠️  Skipping - No token available" -ForegroundColor Yellow
    } else {
        $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/departments?page=1&page_size=10" `
            -Method Get `
            -Headers @{"Authorization" = "Bearer $script:token"} `
            -UseBasicParsing

        Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
        $data = $response.Content | ConvertFrom-Json
        Write-Host "Total Departments: $($data.data.total)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n✨ API Tests Complete!" -ForegroundColor Green