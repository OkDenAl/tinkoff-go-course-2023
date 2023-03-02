package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	tags       []TagStat
	tagToIndex map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	return &TagCloud{
		tags:       make([]TagStat, 0),
		tagToIndex: make(map[string]int),
	}
}

func (tc *TagCloud) isTagInCloud(tag string) bool {
	_, status := tc.tagToIndex[tag]
	return status
}

func binsearch(arr []TagStat, target int) int {
	l := -1
	r := len(arr)
	for r-l > 1 {
		m := (l + r) / 2
		if arr[m].OccurrenceCount > arr[target].OccurrenceCount {
			l = m
		} else {
			r = m
		}
	}
	return r
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (tc *TagCloud) AddTag(tag string) {
	if tc.isTagInCloud(tag) {
		tc.tags[tc.tagToIndex[tag]].OccurrenceCount++
	} else {
		tc.tags = append(tc.tags, TagStat{Tag: tag, OccurrenceCount: 1})
		tc.tagToIndex[tag] = len(tc.tags) - 1
	}
	ind := binsearch(tc.tags, tc.tagToIndex[tag])
	indexTag := tc.tags[ind].Tag
	tc.tags[ind], tc.tags[tc.tagToIndex[tag]] = tc.tags[tc.tagToIndex[tag]], tc.tags[ind]
	tc.tagToIndex[indexTag], tc.tagToIndex[tag] = tc.tagToIndex[tag], tc.tagToIndex[indexTag]
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (tc *TagCloud) TopN(n int) []TagStat {
	if n > len(tc.tags) {
		return tc.tags
	} else {
		return tc.tags[:n]
	}
}
