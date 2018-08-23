package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/julienschmidt/httprouter"
)
// path为 ‘/’的handler 
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "welcome!\n")
}
// 请求获取的实体
type Book struct {
    ISDN   string "'json':'isdn'"
    Title  string "'json':'title'"
    Author string "'json':'author'"
    Pages  int    "'json':''pages"
}
//响应时返回的数据 Data
type JsonResponse struct {
    Meta interface{} "json:'status'"
    Data interface{} "json:'data'"
}
//请求错误的响应处理
type JsonErrorResponse struct {
    Error *ApiError "json:'error'"
}

type ApiError struct {
    Status int16  "json:'status'"
    Title  string "json:'title'"
}

var bookstore = make(map[string]*Book)

// path为 ‘/books’的handler 

func BookIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    books := []*Book{}
    for _, book := range bookstore {
        books = append(books, book)
    }
    response := &JsonResponse{Data: &books}
    w.Header().Set("Content-Type", "application/json;charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        panic(err)
    }
}
// path为 ‘/books/：json’的handler 
func BookShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    isdn := params.ByName("isdn")
    book, ok := bookstore[isdn]
    w.Header().Set("Content-Type", "application/json;charset=UTF-8")
    if !ok {
        w.WriteHeader(http.StatusNotFound)
        response := JsonErrorResponse{Error: &ApiError{Status: 404, Title: "Record Not Found"}}
        if err := json.NewEncoder(w).Encode(response); err != nil {
            panic(err)
        }
    }
    response := JsonResponse{Data: book}
    if err := json.NewEncoder(w).Encode(response); err != nil {
        panic(err)
    }
}


func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/books/", BookIndex)
    router.GET("/books/:isdn", BookShow)
    //测试数据
    bookstore["123"] = &Book{
        ISDN:   "123",
        Title:  "Silence of the Lambs",
        Author: "Thomas Harris",
        Pages:  367,
    }
    //测试数据
    bookstore["124"] = &Book{
        ISDN:   "124",
        Title:  "tO KILL a mocking bird",
        Author: "Thomas Harris",
        Pages:  320,
    }
    log.Fatal(http.ListenAndServe(":8080", router))
}