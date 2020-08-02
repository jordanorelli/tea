package tea

// generator functions for creating randomized test data.

import (
	"math/rand"
	"time"
)

// alpha is an alphabet of letters to pick from when producing random strings.
// the selected characters are human-readable without much visual ambiguity and
// include some code points beyond the ascii range to make sure things don't
// break on unicode input.
var alpha = []rune("" +
	// lower-case ascii letters
	"abcdefghjkmnpqrstwxyz" +
	//  upper-case ascii letters
	"ABCDEFGHIJKLMNPQRSTWXYZ" +
	// some digits
	"23456789" +
	// miscellaneous non-ascii characters
	"¢£¥ÐÑØæñþÆŁřƩλЖд")

func rstring(n int) string {
	r := make([]rune, n)
	for i, _ := range r {
		r[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(r)
}

func init() {
	rand.Seed(time.Now().Unix())
}
