package vanity

import (
	"math"
	"strings"
)

const halfAddressLength = 45

var (
	minPriv = &[32]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	maxPriv = &[32]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}

	minKey = &Key{minPriv, minPriv}
	maxKey = &Key{maxPriv, maxPriv}
)

func init() {
	scReduce32(minPriv)
	scReduce32(maxPriv)
}

func IsValidPrefix(prefix string, network []byte, initIndex int) bool {
	// check length
	if len(prefix) == 0 || len(prefix) >= 95-initIndex {
		return false
	}

	// check base58
	for _, v := range prefix {
		switch {
		case v < '1':
			return false
		case v > '9' && v < 'A':
			return false
		case v > 'Z' && v < 'a':
			return false
		case v > 'z':
			return false
		}

		switch v {
		case 'I', 'O', 'l':
			return false
		}
	}

	if initIndex < 2 {
		// check elliptic curve property
		minAddress := minKey.Address(network)
		maxAddress := maxKey.Address(network)

		if strings.Compare(prefix, minAddress[initIndex:]) < 0 {
			return false
		}
		if strings.Compare(prefix, maxAddress[initIndex:]) > 0 {
			return false
		}
	}

	return true
}

func NeedOnlySpendKey(prefix string) bool {
	return len(prefix) <= halfAddressLength-3
}

func EstimatedDifficulty(prefix string, caseSensitive, includeNetwork bool) uint64 {
	diff := uint64(0)

	// len(prefix) is byte count
	// origLen is bit count
	origLen := int(math.Round(float64(len(prefix)) * 5.858))
	if includeNetwork {
		origLen -= 8
	}
	for i := 0; i < origLen; i++ {
		diff |= 1 << uint(i)
	}

	if !caseSensitive {
		for _, v := range prefix {
			if 'a' <= v && v <= 'z' || 'A' <= v && v <= 'Z' {
				diff >>= 1
			}
		}
	}

	if diff < 1 {
		diff = 1
	}

	return diff
}
