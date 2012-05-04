package pindex

import (
	"bytes"
	"github.com/peterbourgon/exp-html"
	"strings"
)

func SimpleCount(r Review) int { return 1 }

func baseWord(word string) string {
	return strings.ToLower(strings.Trim(word, ` ,.;!?"-`))
}

func WordCount(r Review) int {
	valid := 0
	for _, word := range strings.Split(r.Body, " ") {
		if baseWord(word) != "" {
			valid += 1
		}
	}
	return valid
}

func CharacterCount(r Review) int {
	return len(stripHTML(r.Body))
}

func InventedWordsFunc(dictfile string) func(r Review) int {
	dict := NewDict(dictfile)
	return func(r Review) int {
		count := 0
		for _, word := range strings.Split(r.Body, " ") {
			if !dict.Has(baseWord(word)) {
				// fmt.Printf("invented '%s'\n", baseWord(word))
				count += 1
			}
		}
		return count
	}
}

func NaïveSentenceLength(r Review) int {
	i, sentences := 0, 0
	for {
		j := strings.Index(r.Body[i:], ".")
		if j < 0 {
			break
		}
		sentences += 1
		i = i + j + 1
	}
	return int(float64(WordCount(r)) / float64(sentences))
}

func Pitchformulaity(r Review) int {
	score := 0
	for _, word := range strings.Split(r.Body, " ") {
		word = strings.Trim(strings.ToLower(word), ",.;!?")
		if n, ok := PitchformulaWords[word]; ok {
			score += n
		}
	}
	return score
}

func stripHTML(s string) string {
	z := html.NewTokenizer(bytes.NewBufferString(s))
	results := []string{}
	done := false
	for !done {
		tt := z.Next()
		switch tt {
		case html.TextToken:
			results = append(results, string(z.Text()))
		case html.ErrorToken:
			done = true
		}
	}
	return strings.Join(results, "")
}

var (
	PitchformulaWords = map[string]int{
		// word          triteness (higher=more trite)
		"warm":          1,
		"warmth":        1,
		"distorted":     1,
		"distortion":    1,
		"echoes":        1,
		"echoing":       1,
		"echo":          1,
		"organic":       5,
		"percussive":    1,
		"gentle":        1,
		"heavy":         1,
		"soft":          1,
		"lush":          9,
		"ethereal":      9,
		"delicate":      3,
		"plucked":       1,
		"buzzing":       1,
		"shimmering":    5,
		"fragile":       5,
		"swirling":      1,
		"cutting":       1,
		"understated":   2,
		"clouds":        1,
		"chiming":       1,
		"pounding":      1,
		"pulsing":       1,
		"fluid":         1,
		"rolling":       1,
		"skittering":    3,
		"rumbling":      1,
		"dreamy":        1,
		"hushed":        1,
		"backwards":     1,
		"murky":         1,
		"fuzzy":         1,
		"subtle":        2,
		"subtly":        2,
		"subtlety":      2,
		"layer":         1,
		"layers":        1,
		"layered":       1,
		"builds":        1,
		"swells":        1,
		"crescendos":    1,
		"rising":        1,
		"blast":         1,
		"blasts":        1,
		"crashing":      1,
		"explosive":     1,
		"complex":       1,
		"complexity":    1,
		"complicated":   1,
		"simple":        1,
		"massive":       1,
		"vast":          1,
		"expansive":     1,
		"tension":       1,
		"unexpected":    1,
		"unpredictable": 1,
		"chaos":         1,
		"chaotic":       1,
		"dense":         1,
		"structure":     1,
		"structured":    1,
		"abstract":      1,
		"accessible":    1,
		"detail":        1,
		"detailed":      1,
		"seamlessly":    1,
		"hypnotic":      4,
		"shifting":      1,
		"drifting":      1,
		"drops":         1,
		"sky":           1,
		"storm":         1,
		"storms":        1,
		"stormy":        1,
		"twists":        1,
		"dynamic":       1,
		"disparate":     2,
		"counterpoint":  1,
		"sweeping":      1,
		"surreal":       1,
		"dissonance":    1,
		"glowing":       3,
		"vibrant":       1,
		"controlled":    1,
		"faded":         1,
		"winds":         1,
		"skeletal":      5,
		"repeatedly":    1,
		"glow":          1,
		"spacious":      1,
		"ocean":         1,
		"oceans":        1,
		"oceanic":       1,
		"rough":         1,
		"primitive":     1,
		"lone":          1,
		"dominated":     1,
		"unstructured":  1,
		"rehearsed":     1,
		"polished":      1,
		"shiny":         1,
		"predictable":   1,
		"melancholy":    3,
		"sadness":       1,
		"plaintive":     1,
		"somber":        1,
		"dirge":         2,
		"depression":    1,
		"dark":          1,
		"frenetic":      2,
		"manic":         1,
		"frantic":       2,
		"frenzied":      1,
		"wild":          1,
		"madness":       1,
		"crazed":        1,
		"strange":       1,
		"mysterious":    2,
		"ghostly":       4,
		"violent":       1,
		"brutal":        1,
		"violence":      1,
		"aggression":    1,
		"joy":           1,
		"happy":         1,
		"bliss":         1,
		"sinister":      1,
		"ominous":       1,
		"menacing":      1,
		"frightening":   1,
		"tense":         1,
		"anxiety":       1,
		"anxious":       1,
		"restless":      1,
		"furious":       1,
		"fury":          1,
		"anger":         1,
		"emotional":     1,
		"playful":       1,
		"personal":      1,
		"affecting":     1,
		"assured":       1,
		"confidence":    1,
		"confident":     1,
		"romantic":      1,
		"relentless":    1,
		"despair":       2,
		"loss":          1,
		"regret":        1,
		"abandon":       1,
		"modest":        1,
		"dangerous":     1,
		"tortured":      1,
		"alienation":    2,
		"insecurity":    1,
		"whisper":       1,
		"whispers":      1,
		"whispering":    1,
		"croon":         5,
		"croons":        5,
		"crooning":      5,
		"wail":          3,
		"wails":         3,
		"wailing":       3,
		"tenor":         1,
		"chant":         1,
		"chants":        1,
		"chanting":      1,
		"baritone":      1,
		"soprano":       1,
		"alto":          1,
		"choir":         1,
		"scream":        1,
		"screams":       1,
		"screaming":     1,
		"off-key":       1,
		"yell":          1,
		"yelling":       1,
		"nasal":         1,
		"nasally":       1,
	}
)
