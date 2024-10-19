# warehouse-server

warehouse-server is a database to manage your warehouse.

## Getting Started

## Blog

Mount auth volume to `/usr/src/app/data/auth/` and data volume to `/usr/src/app/data/content/`.

The routes needs the API-Key to contain the `articles` permission.

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

apikey.IsValidAPIKey("key")
apikey.HasPermission("key", "permission")
```

### Cors

```go
import "github.com/abaldeweg/warehouse-server/framework/cors"

r := gin.Default()
r.Use(cors.SetDefaultCorsHeaders())
```

## Settings

|Var                    |Description                                |Used by
|-----------------------|-------------------------------------------|--------------------------------
|CORS_ALLOW_ORIGIN      |Allowed origins                            |gateway, blog
|API_CORE               |API endpoint for the core                  |gateway
|project_dir            |Path to docker compose                     |admincli
|database               |Database name to dump                      |admincli
|AUTH_API_ME            |Authentication API endpoint                |gateway

admincli will read a config file from following paths:

- /etc/admincli/admincli.yaml

Example

```yaml
// admincli.yaml
project_dir: .
database: db-1
```
