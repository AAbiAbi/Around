// 处理不同的关于post的操作
package handler

import (
	"encoding/json" //把client传过来的json进行decode
	"fmt"           //输出
	"net/http"      //ResponseWriter

	"around/model"
	"around/service"

	"path/filepath"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	//jwt "github.com/form3tech-oss/jwt-go"
)
var (
    mediaTypes = map[string]string{
        ".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
    }
)
// hash map 


func uploadHandler(w http.ResponseWriter, r *http.Request) {//r是指针，根据pass by value原则，同时也可以省空间
	// w 不是pointer， 是interface，没有object，不可以用pointer
    // Parse from body of request to get a json object.
    // fmt.Println("Received one upload request")
    // decoder := json.NewDecoder(r.Body)
    // var p model.Post
    // if err := decoder.Decode(&p); err != nil {
    //     panic(err)
    // }

   //fmt.Fprintf(w, "Post received: %s\n", p.Message)
	//handler控制往writer里写东西
	// writer只是buffer，最终要写入request.body
    //read data from form data
    fmt.Println("Received one upload request")
    token := r.Context().Value("user")//jwt token
    claims := token.(*jwt.Token).Claims//key value pair
    username := claims.(jwt.MapClaims)["username"]
    p := model.Post{
        Id:      uuid.New(),
        User:    username.(string),

        Message: r.FormValue("message"),
    }
// user不从前端读，从token来获得就可以
    file, header, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    suffix := filepath.Ext(header.Filename)
    if t, ok := mediaTypes[suffix]; ok {
        p.Type = t
    } else {
        p.Type = "unknown"
    }

    err = service.SavePost(&p, file)
    if err != nil {
        http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save post to backend %v\n", err)
        return
    }

    fmt.Println("Post is saved successfully.")


}
func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for search")
    w.Header().Set("Content-Type", "application/json")

    user := r.URL.Query().Get("user")
    // 既声明又赋值
    keywords := r.URL.Query().Get("keywords")

    var posts []model.Post
    var err error
    if user != "" {
        // user传进来一个数
        posts, err = service.SearchPostsByUser(user)
    } else {

        posts, err = service.SearchPostsByKeywords(keywords)
    }
    
    if err != nil {
        http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
        // 500 error
        fmt.Printf("Failed to read post from backend %v.\n", err)
        return
    }

    js, err := json.Marshal(posts)
    // 拿到json格式
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js)
    // 写回去
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeletePost(id, username); err != nil {
        http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete post from backend %v\n", err)
        return
    }
    fmt.Println("Post is deleted successfully")
}
