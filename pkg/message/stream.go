package message

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type MsgStream struct {
	in  *bufio.Reader
	out *bufio.Writer
	tmp []byte
	msp *MsgStreamPool
}

func NewMsgStream(in io.Reader, out io.Writer) (*MsgStream, error) {
	return &MsgStream{bufio.NewReader(in), bufio.NewWriter(out), make([]byte, 32), nil}, nil
}

func NewMsgStreamForPool(msp *MsgStreamPool) (*MsgStream, error) {
	conn, err := net.Dial(msp.connType, net.JoinHostPort(msp.connHost, msp.connPort))
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return nil, err
	}
	return &MsgStream{bufio.NewReader(conn), bufio.NewWriter(conn), make([]byte, 32), msp}, nil
}

func (ms *MsgStream) Close() {
	if ms.msp == nil {
		return
	}
	ms.msp.pool <- ms
}

func (ms *MsgStream) Destroy() {
	if ms.msp == nil {
		return
	}
	ms.msp.size--
}

func (ms *MsgStream) readVint() (uint, error) {
	val, err := binary.ReadUvarint(ms.in)
	return uint(val), err
}

func (ms *MsgStream) readLenDelimData() ([]byte, error) {
	size, err := binary.ReadUvarint(ms.in)
	if err != nil {
		return nil, err
	}
	data := make([]byte, size)
	ptr := uint64(0)
	var n int
	for ptr < size {
		n, err = ms.in.Read(data[ptr:len(data)])
		if err != nil {
			return nil, err
		}
		ptr += uint64(n)
	}
	return data, nil
}

func (ms *MsgStream) writeVint(val uint) error {
	n := binary.PutUvarint(ms.tmp, uint64(val))
	_, err := ms.out.Write(ms.tmp[:n])
	return err
}

func (ms *MsgStream) writeLenDelimData(data []byte) error {
	n := binary.PutUvarint(ms.tmp, uint64(len(data)))
	_, err := ms.out.Write(ms.tmp[:n])
	if err != nil {
		return err
	}
	_, err = ms.out.Write(data)
	return err
}

func (ms *MsgStream) ReadMsg() (interface{}, error) {
	msgType, err := ms.readVint()
	if err != nil {
		return nil, err
	}
	return ReadMsgFunc[msgType](ms)
}

func (ms *MsgStream) WriteMsg(msg interface{}) error {
	msgNum, err := GetMsgNum(msg)
	if err != nil {
		return err
	}
	err = ms.writeVint(uint(msgNum))
	if err != nil {
		return err
	}
	err = WriteMsgFunc[msgNum](ms, msg)
	if err != nil {
		return err
	}
	err = ms.out.Flush()
	return err
}
