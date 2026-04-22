param(
  [Parameter(Position = 0)]
  [ValidateSet("up", "down", "version", "force", "goto", "drop")]
  [string]$Command = "up",

  [string]$EnvFile = ".env",

  [Parameter(ValueFromRemainingArguments = $true)]
  [string[]]$CommandArgs
)

$ErrorActionPreference = "Stop"

function Read-DotEnv {
  param([string]$Path)

  if (-not (Test-Path -LiteralPath $Path)) {
    throw "Env file not found: $Path"
  }

  $map = @{}
  foreach ($rawLine in Get-Content -LiteralPath $Path) {
    $line = $rawLine.Trim()
    if ([string]::IsNullOrWhiteSpace($line)) { continue }
    if ($line.StartsWith("#")) { continue }

    $idx = $line.IndexOf("=")
    if ($idx -lt 1) { continue }

    $key = $line.Substring(0, $idx).Trim()
    $value = $line.Substring($idx + 1).Trim()

    if (
      ($value.StartsWith('"') -and $value.EndsWith('"')) -or
      ($value.StartsWith("'") -and $value.EndsWith("'"))
    ) {
      $value = $value.Substring(1, $value.Length - 2)
    }

    $map[$key] = $value
  }

  return $map
}

function Get-ConfigValue {
  param(
    [hashtable]$Map,
    [string]$Key,
    [string]$Default = ""
  )

  if ($Map.ContainsKey($Key) -and -not [string]::IsNullOrWhiteSpace($Map[$Key])) {
    return $Map[$Key]
  }

  return $Default
}

function Resolve-MigrationPath {
  param(
    [string]$Candidate,
    [string]$BackendDir,
    [string]$RepoDir
  )

  if ([string]::IsNullOrWhiteSpace($Candidate)) {
    return (Join-Path $BackendDir "migration")
  }

  if ([System.IO.Path]::IsPathRooted($Candidate)) {
    return $Candidate
  }

  $fromBackend = Join-Path $BackendDir $Candidate
  if (Test-Path -LiteralPath $fromBackend) {
    return $fromBackend
  }

  $fromRepo = Join-Path $RepoDir $Candidate
  if (Test-Path -LiteralPath $fromRepo) {
    return $fromRepo
  }

  return $fromBackend
}

$migrateCmd = Get-Command migrate -ErrorAction SilentlyContinue
if (-not $migrateCmd) {
  throw "'migrate' command not found. Install golang-migrate CLI with MySQL driver first."
}

$backendDir = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$repoDir = (Resolve-Path (Join-Path $backendDir "..")).Path
$envPath = Join-Path $backendDir $EnvFile
$config = Read-DotEnv -Path $envPath

$mysqlHost = Get-ConfigValue -Map $config -Key "MYSQL_HOST" -Default "127.0.0.1"
$mysqlPort = Get-ConfigValue -Map $config -Key "MYSQL_PORT" -Default "3306"
$mysqlUser = Get-ConfigValue -Map $config -Key "MYSQL_USER" -Default "root"
$mysqlPassword = Get-ConfigValue -Map $config -Key "MYSQL_PASSWORD"
$mysqlDatabase = Get-ConfigValue -Map $config -Key "MYSQL_DATABASE" -Default "gallery_db"
$mysqlParams = Get-ConfigValue -Map $config -Key "MYSQL_PARAMS"

if ([string]::IsNullOrWhiteSpace($mysqlPassword)) {
  throw "MYSQL_PASSWORD is required in $envPath"
}

Push-Location $repoDir
try {
  $migrationPathRaw = Get-ConfigValue -Map $config -Key "MIGRATION_PATH" -Default "backend/migration"
  $migrationPath = Resolve-MigrationPath -Candidate $migrationPathRaw -BackendDir $backendDir -RepoDir $repoDir
  $migrationPathAbsolute = (Resolve-Path -LiteralPath $migrationPath).Path
  $sourcePathRelative = (Resolve-Path -LiteralPath $migrationPathAbsolute -Relative)
  $sourcePathRelative = $sourcePathRelative.TrimStart('.', '\', '/').Replace('\', '/')
  $sourceUri = "file://$sourcePathRelative"

  $userEncoded = [System.Uri]::EscapeDataString($mysqlUser)
  $passEncoded = [System.Uri]::EscapeDataString($mysqlPassword)

  $queryItems = New-Object System.Collections.Generic.List[string]
  $queryItems.Add("x-multi-statement=true")

  if (-not [string]::IsNullOrWhiteSpace($mysqlParams)) {
    $queryItems.Add($mysqlParams.TrimStart("?"))
  }

  $queryString = [string]::Join("&", $queryItems)
  $databaseUrl = "mysql://${userEncoded}:${passEncoded}@tcp(${mysqlHost}:${mysqlPort})/${mysqlDatabase}?$queryString"

  $maskedUrl = "mysql://${mysqlUser}:***@tcp(${mysqlHost}:${mysqlPort})/${mysqlDatabase}?$queryString"
  Write-Host "Migration path: $migrationPathAbsolute"
  Write-Host "Migration source: $sourceUri"
  Write-Host "Database URL : $maskedUrl"

  $invokeArgs = @("-source", $sourceUri, "-database", $databaseUrl, $Command)
  if ($CommandArgs) {
    $invokeArgs += $CommandArgs
  }

  & migrate @invokeArgs
  exit $LASTEXITCODE
}
finally {
  Pop-Location
}
