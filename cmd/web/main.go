package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"git.garena.com/youngiek.song/backend/pkg/message"
	"github.com/dgrijalva/jwt-go"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3233"
	CONN_TYPE = "tcp"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("../../web/template/login.html")
	t.Execute(w, nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")
	passwd := r.PostFormValue("pwd")

	conn, err := net.Dial(CONN_TYPE, net.JoinHostPort(CONN_HOST, CONN_PORT))
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
	}

	// defer conn.Close()

	stream := message.NewMsgStream(conn, conn)
	stream.WriteMsg(message.LoginRequest{id, passwd})
	res, _ := stream.ReadMsg()
	token, _ := jwt.Parse(res.(message.LoginResponse).Token, nil)
	claims := token.Claims.(jwt.MapClaims)

	fmt.Fprintln(w, "code : ", res.(message.LoginResponse).Code)
	fmt.Fprintln(w, "token : ", res.(message.LoginResponse).Token)
	fmt.Fprintln(w, "claim : ", claims)
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
