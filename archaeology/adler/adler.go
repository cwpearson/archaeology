package adler

type Sum struct {
	a, b uint16

	k, l uint64
}

// NewSum produces a Sum of the buffer buf assuming buf is a window starting at offset k
func NewSum(buf []byte, k uint64) *Sum {

	if len(buf) == 0 {
		return nil
	}

	s := &Sum{}

	// Checksum of bytes k through l
	s.k = k
	s.l = s.k + uint64(len(buf)-1)

	for i, data := range buf {
		s.a += uint16(data)
		s.b += uint16((s.l - (uint64(i) + s.k) + 1) * uint64(data))
	}
	return s
}

func (s *Sum) Recurrence(add, sub byte) {

	s.a = uint16(uint32(s.a) - uint32(sub) + uint32(add))
	s.b = uint16(uint64(s.b) - (uint64(s.l-s.k+1)*uint64(sub) + uint64(s.a)))

	s.k++
	s.l++
}

func (s *Sum) Current() uint32 {
	return uint32(s.a) + (uint32(s.b) << 16)
}
