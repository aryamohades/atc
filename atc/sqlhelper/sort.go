package sqlhelper

import "strings"

const (
	SortDirAsc  = "ASC"
	SortDirDesc = "DESC"
)

var sortMap = map[string]string{
	"asc":  SortDirAsc,
	"desc": SortDirDesc,
}

func GetSortDir(sort string) (string, bool) {
	v, ok := sortMap[strings.ToLower(sort)]
	return v, ok
}
