package save

import (
	"encoding/binary"
	"io"
)

func readProperty(r io.ReadSeeker) (Property, error) {
	propertyName, err := readString(r)
	if err != nil {
		return Property{}, err
	}
	propertyType, err := readString(r)
	if err != nil {
		return Property{}, err
	}
	propertyLength, err := readInt64(r)
	if err != nil {
		return Property{}, err
	}

	var value interface{}

	switch propertyType {
	case "StrProperty":
		{
			value, err = readPaddedString(r)
		}
	case "IntProperty":
		{
			value, err = readPaddedInt32(r)
		}
	case "FloatProperty":
		{
			value, err = readPaddedFloat(r)
		}
	case "BoolProperty":
		{
			value, err = readBool(r)
		}
	default:
		value, err = readPaddedRaw(r, propertyLength)
	}
	return Property{
		length:    propertyLength,
		name:      propertyName,
		valueType: propertyType,
		Value:     value,
	}, nil
}

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

func readInt64(r io.Reader) (int64, error) {
	var i int64
	err := binary.Read(r, binary.LittleEndian, &i)
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

func readPaddedRaw(r io.Reader, len int64) ([]byte, error) {
	_, err := readByte(r)
	if err != nil {
		return nil, err
	}
	bs := make([]byte, len)
	err = binary.Read(r, binary.LittleEndian, &bs)
	return bs, err
}

func readHeader(r io.Reader) ([metaOffset]byte, error) {
	var bs [metaOffset]byte
	err := binary.Read(r, binary.LittleEndian, &bs)
	return bs, err
}

func readFooter(r io.Reader) ([eofOffset]byte, error) {
	var bs [eofOffset]byte
	err := binary.Read(r, binary.LittleEndian, &bs)
	return bs, err
}
