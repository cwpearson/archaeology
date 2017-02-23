package archaeology

type InstrTy int

const (
	BlockRef InstrTy = iota
	NewData
)

type Instruction struct {
	Ty       InstrTy
	blockRef int64
	newData  []byte
}

func (i *Instruction) String() string {
	switch i.Ty {
	case BlockRef:
		return "Block=" + string(i.blockRef)
	case NewData:
		return "Data len=" + string(len(i.newData))
	default:
		return "UNEXPECTED"
	}
}

func NewBlockRef(block int64) *Instruction {
	return &Instruction{BlockRef, block, nil}
}

func NewNewData(data []byte) *Instruction {
	return &Instruction{NewData, -1, data}
}
