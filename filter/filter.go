package filter

type TagMap map[string][][]byte

type F struct {
	Ids          [][]byte
	Authors      [][]byte
	Tags         TagMap
	Since, Until int64
	Sort         string
}
