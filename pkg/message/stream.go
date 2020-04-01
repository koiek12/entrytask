package message

import (
	"bufio"
	"encoding/binary"
	"io"
)

type MsgStream struct {
	in  *bufio.Reader
	out *bufio.Writer
	tmp []byte
}

func NewMsgStream(in io.Reader, out io.Writer) *MsgStream {
	return &MsgStream{bufio.NewReader(in), bufio.NewWriter(out), make([]byte, 32)}
}

func (st *MsgStream) readVint() uint {
	reqType, _ := binary.ReadUvarint(st.in)
	return uint(reqType)
}

func (st *MsgStream) readData() []byte {
	size, _ := binary.ReadUvarint(st.in)
	data := make([]byte, size)
	ptr := uint64(0)
	for ptr < size {
		n, _ := st.in.Read(data[ptr:len(data)])
		ptr += uint64(n)
	}
	return data
}

func (st *MsgStream) writeVint(reqType uint) {
	n := binary.PutUvarint(st.tmp, uint64(reqType))
	st.out.Write(st.tmp[:n])
}

func (st *MsgStream) writeData(data []byte) {
	n := binary.PutUvarint(st.tmp, uint64(len(data)))
	st.out.Write(st.tmp[:n])
	st.out.Write(data)
}

func (st *MsgStream) ReadMsg() (interface{}, error) {
	msgType := st.readVint()
	return ReadMsgFunc[msgType](st)
}

func (st *MsgStream) WriteMsg(msg interface{}) {
	msgNum := GetMsgNum(msg)
	WriteMsgFunc[msgNum](st, msg)
	st.out.Flush()
}
