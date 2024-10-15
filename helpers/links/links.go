package links

import (
	"cmp"
	"slices"

	"github.com/pakkerman/web-crawler-go/types"
)

func SortLinks(input map[string]int) []types.SortedLinks {
	var sortedSlice []types.SortedLinks

	for k, v := range input {
		sortedSlice = append(sortedSlice, types.SortedLinks{k, v})
	}

	slices.SortFunc(sortedSlice, func(a, b types.SortedLinks) int {
		return cmp.Or(
			cmp.Compare(b.Visits, a.Visits),
			cmp.Compare(a.Url, b.Url),
		)
	})

	return sortedSlice
}
