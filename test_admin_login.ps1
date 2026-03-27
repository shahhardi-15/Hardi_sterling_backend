$uri = "http://localhost:5000/api/admin/login"
$headers = @{"Content-Type" = "application/json"}
$body = @{
    email = "adminsterling@gmail.com"
    password = "admin@123"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri $uri -Method Post -Headers $headers -Body $body
    Write-Host "Admin Login Test Response:"
    $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $($_.Exception.Message)"
}
