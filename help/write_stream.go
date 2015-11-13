package help

import (
	"errors"
	"math"
)

type Alloc struct {
}

func (this *Alloc) New(s int) []byte {
	b := make([]byte, s)
	return b
}

func (this *Alloc) Del(b []byte) {
}

const (
	MsgHeaderSize    = 4
	PacketHeaderSize = 4
	PacketSize       = 4096
)

type WriteStream struct {
	Stream
	Allocor    *Alloc   // 内存分配器
	Packets    [][]byte // 包
	CurrPacket int32    // 当前包
	LastMsgPos uint64   // 最后一个完整写入的消息位置, LastMsgPos<=Pos
	MsgCount   uint32   // 消息数量
	Writting   bool     // 正在写消息中
}

func (t *WriteStream) Init(d []byte) {
	t.Data = d
	t.MaxLen = uint64(len(t.Data))
	t.CurrPos = 0
	t.LastMsgPos = 0
	t.Writting = false
	t.MsgCount = 0
	t.Packets = make([][]byte, 10)
	t.CurrPacket = -1
}

func (t *WriteStream) packetBegin() {
	t.CurrPacket++
	t.Packets[t.CurrPacket] = t.Allocor.New(PacketSize)

	t.Data = t.Packets[t.CurrPacket]
	t.MaxLen = uint64(len(t.Data))
	t.CurrPos = 0
	t.LastMsgPos = 0
	t.Writting = false
	t.MsgCount = 0

	t.Seek(0)
	t.WriteUint32(0)
	t.LastMsgPos = t.CurrPos
}

func (t *WriteStream) packetEnd() {
	// LastMsgPos
	// MsgCount
	// Key
	// 写入包头
}

func (t *WriteStream) WriteBegin(id uint16) {
	if !t.Writting {
		t.Writting = true
		t.WriteUint16(0)
		t.WriteUint16(id)
	} else {
		panic(errors.New("WriteStream:WriteBegin no end"))
	}
}

func (t *WriteStream) WriteEnd() {
	if t.Writting {
		currPos := t.CurrPos
		s.Seek(t.LastMsgPos)
		t.WriteUint16(currPos - t.LastMsgPos)
		s.Seek(currPos)
		t.Writting = false
		t.MsgCount++
		t.LastMsgPos = currPos
	} else {
		panic(errors.New("WriteStream:WriteEnd no begin"))
	}
}

func (t *WriteStream) MoveFailedMsg(tx *WriteStream) {
	if t.Writting {
		// 这个消息需要挪动
		// 挪动完成
		// t.Writting => false
		// 可是, 消息write怎么继续?
	}
}
