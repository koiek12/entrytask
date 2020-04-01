package message

import "reflect"

var msgNums = map[string]int{
	"message.LoginRequest":  0,
	"message.LoginResponse": 1,
}

func GetMsgNum(msg interface{}) int {
	msgType := reflect.TypeOf(msg).String()
	return msgNums[msgType]
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
	id := string(st.readLenDelimData())
	password := string(st.readLenDelimData())
	return LoginRequest{id, password}, nil
}

func WriteLoginRequest(st *MsgStream, req interface{}) error {
	st.writeLenDelimData([]byte(req.(LoginRequest).Id))
	st.writeLenDelimData([]byte(req.(LoginRequest).Password))
	return nil
}

type LoginResponse struct {
	Code  uint
	Token string
}

func ReadLoginResponse(st *MsgStream) (interface{}, error) {
	code := st.readVint()
	token := st.readLenDelimData()
	return LoginResponse{code, string(token)}, nil
}

func WriteLoginResponse(st *MsgStream, req interface{}) error {
	st.writeVint(req.(LoginResponse).Code)
	st.writeLenDelimData([]byte(req.(LoginResponse).Token))
	return nil
}
