package message

import "reflect"

var MsgNums = map[string]int{
	"message.LoginRequest":  0,
	"message.LoginResponse": 1,
}

func GetMsgNum(msg interface{}) int {
	msgType := reflect.TypeOf(msg).String()
	return MsgNums[msgType]
}

var ReadMsgFunc = []func(*MsgStream) (interface{}, error){
	ReadLoginRequest,
	ReadLoginResponse,
}

var WriteMsgFunc = []func(*MsgStream, interface{}) error{
	WriteLoginRequest,
	WriteLoginResponse,
}

type LoginRequest struct {
	Id       string
	Password string
}

func ReadLoginRequest(st *MsgStream) (interface{}, error) {
	id := string(st.readData())
	password := string(st.readData())
	return LoginRequest{id, password}, nil
}

func WriteLoginRequest(st *MsgStream, req interface{}) error {
	st.writeVint(uint(0))
	st.writeData([]byte(req.(LoginRequest).Id))
	st.writeData([]byte(req.(LoginRequest).Password))
	return nil
}

type LoginResponse struct {
	Code  uint
	Token string
}

func ReadLoginResponse(st *MsgStream) (interface{}, error) {
	code := st.readVint()
	token := st.readData()
	return LoginResponse{code, string(token)}, nil
}

func WriteLoginResponse(st *MsgStream, req interface{}) error {
	st.writeVint(uint(1))
	st.writeVint(req.(LoginResponse).Code)
	st.writeData([]byte(req.(LoginResponse).Token))
	return nil
}
