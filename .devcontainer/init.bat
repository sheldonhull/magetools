powershell -NoLogo -NoProfile -ExecutionPolicy Bypass -Command "New-Item (Join-Path $ENV:USERPROFILE '.ssh') -ItemType Directory -EA 0; New-Item (Join-Path $ENV:USERPROFILE '.envrc') -ItemType File -EA 0; New-Item (Join-Path $ENV:USERPROFILE '.kube') -ItemType Directory -EA 0;"
echo '✅ prep setup complete'