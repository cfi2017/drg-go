package save

import (
	"os"
)

const (
	metaOffset = 1037
	eofOffset  = 13
)

type Property struct {
	Value     interface{}
	length    int64
	offset    int64
	valueType string
	name      string
}

type SaveFile struct {
	Path       string
	properties map[string]Property
	header     [metaOffset]byte
	footer     [eofOffset]byte
}

func (s *SaveFile) ReadProperties() error {
	file, err := os.Open(s.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	length := info.Size()
	s.header, err = readHeader(file)
	if err != nil {
		return err
	}
	_, err = readString(file)
	if err != nil {
		return err
	}

	s.properties = make(map[string]Property)
	for pos := getPos(file); pos < length-eofOffset; {
		property, err := readProperty(file)
		if err != nil {
			return err
		}
		if property.name == "" {
			continue
		}
		property.offset = pos
		s.properties[property.name] = property
	}

	s.footer, err = readFooter(file)
	return err
}

func (s *SaveFile) SaveProperties() error {
	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	err = writeHeader(file, s.header)
	if err != nil {
		return err
	}

	for _, property := range s.properties {
		err = writeProperty(file, property)
		if err != nil {
			return err
		}
	}

	err = writeFooter(file, s.footer)
	if err != nil {
		return err
	}

	panic("implement me")
}
