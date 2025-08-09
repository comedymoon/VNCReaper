# ===== CONFIG =====
$Version = "v1.1"
$AppName = "vncreaper"
$EntryPoint = "./cmd/vncreaper"
$OutDir = "./bin"

# ===== PREP =====
Write-Host "[*] Cleaning output directory..."
if (Test-Path $OutDir) { Remove-Item $OutDir -Recurse -Force }
New-Item -ItemType Directory -Path $OutDir | Out-Null

# ===== TARGETS =====
$Platforms = @(
    @{ GOOS="windows"; GOARCH="amd64" },
    @{ GOOS="windows"; GOARCH="arm64" },
    @{ GOOS="linux";   GOARCH="amd64" },
    @{ GOOS="linux";   GOARCH="arm64" }
)

# ===== BUILD LOOP =====
foreach ($Target in $Platforms) {
    $GOOS = $Target.GOOS
    $GOARCH = $Target.GOARCH
    Write-Host "[*] Building for $GOOS/$GOARCH..."

    $Output = "$OutDir/$AppName-$Version-$GOOS-$GOARCH"
    if ($GOOS -eq "windows") { $Output += ".exe" }

    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH

    go build `
        -trimpath `
        -ldflags="-s -w -buildid=" `
        -mod=vendor `
        -o $Output `
        $EntryPoint
}

Write-Host "[+] All builds completed! Files are in: $OutDir"