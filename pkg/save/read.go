package save

import (
	"encoding/binary"
	"io"
)

func readString(r io.Reader) (string, error) {
	var length int32
	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}

	values := make([]byte, length)
	err = binary.Read(r, binary.LittleEndian, &values)
	if err != nil {
		return "", err
	}

	return string(values), nil
}

func readPaddedString(r io.Reader) (string, error) {
	_, err := readByte(r)
	if err != nil {
		return "", err
	}
	return readString(r)
}

func readByte(r io.Reader) (byte, error) {
	var b byte
	err := binary.Read(r, binary.LittleEndian, &b)
	return b, err
}

func readPaddedInt32(r io.Reader) (int32, error) {
	_, err := readByte(r)
	if err != nil {
		return 0, err
	}
	var i int32
	err = binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func readPaddedInt64(r io.Reader) (int64, error) {
	_, err := readByte(r)
	if err != nil {
		return 0, err
	}
	var i int64
	err = binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func readPaddedFloat(r io.Reader) (float32, error) {
	_, err := readByte(r)
	if err != nil {
		return 0, err
	}
	var i float32
	err = binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func readBool(r io.Reader) (bool, error) {
	var i int16
	err := binary.Read(r, binary.LittleEndian, &i)
	return i == 1, err
}
