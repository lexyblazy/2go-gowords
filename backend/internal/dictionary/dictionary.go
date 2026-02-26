package dictionary

import (
	"bufio"
	// "fmt"
	"math/rand"
	"os"
	"strings"
)

type Dictionary struct {
	filePath      string
	uniqueWords   map[string]struct{}
	WordsByLength map[int][]Word
}

type Word struct {
	Text string
	Freq [26]uint8
}

func NewWord(text string) Word {
	return Word{
		Text: text,
		Freq: getFreq(text),
	}
}

func NewDictionary(filePath string) *Dictionary {

	// initialize the dictionary
	d := &Dictionary{
		filePath:      filePath,
		uniqueWords:   make(map[string]struct{}),
		WordsByLength: make(map[int][]Word),
	}

	err := d.loadFromFile(filePath)
	if err != nil {
		panic(err)
	}

	return d
}

func (d *Dictionary) IsValidWord(word string) bool {

	if len(word) < 3 {
		return false
	}

	for _, letter := range word {
		if letter < 'a' || letter > 'z' {
			return false
		}
	}

	return true
}

func (d *Dictionary) NormalizeWord(word string) string {
	return strings.TrimSpace(strings.ToLower(word))

}

func (d *Dictionary) loadFromFile(filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := d.NormalizeWord(scanner.Text())
		if d.IsValidWord(word) {
			d.uniqueWords[word] = struct{}{}
			d.WordsByLength[len(word)] = append(d.WordsByLength[len(word)], NewWord(word))

		}
	}

	return nil
}

func getFreq(word string) [26]uint8 {

	freq := [26]uint8{}
	for _, letter := range word {
		freq[letter-'a']++
	}
	return freq
}

func (d *Dictionary) fitsWithin(subset [26]uint8, superset [26]uint8) bool {
	for i := 0; i < 26; i++ {

		if subset[i] > superset[i] {
			return false
		}
	}
	return true
}

// merge words into a single frequency array
func mergeWords(words ...Word) [26]uint8 {
	var out [26]uint8

	for _, w := range words {
		for i := 0; i < 26; i++ {
			out[i] += w.Freq[i]
		}
	}

	return out
}

func (d *Dictionary) GenerateValidWords(words ...Word) map[string]struct{} {

	pool := mergeWords(words...)

	maxLength := 0
	results := make(map[string]struct{})

	for _, val := range pool {
		maxLength += int(val)
	}

	for length := 3; length <= maxLength; length++ {
		words := d.WordsByLength[length]
		for _, word := range words {
			if d.fitsWithin(word.Freq, pool) {
				results[word.Text] = struct{}{}
			}
		}
	}

	return results
}

func (d *Dictionary) CountDistinctCharacters(words ...Word) int {

	distinct := 0

	pool := mergeWords(words...)
	for _, val := range pool {
		if val > 0 {
			distinct++
		}
	}
	return distinct
}

func (d *Dictionary) GetExpansionWord(rng *rand.Rand, words ...Word) Word {

	basePool := mergeWords(words...)

	unusedLetters := []int{}
	existingLetters := []int{}

	for i, val := range basePool {
		if val == 0 {
			unusedLetters = append(unusedLetters, i)
		} else {
			existingLetters = append(existingLetters, i)
		}
	}

	// Safety check
	if len(unusedLetters) == 0 {
		panic("no unused letters available")
	}

	// First letter must be new
	firstIdx := unusedLetters[rng.Intn(len(unusedLetters))]

	// Second letter may be new or existing
	var secondIdx int
	if rng.Float64() < 0.5 && len(existingLetters) > 0 {
		secondIdx = existingLetters[rng.Intn(len(existingLetters))]
	} else {
		secondIdx = unusedLetters[rng.Intn(len(unusedLetters))]
	}

	letters := []rune{
		rune(firstIdx + 'a'),
		rune(secondIdx + 'a'),
	}

	return NewWord(string(letters))
}
