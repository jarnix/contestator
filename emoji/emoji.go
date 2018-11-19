package emoji

import (
	"math/rand"
	"os"
	"time"

	"github.com/kyokomi/emoji"
)

var emojis = [...]string{
	":+1:",
	":beer:",
	":burrito:",
	":cocktail:",
	":computer:",
	":cookie:",
	":cool:",
	":dollar:",
	":duck:",
	":frowning_face:",
	":game_die:",
	":ghost:",
	":heart:",
	":kiss:",
	":lollipop:",
	":mask:",
	":peach:",
	":poop:",
	":potato:",
	":roll_eyes:",
	":eggplant:",
}

// GenerateText generates a series of emojis from the list
func GenerateText() string {
	rand.Seed(time.Now().Unix())
	var stupidText = ""
	p := rand.Perm(len(emojis))
	var randomInts []int
	for _, randomInt := range p[:3] {
		randomInts = append(randomInts, randomInt)
	}
	stupidText += emojis[randomInts[0]] + " + "
	stupidText += emojis[randomInts[1]] + " = "
	stupidText += emojis[randomInts[2]]
	stupidText += os.Getenv("HASHTAGS_EMOJIS")
	return emoji.Sprint(stupidText)
}
