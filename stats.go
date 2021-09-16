package dupi

import "fmt"

type Stats struct {
	Root      string
	NumDocs   uint64
	NumPaths  uint64
	NumPosts  uint64
	NumBlots  uint64
	BlotMean  float64
	BlotSigma float64
}

const stFmt = `dupi index at %s:
	- %d docs
	- %d nodes in path tree
	- %d posts
	- %d blots
	- %.2f mean docs per blot
	- %.2f sigma (std deviation)
`

func (st *Stats) String() string {
	return fmt.Sprintf(stFmt, st.Root, st.NumDocs,
		st.NumPaths, st.NumPosts, st.NumBlots,
		st.BlotMean, st.BlotSigma)
}
