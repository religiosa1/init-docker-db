// Package randomname provides an utility function for generating a random
// adj-noun name for a dcker container
package randomname

import (
	"fmt"
	"math/rand"
)

/* List of tokens (adjectives and nouns) borrowed from the boring-name-generator
 * npm package.
 * https://github.com/boringprotocol/boring-name-generator
 */

// Generate a random name for docker container
func Generate() string {
	adjIdx := rand.Intn(len(adjectives))
	nounIdx := rand.Intn(len(nouns))
	return fmt.Sprintf("%s-%s", adjectives[adjIdx], nouns[nounIdx])
}
