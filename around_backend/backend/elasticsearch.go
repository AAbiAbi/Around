package backend

import (
    "context"
    "fmt"

    "around/constants"
    "around/util"

    "github.com/olivere/elastic/v7"
)

var (
    ESBackend *ElasticsearchBackend
    //只能实例化一个，singleton
)

type ElasticsearchBackend struct {
    client *elastic.Client
}

func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
    //初始化 elastic client
    client, err := elastic.NewClient(
        elastic.SetURL(config.Address),
        elastic.SetBasicAuth(config.Username, config.Password))
        // 创建client
    if err != nil {
        panic(err)
        //Handle error
    }

    exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
            "mappings": {
               
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
        // 创建一个json格式的String，因为elasticSearch是用json格式进行交互的
         // 五个column
         //keyword：搜索时必须完全匹配； text：类似于contains
        //  index是索引的意思，在树里面查找比较快，此处关闭索引功能，减少properties的消耗
        _, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
        // POST_INDEX是名字
        if err != nil {
            panic(err)
        }
    }

    exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    // context是一些环境变量，给一个deadline，若未完成request在ddl内，则认为error
    // exists是boolean值
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
                // user profile
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")

    ESBackend = &ElasticsearchBackend{client: client}
    // 给ElasticsearchBackend赋值
    //ESBackend 是一个zhi zhen需要dereference
}

// 连接elasticSearch的后端
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).
        // index比较灵活，搜谁都行
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())
    return err
}
func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
    _, err := backend.client.DeleteByQuery().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())

    return err
}
