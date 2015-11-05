// easybuff
// 不要修改本文件, 每次消息有变动, 请手动生成本文件
// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)

package proto

import (
	. "github.com/toophy/login/help"
)

type ActorBase2 struct {
	Name  string    // string
	Age   int8      // int8
	Hp    int32     // int32
	Mp    int32     // int32
	Maxhp int32     // int32
	Maxmp int32     // int32
	Ox    ActorBase // ActorBase
}

func (this *ActorBase2) Read(s *Stream) {
	s.ReadString(this.Name)
	s.ReadInt8(this.Age)
	s.ReadInt32(this.Hp)
	s.ReadInt32(this.Mp)
	s.ReadInt32(this.Maxhp)
	s.ReadInt32(this.Maxmp)
	this.Ox.Read(s)
}

func (this *ActorBase2) Write(s *Stream) {
	s.WriteString(this.Name)
	s.WriteInt8(this.Age)
	s.WriteInt32(this.Hp)
	s.WriteInt32(this.Mp)
	s.WriteInt32(this.Maxhp)
	s.WriteInt32(this.Maxmp)
	this.Ox.Write(s)
}
