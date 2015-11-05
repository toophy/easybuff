package help

// import ()

// type Stream struct {
// 	Data   []byte
// 	MaxLen int32
// 	Pos    int32
// }

// func (t *Stream) Init(d []byte, d_len int32) {
// 	t.Data = d
// 	t.MaxLen = d_len
// 	t.Pos = 0
// }

// func (t *Stream) Seek(p int32) {
// 	if t.Data != nil && p >= 0 && p < t.MaxLen {
// 		t.Pos = p
// 	}
// }

// func (t *Stream) ReadInt8() int64 {
// 	if (t.Pos + 1) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 1
// 		return int64(t.Data[old_pos])
// 	}
// 	return 0
// }

// func (t *Stream) ReadInt16() int64 {
// 	if (t.Pos + 2) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 2
// 		return int64(t.Data[old_pos])<<8 + int64(t.Data[old_pos+1])
// 	}
// 	return 0
// }

// func (t *Stream) ReadInt24() int64 {
// 	if (t.Pos + 3) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 3
// 		return (int64(t.Data[old_pos]) << 16) +
// 			(int64(t.Data[old_pos+1]) << 8) +
// 			(int64(t.Data[old_pos+2]))
// 	}
// 	return 0
// }

// func (t *Stream) ReadInt32() int64 {
// 	if (t.Pos + 4) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 4
// 		return (int64(t.Data[old_pos]) << 24) +
// 			(int64(t.Data[old_pos+1]) << 16) +
// 			(int64(t.Data[old_pos+2]) << 8) +
// 			(int64(t.Data[old_pos+3]))
// 	}
// 	return 0
// }

// func (t *Stream) ReadInt40() int64 {
// 	if (t.Pos + 5) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 5
// 		return (int64(t.Data[old_pos]) << 32) +
// 			(int64(t.Data[old_pos+1]) << 24) +
// 			(int64(t.Data[old_pos+2]) << 16) +
// 			(int64(t.Data[old_pos+3]) << 8) +
// 			(int64(t.Data[old_pos+4]))
// 	}
// 	return 0
// }

// func (t *Stream) ReadInt48() int64 {
// 	if (t.Pos + 6) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 6
// 		return (int64(t.Data[old_pos]) << 40) +
// 			(int64(t.Data[old_pos+1]) << 32) +
// 			(int64(t.Data[old_pos+2]) << 24) +
// 			(int64(t.Data[old_pos+3]) << 16) +
// 			(int64(t.Data[old_pos+4]) << 8) +
// 			(int64(t.Data[old_pos+5]))
// 	}
// 	return 0
// }

// func (t *Stream) ReadInt56() int64 {
// 	if (t.Pos + 7) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 7
// 		return (int64(t.Data[old_pos]) << 48) +
// 			(int64(t.Data[old_pos+1]) << 40) +
// 			(int64(t.Data[old_pos+2]) << 32) +
// 			(int64(t.Data[old_pos+3]) << 24) +
// 			(int64(t.Data[old_pos+4]) << 16) +
// 			(int64(t.Data[old_pos+5]) << 8) +
// 			(int64(t.Data[old_pos+6]))
// 	}
// 	return 0
// }

// func (t *Stream) ReadInt64() int64 {
// 	if (t.Pos + 8) < (t.MaxLen + 1) {
// 		old_pos := t.Pos
// 		t.Pos = t.Pos + 8
// 		return (int64(t.Data[old_pos]) << 56) +
// 			(int64(t.Data[old_pos]) << 48) +
// 			(int64(t.Data[old_pos+1]) << 40) +
// 			(int64(t.Data[old_pos+2]) << 32) +
// 			(int64(t.Data[old_pos+3]) << 24) +
// 			(int64(t.Data[old_pos+4]) << 16) +
// 			(int64(t.Data[old_pos+5]) << 8) +
// 			(int64(t.Data[old_pos+6]))
// 	}
// 	return 0
// }

// func (t *Ty_msg_stream) ReadStr() string {
// 	data_len := t.ReadU2()
// 	if data_len > 0 && (t.pos+data_len) < (t.msg.Len+1) {
// 		old_pos := t.pos
// 		t.pos = t.pos + data_len
// 		return string(t.msg.Data[old_pos : old_pos+data_len])
// 	}
// 	return ""
// }

// func (t *Ty_msg_stream) WriteU1(d int) bool {
// 	if t.pos+1 < MaxDataLen {
// 		t.msg.Data[t.pos] = byte(d & 0xFF)
// 		t.pos = t.pos + 1
// 		t.msg.Len = t.msg.Len + 1
// 		return true
// 	}

// 	return false
// }

// func (t *Ty_msg_stream) WriteU2(d int) bool {
// 	if t.pos+2 < MaxDataLen {
// 		// 65280
// 		t.msg.Data[t.pos] = byte((d & 0xFF00) >> 8)
// 		//
// 		t.msg.Data[t.pos+1] = byte(d & 0xFF)
// 		t.pos = t.pos + 2
// 		t.msg.Len = t.msg.Len + 2
// 		return true
// 	}

// 	return false
// }

// func (t *Ty_msg_stream) WriteU4(d int) bool {
// 	nd := uint(d)
// 	if t.pos+4 < MaxDataLen {
// 		// 4278190080
// 		t.msg.Data[t.pos] = byte((nd & 0xFF000000) >> 24)
// 		// 16711680
// 		t.msg.Data[t.pos+1] = byte((nd & 0xFF0000) >> 16)
// 		// 65280
// 		t.msg.Data[t.pos+2] = byte((nd & 0xFF00) >> 8)
// 		//
// 		t.msg.Data[t.pos+3] = byte(nd & 0xFF)
// 		t.pos = t.pos + 4
// 		t.msg.Len = t.msg.Len + 4
// 		return true
// 	}

// 	return false
// }

// func (t *Ty_msg_stream) WriteString(d *string) bool {
// 	d_len := len(*d)

// 	if t.pos+2+d_len < MaxDataLen {
// 		if t.WriteU2(d_len) {
// 			ds := (*d)[:]
// 			dx := t.msg.Data[t.pos : t.pos+d_len]
// 			copy(dx, ds)
// 			t.pos = t.pos + d_len
// 			t.msg.Len = t.msg.Len + d_len
// 			return true
// 		}
// 	}

// 	//	println("write string too long")

// 	// if t.WriteU2(d_len) {

// 	// 	if t.pos+d_len < MaxDataLen {
// 	// 		ds := (*d)[:]
// 	// 		dx := t.msg.Data[t.pos : t.pos+d_len]
// 	// 		copy(dx, ds)
// 	// 		t.pos = t.pos + d_len
// 	// 		t.msg.Len = t.msg.Len + d_len
// 	// 		return true
// 	// 	} else {
// 	// 		t.pos = t.pos - 2
// 	// 		t.msg.Len = t.msg.Len - 2
// 	// 	}
// 	// }

// 	return false
// }
