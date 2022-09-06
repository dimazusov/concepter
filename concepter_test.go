package main

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"optimization/internal/pkg/morph"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"optimization/internal/pkg/sentence"
)

func TestNewConcepterAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// givenSentence    перемести - глагол в повелительном наклонении

	partStr := "глагол в повелительном наклонении"
	fullStr := "необходимо выполнить команду для глагола в повелительном наклонении"
	findStr := "перемести глагол в повелительном наклонении"
	expectedStr := "необходимо выполнить команду для перемещения"

	partSentence := getSentence(partStr).Sentence[partStr]
	fullSentence := getSentence(fullStr).Sentence[fullStr]
	findTemplate := getSentence(findStr)
	findSentence := findTemplate.Sentence[findStr]

	parts := splitSentence(fullSentence)
	for i, part := range parts { // 1
		newPart := getFirstNounCase(*part)
		if newPart != nil && newPart.Case != nil {
			parts[i] = newPart
		}
	}
	parts = filterNounless(parts)
	require.NotNil(t, parts)
	for i, part := range parts { // 2
		parts[i] = changeFirstNoun(*part, morph.CaseNomn)
	}

	rep := NewMockRepository(ctrl)
	client := NewMockMorphClient(ctrl)
	for n, part := range parts {
		if part.Sentence.Sentence() == partSentence.Sentence() {
			rep.EXPECT().
				GetByTemplate(context.Background(), partSentence).
				AnyTimes().
				Return(&findSentence, nil)
			client.EXPECT().
				ChangePOS(context.Background(), findSentence.Words[0], "NOUN")
			client.EXPECT().
				Inflect(context.Background(), findSentence.Words[0], "gent")
		} else {
			sent := sentence.Sentence{
				ID:        uint(n),
				CountWord: uint(n),
				Words:     nil,
			}
			rep.EXPECT().
				GetByTemplate(context.Background(), part.Sentence).
				AnyTimes().
				Return(&sent, errors.New(strconv.Itoa(n)))
		}
	}

	c := NewConcepterAction(rep, client)
	givenSentence, err := c.Handle(context.Background(), &fullSentence)
	require.Nil(t, err)

	expectedSentence := []sentence.Sentence{getSentence(expectedStr).
		Sentence[expectedStr]}
	require.Equal(t, true, reflect.DeepEqual(givenSentence, expectedSentence))
}

func getSentence(str string) sentence.Template {
	m := make(map[string]string)
	// partStr
	m["глагол в повелительном наклонении"] = `{
	"left": false,
	"right": false,
	"sent": {
		"глагол в повелительном наклонении": {
			"id": 0,
			"count_words": 4,
			"words": [{
				"word": "глагол",
				"normalForm": "глагол",
				"score": 1.0,
				"positionInSentence": 0,
				"tag": {
					"pos": "NOUN",
					"animacy": "inan",
					"aspect": "",
					"case": "nomn",
					"gender": "masc",
					"involvement": "",
					"mood": "",
					"number": "sing",
					"person": "",
					"tense": "",
					"transitivity": "",
					"voice": ""
				}
			}, {
				"word": "в",
				"normalForm": "в",
				"score": 0.999327,
				"positionInSentence": 0,
				"tag": {
					"pos": "PREP",
					"animacy": "",
					"aspect": "",
					"case": "",
					"gender": "",
					"involvement": "",
					"mood": "",
					"number": "",
					"person": "",
					"tense": "",
					"transitivity": "",
					"voice": ""
				}
			}, {
				"word": "повелительном",
				"normalForm": "повелительный",
				"score": 0.5,
				"positionInSentence": 0,
				"tag": {
					"pos": "ADJF",
					"animacy": "",
					"aspect": "",
					"case": "loct",
					"gender": "neut",
					"involvement": "",
					"mood": "",
					"number": "sing",
					"person": "",
					"tense": "",
					"transitivity": "",
					"voice": ""
				}
			}, {
				"word": "наклонении",
				"normalForm": "наклонение",
				"score": 1.0,
				"positionInSentence": 0,
				"tag": {
					"pos": "NOUN",
					"animacy": "inan",
					"aspect": "",
					"case": "loct",
					"gender": "neut",
					"involvement": "",
					"mood": "",
					"number": "sing",
					"person": "",
					"tense": "",
					"transitivity": "",
					"voice": ""
				}
			}]
		}
	}
}`
	// findStr
	m["перемести глагол в повелительном наклонении"] = ``
	// fullStr
	m["необходимо выполнить команду для глагола в повелительном наклонении"] = `{
	"left": false,
	"right": false,
	"sent": {
		"необходимо выполнить команду для глагола в повелительном наклонении": {
			"id": 0,
			"count_words": 8,
			"words": [{
					"word": "необходимо",
					"normalForm": "необходимо",
					"score": 0.5,
					"positionInSentence": 0,
					"tag": {
						"pos": "PRED",
						"animacy": "",
						"aspect": "",
						"case": "",
						"gender": "",
						"involvement": "",
						"mood": "",
						"number": "",
						"person": "",
						"tense": "pres",
						"transitivity": "",
						"voice": ""
					}
				}, {
					"word": "выполнить",
					"normalForm": "выполнить",
					"score": 1.0,
					"positionInSentence": 0,
					"tag": {
						"pos": "INFN",
						"animacy": "",
						"aspect": "perf",
						"case": "",
						"gender": "",
						"involvement": "",
						"mood": "",
						"number": "",
						"person": "",
						"tense": "",
						"transitivity": "tran",
						"voice": ""
					}
				},
				{
					"word": "команду",
					"normalForm": "команда",
					"score": 1.0,
					"positionInSentence": 0,
					"tag": {
						"pos": "NOUN",
						"animacy": "inan",
						"aspect": "",
						"case": "accs",
						"gender": "femn",
						"involvement": "",
						"mood": "",
						"number": "sing",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				},
				{
					"word": "для",
					"normalForm": "для",
					"score": 0.999843,
					"positionInSentence": 0,
					"tag": {
						"pos": "PREP",
						"animacy": "",
						"aspect": "",
						"case": "",
						"gender": "",
						"involvement": "",
						"mood": "",
						"number": "",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				},
				{
					"word": "глагола",
					"normalForm": "глагол",
					"score": 1.0,
					"positionInSentence": 0,
					"tag": {
						"pos": "NOUN",
						"animacy": "inan",
						"aspect": "",
						"case": "gent",
						"gender": "masc",
						"involvement": "",
						"mood": "",
						"number": "sing",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				},
				{
					"word": "в",
					"normalForm": "в",
					"score": 0.999327,
					"positionInSentence": 0,
					"tag": {
						"pos": "PREP",
						"animacy": "",
						"aspect": "",
						"case": "",
						"gender": "",
						"involvement": "",
						"mood": "",
						"number": "",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				},
				{
					"word": "повелительном",
					"normalForm": "повелительный",
					"score": 0.5,
					"positionInSentence": 0,
					"tag": {
						"pos": "ADJF",
						"animacy": "",
						"aspect": "",
						"case": "loct",
						"gender": "neut",
						"involvement": "",
						"mood": "",
						"number": "sing",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				},
				{
					"word": "наклонении",
					"normalForm": "наклонение",
					"score": 1.0,
					"positionInSentence": 0,
					"tag": {
						"pos": "NOUN",
						"animacy": "inan",
						"aspect": "",
						"case": "loct",
						"gender": "neut",
						"involvement": "",
						"mood": "",
						"number": "sing",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				}
			]
		}
	}
}`
	// expectedStr
	m["необходимо выполнить команду для перемещения"] = `{
	"left": false,
	"right": false,
	"sent": {
		"необходимо выполнить команду для перемещения": {
			"id": 0,
			"count_words": 5,
			"words": [{
					"word": "необходимо",
					"normalForm": "необходимо",
					"score": 0.5,
					"positionInSentence": 0,
					"tag": {
						"pos": "PRED",
						"animacy": "",
						"aspect": "",
						"case": "",
						"gender": "",
						"involvement": "",
						"mood": "",
						"number": "",
						"person": "",
						"tense": "pres",
						"transitivity": "",
						"voice": ""
					}
				}, {
					"word": "выполнить",
					"normalForm": "выполнить",
					"score": 1.0,
					"positionInSentence": 0,
					"tag": {
						"pos": "INFN",
						"animacy": "",
						"aspect": "perf",
						"case": "",
						"gender": "",
						"involvement": "",
						"mood": "",
						"number": "",
						"person": "",
						"tense": "",
						"transitivity": "tran",
						"voice": ""
					}
				},
				{
					"word": "команду",
					"normalForm": "команда",
					"score": 1.0,
					"positionInSentence": 0,
					"tag": {
						"pos": "NOUN",
						"animacy": "inan",
						"aspect": "",
						"case": "accs",
						"gender": "femn",
						"involvement": "",
						"mood": "",
						"number": "sing",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				},
				{
					"word": "для",
					"normalForm": "для",
					"score": 0.999843,
					"positionInSentence": 0,
					"tag": {
						"pos": "PREP",
						"animacy": "",
						"aspect": "",
						"case": "",
						"gender": "",
						"involvement": "",
						"mood": "",
						"number": "",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				},
				{
					"word": "перемещения",
					"normalForm": "перемещение",
					"score": 0.878787,
					"positionInSentence": 0,
					"tag": {
						"pos": "NOUN",
						"animacy": "inan",
						"aspect": "",
						"case": "gent",
						"gender": "neut",
						"involvement": "",
						"mood": "",
						"number": "sing",
						"person": "",
						"tense": "",
						"transitivity": "",
						"voice": ""
					}
				}
			]
		}
	}
}`

	t := sentence.Template{}
	err := json.Unmarshal([]byte(m[str]), &t)
	if err != nil {
		log.Fatalln(err)
	}
	return t
}
