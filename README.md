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

### Storage

```go
package main

import "github.com/abaldeweg/warehouse-server/framework/storage"

// Create a new storage instance
s := storage.NewStorage("filesystem", "data/auth", "api_keys.json")

// Overwrite settings as needed
s.FileSystem.Path = "./data/auth"
s.FileSystem.FileName = "api_keys.json"
s.Cloud.BucketName = "storage-bucket"
s.Cloud.Path = "."
s.Cloud.FileName = "api_keys.json"

// Methods to work with the storage
k, _ := s.Save()
k, _ := s.Load()
k, _ := s.Remove()
k, _ := s.Exists()
```

## gateway

|Var                    |Description
|-----------------------|-----------
|CORS_ALLOW_ORIGIN      |Allowed origins
|API_CORE               |API endpoint for the core
|AUTH_API_ME            |Authentication API endpoint

### core

|Var|Description|Default
|---|-----------|-------
|DATABASE|Define which database to use (sqlite or mysql)|`sqlite`
|MYSQL_URL|Databse config string for MySQL|`adm:pass@tcp(localhost:3306)/warehouse?charset=utf8mb4&parseTime=True&loc=Local`
|SQLITE_NAME|Database name for SQLite (without file extension)|`warehouse`

## Static

The module sets up a simple HTTP file server that serves files from the `data` directory on port 8080.

Mount data volume to `/usr/src/app/data/`.
