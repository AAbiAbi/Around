//
// entity
package model

type Post struct {
    Id      string `json:"id"`//json中key value pair中拿到值，uuid
    User    string `json:"user"`
    Message string `json:"message"`
    Url     string `json:"url"`//gcs返回的
    Type    string `json:"type"`
	//public 就是大写开头
}
// mapping from json format

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int64  `json:"age"`
    Gender   string `json:"gender"`
    //可以放很多信息，不一定用
    //通过saveToES存入es

}