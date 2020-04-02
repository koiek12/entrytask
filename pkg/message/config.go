package message

import "reflect"

var msgNums = map[string]int{
	"message.LoginRequest":  0,
	"message.LoginResponse": 1,
}

func GetMsgNum(msg interface{}) (int, error) {
	msgType := reflect.TypeOf(msg).String()
	num, ok := msgNums[msgType]
	if !ok {
		return 0, nil
	}
	return num, nil
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
	id, err := st.readLenDelimData()
	if err != nil {
		return nil, err
	}
	password, err := st.readLenDelimData()
	return LoginRequest{string(id), string(password)}, err
}

func WriteLoginRequest(st *MsgStream, req interface{}) error {
	err := st.writeLenDelimData([]byte(req.(LoginRequest).Id))
	if err != nil {
		return err
	}
	err = st.writeLenDelimData([]byte(req.(LoginRequest).Password))
	return err
}

type LoginResponse struct {
	Code  uint
	Token string
}

func ReadLoginResponse(st *MsgStream) (interface{}, error) {
	code, err := st.readVint()
	if err != nil {
		return nil, err
	}
	token, err := st.readLenDelimData()
	return LoginResponse{code, string(token)}, err
}

func WriteLoginResponse(st *MsgStream, req interface{}) error {
	err := st.writeVint(req.(LoginResponse).Code)
	if err != nil {
		return err
	}
	err = st.writeLenDelimData([]byte(req.(LoginResponse).Token))
	return err
}
