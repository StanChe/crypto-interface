package script

import (
	"bytes"
	"log"
	"sort"
)

// implement sort package Interface
type sortedPubkeys [][]byte

func (pks sortedPubkeys) Len() int {
	return len(pks)
}

func (pks sortedPubkeys) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(pks[i], pks[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (pks sortedPubkeys) Swap(i, j int) {
	pks[j], pks[i] = pks[i], pks[j]
}

// SortedPubkeys - exported function to get pubkeys sorted
func SortedPubkeys(src [][]byte) [][]byte {
	sorted := sortedPubkeys(src)
	sort.Sort(sorted)
	return sorted
}
