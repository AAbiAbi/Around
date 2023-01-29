package main

// go m
import (
	"fmt"
	"log"
	"net/http"
	"around/util"
	// "time"
	"around/backend"
	"around/handler"
)
func main() {
    fmt.Println("started-service")
	config, err := util.LoadApplicationConfig("conf", "deploy.yml")
	if err != nil {
        panic(err)
    }

    backend.InitElasticsearchBackend(config.ElasticsearchConfig)
	backend.InitGCSBackend(config.GCSConfig)
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter(config.TokenConfig)))
    // 监听一个8080端口，然后把router拿进来
}

// main package
