package adler

type Adler struct {
	a, b uint16

	k, l uint64

	window []byte
}

// NewSum produces a Sum of the buffer buf assuming buf is a window starting at offset k
func NewSum(buf []byte) *Adler {

	s := &Adler{}
	s.window = buf
	// Checksum of bytes k through l
	s.k = 0
	s.l = s.k + uint64(len(buf))

	for i, data := range s.window {
		s.a += uint16(data)
		s.b += uint16((s.l - (uint64(i) + s.k)) * uint64(data))
	}
	return s
}

func Sum(buf []byte) uint32 {

	s := &Adler{}
	s.window = buf
	// Checksum of bytes k through l
	s.k = 0
	s.l = s.k + uint64(len(buf))

	for i, data := range s.window {
		s.a += uint16(data)
		s.b += uint16((s.l - (uint64(i) + s.k)) * uint64(data))
	}
	return s.Current()
}

func (s *Adler) Roll(add byte) uint32 {

	sub := s.window[0] // leaving the window

	c := uint64(s.l - s.k)

	s.a = uint16(uint32(s.a) - uint32(sub) + uint32(add))
	s.b = uint16(uint64(s.b) - (c * uint64(sub)) + uint64(s.a))

	s.k++
	s.l++

	// update the window
	s.window = append(s.window[1:], add)

	return s.Current()
}

func (s *Adler) Current() uint32 {
	return uint32(s.a) + (uint32(s.b) << 16)
}
