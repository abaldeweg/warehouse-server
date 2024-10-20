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

## admincli

### admincli Settings

|Var                    |Description
|-----------------------|-----------
|project_dir            |Path to docker compose
|database               |Database name to dump

admincli will read a config file from following paths:

- /etc/admincli/admincli.yaml

Example

```yaml
// admincli.yaml
project_dir: .
database: db-1
```

## gateway

|Var                    |Description
|-----------------------|-----------
|CORS_ALLOW_ORIGIN      |Allowed origins
|API_CORE               |API endpoint for the core
|AUTH_API_ME            |Authentication API endpoint

## Blog

Mount auth volume to `/usr/src/app/data/auth/` and data volume to `/usr/src/app/data/content/`.

The routes needs the API-Key to contain the `articles` permission.

### Blog Settings

|Var                    |Description
|-----------------------|-----------
|CORS_ALLOW_ORIGIN      |Allowed origins
