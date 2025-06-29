package filter

type TagMap map[string][][]byte

type F struct {
	Ids          [][]byte
	NotIds       [][]byte
	Authors      [][]byte
	NotAuthors   [][]byte
	Tags         TagMap
	NotTags      TagMap
	Since, Until int64
	Sort         string
}
