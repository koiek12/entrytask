package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"git.garena.com/youngiek.song/backend/pkg/message"
	"github.com/dgrijalva/jwt-go"
)

var msp = message.NewMsgStreamPool("tcp", "localhost", "3233", 100)

func handleMain(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("../../web/template/login.html")
	t.Execute(w, nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")
	passwd := r.PostFormValue("pwd")

	stream, _ := msp.GetMsgStream()
	stream.WriteMsg(message.LoginRequest{id, passwd})
	res, err := stream.ReadMsg()
	if err != nil {
		stream.Destroy()
		fmt.Fprintln(w, err)
	}
	token, _ := jwt.Parse(res.(message.LoginResponse).Token, nil)
	claims := token.Claims.(jwt.MapClaims)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<html><body>")
	fmt.Fprintln(w, "code : ", res.(message.LoginResponse).Code)
	fmt.Fprintln(w, "token : ", res.(message.LoginResponse).Token)
	fmt.Fprintln(w, "claim : ", claims)
	fmt.Fprintln(w, "<img src=\"image/a.PNG\">")
	fmt.Fprintln(w, "</body></html>")
	stream.Close()
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()
	http.Handle("/image", http.FileServer(http.Dir(*directory)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
