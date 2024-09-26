package main

import (
	"cmp"
	"slices"
)

func sortLinks(input map[string]int) []sortedLinks {
	var sortedSlice []sortedLinks

	for k, v := range input {
		sortedSlice = append(sortedSlice, sortedLinks{k, v})
	}

	slices.SortFunc(sortedSlice, func(a, b sortedLinks) int {
		return cmp.Or(
			cmp.Compare(b.visits, a.visits),
			cmp.Compare(a.url, b.url),
		)
	})

	return sortedSlice
}
