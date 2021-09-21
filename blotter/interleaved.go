package blotter

type Interleaved struct {
	circs []*Circ
	i     int
}

func NewInterleaved(seqLen, interLeaving int) *Interleaved {
	res := &Interleaved{circs: make([]*Circ, interLeaving)}
	for i := range res.circs {
		res.circs[i] = NewCirc(seqLen)
	}
	return res
}

func (i *Interleaved) Config() *Config {
	example := i.circs[0]
	seqLen := len(example.hashes)
	return &Config{SeqLen: seqLen, Interleave: len(i.circs)}
}

func (i *Interleaved) Interleaving() int {
	return len(i.circs)
}

func (i *Interleaved) Blot(tok []byte) uint32 {
	res := i.circs[i.i].Blot(tok)
	i.i++
	if i.i == len(i.circs) {
		i.i = 0
	}
	return res
}
