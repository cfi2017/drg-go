package save

import (
	"encoding/binary"
	"errors"
	"os"
)

const (
	metaOffset = 893
	eofOffset  = 13
)

var (
	ErrPropertyNotFound = errors.New("property not found")
)

type SaveFile struct {
	Path string
}

func (s *SaveFile) GetValue(key string) (interface{}, error) {
	file, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	length := info.Size()
	_, err = setPos(file, metaOffset)
	if err != nil {
		return nil, err
	}
	_, err = readString(file)
	if err != nil {
		return nil, err
	}

	for getPos(file) < length-eofOffset {
		propertyName, err := readString(file)
		if err != nil {
			return nil, err
		}
		propertyType, err := readString(file)
		if err != nil {
			return nil, err
		}
		var propertyLength int64
		err = binary.Read(file, binary.LittleEndian, &propertyLength)
		if err != nil {
			return nil, err
		}

		var value interface{}

		switch propertyType {
		case "StrProperty":
			{
				value, err = readPaddedString(file)
			}
		case "IntProperty":
			{
				value, err = readPaddedInt32(file)
			}
		case "FloatProperty":
			{
				value, err = readPaddedFloat(file)
			}
		case "BoolProperty":
			{
				value, err = readBool(file)
			}
		default:
			_, err = setPos(file, propertyLength+1)
			if err != nil {
				return nil, err
			}
		}

		if propertyName == key {
			return value, nil
		}

	}

	return nil, ErrPropertyNotFound
}
