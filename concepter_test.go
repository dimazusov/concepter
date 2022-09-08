package main

import (
	"context"
	"encoding/json"
	"log"
	"optimization/internal/pkg/morph"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"optimization/internal/pkg/sentence"
)

func TestNewConcepterAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// givenSentence    перемести - глагол в повелительном наклонении

	partSentence := getSentence("глагол в повелительном наклонении")
	fullSentence := getSentence("необходимо выполнить команду для глагола в повелительном наклонении")
	expectedSentence := getSentence("необходимо выполнить команду для перемещения")
	findTemplate := getSentence("перемести глагол в повелительном наклонении")
	findSentence := findTemplate.Sentence
	replace := getReplace()

	parts := splitSentence(fullSentence.Sentence)
	parts = deepCopy(parts)
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
		if part.Sentence.Sentence() == partSentence.Sentence.Sentence() {
			rep.EXPECT().
				GetByTemplate(context.Background(), partSentence).
				AnyTimes().
				Return(&findSentence, nil)
			client.EXPECT().
				ChangePOS(context.Background(), findSentence.Words[0], "NOUN").
				Return(replace, nil)
			client.EXPECT().
				Inflect(context.Background(), replace, "gent").
				Return(expectedSentence.Sentence.Words[4], nil)
		} else {
			sent := sentence.Sentence{
				ID:        uint(n),
				CountWord: uint(n),
				Words:     nil,
			}
			template := sentence.Template{
				Sentence: part.Sentence,
				Left:     true,
				Right:    false,
			}
			rep.EXPECT().
				GetByTemplate(context.Background(), template).
				AnyTimes().
				Return(&sent, nil)
		}
	}

	c := NewConcepterAction(rep, client)
	givenSentence, err := c.Handle(context.Background(), &fullSentence.Sentence)
	require.Nil(t, err)
	require.Equal(t, true, reflect.DeepEqual(givenSentence, []sentence.Sentence{expectedSentence.Sentence}))
}

func deepCopy(parts []*sentence.Part) []*sentence.Part {
	newParts := make([]*sentence.Part, len(parts))
	for i, part := range parts {
		sent := part.Sentence
		newWords := make([]sentence.Form, len(sent.Words))
		for i, word := range sent.Words {
			tag := word.Tag
			newTag := sentence.Tag{
				POS:          check(tag.POS),
				Animacy:      check(tag.Animacy),
				Aspect:       check(tag.Aspect),
				Case:         check(tag.Case),
				Gender:       check(tag.Gender),
				Involvment:   check(tag.Involvment),
				Mood:         check(tag.Mood),
				Number:       check(tag.Number),
				Person:       check(tag.Person),
				Tense:        check(tag.Tense),
				Transitivity: check(tag.Transitivity),
				Voice:        check(tag.Voice),
			}
			newForm := sentence.Form{
				ID:                 word.ID,
				JudgmentID:         word.JudgmentID,
				Word:               word.Word,
				NormalForm:         word.NormalForm,
				Score:              word.Score,
				PositionInSentence: word.PositionInSentence,
				Tag:                newTag,
			}
			newWords[i] = newForm
		}
		newSent := sentence.Sentence{
			ID:        sent.ID,
			CountWord: sent.CountWord,
			Words:     newWords,
		}
		newPart := sentence.Part{
			Sentence: newSent,
			Case:     check(part.Case),
		}
		newParts[i] = &newPart
	}
	return newParts
}

func check(str *string) *string {
	if str != nil {
		return &(*str)
	}
	return nil
}

func getReplace() sentence.Form {
	str := `{
    "word": "перемещение",
    "normalForm": "перемещение",
    "score": 0.625,
    "positionInSentence": 0,
    "tag": {
      "pos": "NOUN",
      "animacy": "inan",
      "aspect": "",
      "case": "nomn",
      "gender": "neut",
      "involvement": "",
      "mood": "",
      "number": "sing",
      "person": "",
      "tense": "",
      "transitivity": "",
      "voice": ""
    }
  }`
	f := sentence.Form{}
	err := json.Unmarshal([]byte(str), &f)
	if err != nil {
		log.Fatalln(err)
	}
	return f
}

func getSentence(str string) sentence.Template {
	m := make(map[string]string)
	// partStr
	m["глагол в повелительном наклонении"] = `{
	"left": true,
	"right": false,
	"sentence": {
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
				"case": "nomn",
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
}`
	// findStr
	m["перемести глагол в повелительном наклонении"] = `{
	"left": true,
	"right": false,
	"sentence": {
		"id": 0,
		"count_words": 5,
		"words": [{
				"word": "перемести",
				"normalForm": "переместить",
				"score": 0.5,
				"positionInSentence": 0,
				"tag": {
					"pos": "VERB",
					"animacy": "",
					"aspect": "perf",
					"case": "",
					"gender": "",
					"involvement": "excl",
					"mood": "impr",
					"number": "sing",
					"person": "",
					"tense": "",
					"transitivity": "tran",
					"voice": ""
				}

			},
			{
				"word": "глагол",
				"normalForm": "глагол",
				"score": 0.75,
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
}`
	// fullStr
	m["необходимо выполнить команду для глагола в повелительном наклонении"] = `{
	"left": false,
	"right": false,
	"sentence": {
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
}`
	// expectedStr
	m["необходимо выполнить команду для перемещения"] = `{
	"left": false,
	"right": false,
	"sentence": {
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
}`

	t := sentence.Template{}
	err := json.Unmarshal([]byte(m[str]), &t)
	if err != nil {
		log.Fatalln(err)
	}
	return t
}
