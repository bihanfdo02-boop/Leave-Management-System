# Leave Management API Test Script

$baseUrl = "http://localhost:8080"

Write-Host "🧪 Starting API Tests..." -ForegroundColor Green

# Test 1: Health Check
Write-Host "`n✅ Test 1: Health Check"
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/health" -Method Get
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
        phone = "1234567890"
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/auth/register" `
        -Method Post `
        -ContentType "application/json" `
        -Body $body

    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $data = $response.Content | ConvertFrom-Json
    $script:token = $data.data.access_token
    Write-Host "Token saved: $($script:token.Substring(0, 20))..." -ForegroundColor Yellow
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Create Department
Write-Host "`n✅ Test 3: Create Department"
try {
    $body = @{
        name = "Engineering"
        description = "Software Development Team"
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/departments" `
        -Method Post `
        -ContentType "application/json" `
        -Headers @{"Authorization" = "Bearer $script:token"} `
        -Body $body

    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $data = $response.Content | ConvertFrom-Json
    $script:deptId = $data.data.id
    Write-Host "Department ID: $script:deptId" -ForegroundColor Yellow
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: List Departments
Write-Host "`n✅ Test 4: List Departments"
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/departments?page=1&page_size=10" `
        -Method Get `
        -Headers @{"Authorization" = "Bearer $script:token"}

    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $data = $response.Content | ConvertFrom-Json
    Write-Host "Total Departments: $($data.data.total)" -ForegroundColor Yellow
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5: Create Leave Type
Write-Host "`n✅ Test 5: Create Leave Type"
try {
    $body = @{
        name = "Annual Leave"
        description = "Paid annual vacation"
        default_days_per_year = 20
        is_paid = $true
        requires_approval = $true
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/leave-types" `
        -Method Post `
        -ContentType "application/json" `
        -Headers @{"Authorization" = "Bearer $script:token"} `
        -Body $body

    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $data = $response.Content | ConvertFrom-Json
    Write-Host "Leave Type Created: $($data.data.name)" -ForegroundColor Yellow
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6: List Leave Types
Write-Host "`n✅ Test 6: List Leave Types"
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/leave-types?page=1&page_size=10" `
        -Method Get `
        -Headers @{"Authorization" = "Bearer $script:token"}

    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $data = $response.Content | ConvertFrom-Json
    Write-Host "Total Leave Types: $($data.data.total)" -ForegroundColor Yellow
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 7: Login
Write-Host "`n✅ Test 7: Login"
try {
    $body = @{
        email = "john.doe@company.com"
        password = "password123"
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/auth/login" `
        -Method Post `
        -ContentType "application/json" `
        -Body $body

    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Green
    $data = $response.Content | ConvertFrom-Json
    Write-Host "Login successful! User: $($data.data.user.email)" -ForegroundColor Yellow
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n✨ API Tests Complete!" -ForegroundColor Green
