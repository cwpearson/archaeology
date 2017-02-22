package adler

type Sum struct {
	a, b uint16

	k, l uint64

	window []byte
}

// NewSum produces a Sum of the buffer buf assuming buf is a window starting at offset k
func NewSum(buf []byte) *Sum {

	if len(buf) == 0 {
		return nil
	}

	s := &Sum{}
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

func (s *Sum) Roll(add byte) {

	sub := s.window[0] // leaving the window

	c := uint64(s.l - s.k)

	s.a = uint16(uint32(s.a) - uint32(sub) + uint32(add))
	s.b = uint16(uint64(s.b) - (c * uint64(sub)) + uint64(s.a))

	s.k++
	s.l++

	// update the window
	s.window = append(s.window[1:], add)
}

func (s *Sum) Current() uint32 {
	return uint32(s.a) + (uint32(s.b) << 16)
}
