package message

import "fmt"

type MsgStreamPool struct {
	pool                         chan *MsgStream
	size                         int
	cap                          int
	connType, connHost, connPort string
}

func NewMsgStreamPool(connType, connHost, connPort string, cap int) *MsgStreamPool {
	return &MsgStreamPool{make(chan *MsgStream, cap), 0, cap, connType, connHost, connPort}
}

func (msp *MsgStreamPool) GetMsgStream() (*MsgStream, error) {
	if msp.size < msp.cap {
		fmt.Println("create new stream")
		stream, err := NewMsgStreamForPool(msp)
		if err != nil {
			return nil, err
		}
		msp.pool <- stream
		msp.size++
		fmt.Println("pool size :", msp.size)
	}
	return <-msp.pool, nil
}
