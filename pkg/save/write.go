package save

import (
	"encoding/binary"
	"io"
)

func writeProperty(w io.Writer, p Property) error {
	err := writeString(w, p.name)
	if err != nil {
		return err
	}
	err = writeString(w, p.valueType)
	if err != nil {
		return err
	}

	// compute length and value bytes
	var bs []byte
	switch p.valueType {
	case "StrProperty": // string property length is 14 + string length
		{
			bs = encodeString(p.Value.(string))
		}
	case "IntProperty":
		{
			bs = encodeInt32(p.Value.(int32))
		}
	case "FloatProperty":
		{
			bs = encodeFloat(p.Value.(float32))
		}
	case "BoolProperty":
		{
			bs = encodeBool(p.Value.(bool))
		}
	default:
		bs = encodeRaw(p.Value.([]byte))
	}

	return writeRaw(w, bs)
}

func writeRaw(w io.Writer, bs []byte) error {
	return binary.Write(w, binary.LittleEndian, bs)
}

func writeHeader(w io.Writer, header [metaOffset]byte) error {
	return binary.Write(w, binary.LittleEndian, header)
}

func writeFooter(w io.Writer, footer [eofOffset]byte) error {
	return binary.Write(w, binary.LittleEndian, footer)
}

func writeString(w io.Writer, s string) error {
	err := binary.Write(w, binary.LittleEndian, uint64(len(s)))
	if err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, []byte(s))
}

func encodeString(property string) []byte {
	// string property is length (8) + padding (1) + strlen (4) + value (n) + null byte (1) = 14 + n bytes
	// length is strlen (4) + value (n) + null byte (1) = 5 + n bytes
	bs := make([]byte, 14+len(property)) // slice is initialised with 0 bytes so we don't have to worry about padding
	binary.LittleEndian.PutUint64(bs[0:8], uint64(5+len(property)))
	binary.LittleEndian.PutUint32(bs[9:13], uint32(1+len(property)))
	for i, b := range []byte(property) {
		bs[13+i] = b
	}
	return bs
}

func encodeInt32(property int32) []byte {
	// int32 is length (8) + padding (1) + int32 (4)
	bs := make([]byte, 13)
	binary.LittleEndian.PutUint64(bs[0:8], 4)
	binary.LittleEndian.PutUint32(bs[9:13], uint32(property))
	return bs
}

func encodeFloat(property float32) []byte {
	// float32 is length (8) + padding (1) + int32 (4)
	bs := make([]byte, 13)
	binary.LittleEndian.PutUint64(bs[0:8], 4)
	binary.LittleEndian.PutUint32(bs[9:13], uint32(property))
	return bs
}

func encodeBool(property bool) []byte {
	// bool is length (8) + bool (2)
	bs := make([]byte, 10)
	binary.LittleEndian.PutUint64(bs[0:8], 0)
	if property {
		binary.LittleEndian.PutUint16(bs[8:10], 1)
	} else {
		binary.LittleEndian.PutUint16(bs[8:10], 0)
	}
	return bs
}

func encodeRaw(src []byte) []byte {
	bs := make([]byte, uint64(9+len(src)))
	binary.LittleEndian.PutUint64(bs[0:8], uint64(len(src)))
	for i, b := range src {
		bs[9+i] = b
	}
	return bs
}
