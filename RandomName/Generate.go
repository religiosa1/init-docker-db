package RandomName

import (
	"fmt"
	"math/rand"
)

/* List of tokens (adjectives and nouns) borrowed from the boring-name-generator
 * npm package.
 * https://github.com/boringprotocol/boring-name-generator
 */

func Generate() string {
	adjIdx := rand.Intn(len(adjectives))
	nounIdx := rand.Intn(len(nouns))
	return fmt.Sprintf("%s-%s", adjectives[adjIdx], nouns[nounIdx])
}
