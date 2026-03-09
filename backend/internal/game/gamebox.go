package game

import (
	"math/rand"
	// "fmt"
	// "strings"

	"github.com/lexyblazy/gowords/internal/dictionary"
)

type Gamebox struct {
	rng *rand.Rand
	d   *dictionary.Dictionary
}

func NewGamebox(dictionary *dictionary.Dictionary, rng *rand.Rand) *Gamebox {
	return &Gamebox{
		rng: rng,
		d:   dictionary,
	}
}

func (b *Gamebox) GetWordsForRound(wordLength int, count int) []dictionary.Word {

	words := b.d.WordsByLength[wordLength]
	selectedWords := make([]dictionary.Word, count)
	for i := 0; i < count; i++ {
		selectedWords[i] = words[b.rng.Intn(len(words))]
	}

	return selectedWords
}

func (b *Gamebox) GetExpansionWord(words ...dictionary.Word) dictionary.Word {
	return b.d.GetExpansionWord(b.rng, words...)
}

func (b *Gamebox) GetValidWordCombinationsWithExpansion(expansionWord dictionary.Word, words ...dictionary.Word) map[string]struct{} {

	updatedWords := append(words, expansionWord)

	return b.GetValidWordCombinations(updatedWords...)

}

func (b *Gamebox) GetValidWordCombinations(words ...dictionary.Word) map[string]struct{} {
	return b.d.GenerateValidWords(words...)
}

func (b *Gamebox) GetDistinctCharacterCount(words ...dictionary.Word) int {
	return b.d.CountDistinctCharacters(words...)
}
