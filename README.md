# warehouse-server

warehouse-server is a database to manage your warehouse.

## Getting Started

Build the images and start the container.

## Framework

Mount cover directory under `/usr/src/app/uploads/cover`.

### Router

```go
package main

import (
 "log"

 "github.com/abaldeweg/warehouse-server/blog/router"
)

func main() {
    r := router.Routes()
    log.Fatal(r.Run(":8080"))
}
```

### Config

```go
import "github.com/abaldeweg/warehouse-server/framework/config"

config.LoadAppConfig(config.WithName("myconfig"), config.WithFormat("json"), config.WithPaths("./config", "."))

viper.SetDefault("CORS_ALLOW_ORIGIN", "http://127.0.0.1")
```

### ApiKey

```go
import "github.com/abaldeweg/warehouse-server/framework/apikey"

key = apikey.NewAPIKeys([]byte(`{"keys": [{"key": "test-key", "permissions": ["read"]}]}`))
key.IsValidAPIKey("key") // returns true or false
key.HasPermission("key", "permission") // returns true or false
```

### Cors

```go
import "github.com/abaldeweg/warehouse-server/framework/cors"

r := gin.Default()
r.Use(cors.SetDefaultCorsHeaders())
```

#### Cors Settings

|Var                    |Description
|-----------------------|-----------
|CORS_ALLOW_ORIGIN      |Allowed origins

## gateway

|Var                    |Description
|-----------------------|-----------
|CORS_ALLOW_ORIGIN      |Allowed origins
|API_CORE               |API endpoint for the core
|AUTH_API_ME            |Authentication API endpoint

## core

|Var|Description|Default
|---|-----------|-------
|DATABASE|Define which database to use (sqlite or mysql)|`sqlite`
|MYSQL_URL|Databse config string for MySQL|`adm:pass@tcp(localhost:3306)/warehouse?charset=utf8mb4&parseTime=True&loc=Local`
|SQLITE_NAME|Database name for SQLite (without file extension)|`warehouse`

## Blog

Mount auth volume to `/usr/src/app/data/auth/` and data volume to `/usr/src/app/data/content/`.

The routes needs the API-Key to contain the `articles` permission.

### Blog Settings

|Var                    |Description
|-----------------------|-----------
|CORS_ALLOW_ORIGIN      |Allowed origins

## Static

Mount data volume to `/usr/src/app/data/`.

## Release

Run `make release TAG=1.0.0`.
