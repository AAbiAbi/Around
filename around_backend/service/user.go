package service

import (
    "fmt"
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"
)
//import service
func CheckUser(username string, password string) (bool, error) {
	//以user为entity,bool验证用户存不存在
    query := elastic.NewBoolQuery()//很多条件传入，都要满足。mustClauses/mustNotClauses
	//搜索名字，匹配密码
	//搜 username + password， have hits
    query.Must(elastic.NewTermQuery("username", username))
    query.Must(elastic.NewTermQuery("password", password))
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }
	//if searchResult.TotalHits() > 0
	//return true, nil
    var utype model.User
    for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		//searchResult里所有的user拿出来
        u := item.(model.User)
        if u.Password == password {
            fmt.Printf("Login as %s\n", username)
            return true, nil

        }
    }
    return false, nil
}

func AddUser(user *model.User) (bool, error) {
	//user strcut input
	//bool 添加用户是否成功，false：已存在
	//ES中若重名，直接override，所以需要提前check是否存在
    query := elastic.NewTermQuery("username", user.Username)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	//
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 {
        return false, nil
		//user已存在
    }
	//user未存在
    err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
    if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}
