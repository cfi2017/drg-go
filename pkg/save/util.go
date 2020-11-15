package save

import "io"

func setPos(s io.Seeker, offset int64) (int64, error) {
	return s.Seek(offset, io.SeekCurrent)
}

func getPos(s io.Seeker) int64 {
	pos, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}
	return pos
}
