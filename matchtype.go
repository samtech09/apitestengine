package apitestengine

type MatchType int

const (
	MatchExact MatchType = iota
	MatchStartsWith
	MatchEndsWith
	MatchContains
	MatchNotContains
)

// type IMatchType interface {
// 	MatchType() matchType
// }

// // every base must fulfill the Baser interface
// func (b matchType) MatchType() matchType {
// 	return b
// }
