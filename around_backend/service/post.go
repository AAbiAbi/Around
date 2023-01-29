package service

import (
    "mime/multipart"
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"
)

func SearchPostsByUser(user string) ([]model.Post, error) {
    query := elastic.NewTermQuery("user", user)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
    query := elastic.NewMatchQuery("message", keywords)
    query.Operator("AND")
	// 多个keyword传入操作是AND
    if keywords == "" {
        query.ZeroTermsQuery("all")
    }
	//cornercase没有东西返回所有
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
    var ptype model.Post
    var posts []model.Post

    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
        p := item.(model.Post)
        posts = append(posts, p)
    }
	return posts
}
func SavePost(post *model.Post, file multipart.File) error {
    // save to GCS to get url
    medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
    if err != nil {
        return err
    }
    // save to ES
    post.Url = medialink
    //剩下的都是从handler传进来的
    return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
}
func DeletePost(id string, user string) error {
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("id", id))
    query.Must(elastic.NewTermQuery("user", user))

    return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}
