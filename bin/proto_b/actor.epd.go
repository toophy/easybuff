// easybuff
// 不要修改本文件, 每次消息有变动, 请手动生成本文件
// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)

package proto

import (
	. "github.com/toophy/login/help"
)

type ActorBase struct {
	Name   string // string
	Age    int8   // int8
	Hp     int32  // int24
	Mp     int32  // int32
	Maxhp  int32  // int24
	Maxmp  int32  // int32
	Exp    int64  // int48
	MaxExp int64  // int56
}

func (this *ActorBase) Read(s *Stream) {
	s.ReadString(this.Name)
	s.ReadInt8(this.Age)
	s.ReadInt24(this.Hp)
	s.ReadInt32(this.Mp)
	s.ReadInt24(this.Maxhp)
	s.ReadInt32(this.Maxmp)
	s.ReadInt48(this.Exp)
	s.ReadInt56(this.MaxExp)
}

func (this *ActorBase) Write(s *Stream) {
	s.WriteString(this.Name)
	s.WriteInt8(this.Age)
	s.WriteInt24(this.Hp)
	s.WriteInt32(this.Mp)
	s.WriteInt24(this.Maxhp)
	s.WriteInt32(this.Maxmp)
	s.WriteInt48(this.Exp)
	s.WriteInt56(this.MaxExp)
}
