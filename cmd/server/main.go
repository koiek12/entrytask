package main

import (
	"fmt"
	"net"
	"os"

	"git.garena.com/youngiek.song/backend/pkg/message"
	"github.com/dgrijalva/jwt-go"
)

const (
	ConnHost = "localhost"
	ConnPort = "3233"
	ConnType = "tcp"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	stream, _ := message.NewMsgStream(conn, conn)
	for {
		req, err := stream.ReadMsg()
		if err != nil {
			break
		}
		switch req.(type) {
		case message.LoginRequest:
			handleLoginRequest(req.(message.LoginRequest), stream)
		}
	}
}

func generateToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("young"))
	if err != nil {
		fmt.Println("Error generating token:", err)
		return ""
	}
	return tokenString
}

func authenticate(id, password string) bool {
	fmt.Println("authenticate successful : ", id, password)
	return true
}

func handleLoginRequest(req message.LoginRequest, st *message.MsgStream) {
	fmt.Println("login request : ", req)
	id := req.Id
	password := req.Password
	if !authenticate(id, password) {
		fmt.Println("authenticate failed")
		return
	}
	token := generateToken(id)
	st.WriteMsg(message.LoginResponse{uint(0), token})
	fmt.Println("token generated : ", token)
}

func main() {
	l, err := net.Listen(ConnType, net.JoinHostPort(ConnHost, ConnPort))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		fmt.Println("Listening new request..")
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
