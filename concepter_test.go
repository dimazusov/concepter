package main

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
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
	findSentence := getSentence("перемести глагол в повелительном наклонении")
	fullSentence := getSentence("необходимо выполнить команду для глагола в повелительном наклонении")

	parts := SplitSentence(fullSentence)
	findCases(parts) // 1
	parts = removeWithoutNouns(parts)
	require.NotNil(t, parts)
	NounsToNomn(parts) // 2

	rep := NewMockRepository(ctrl)
	client := NewMockMorphClient(ctrl)
	for _, part := range parts {
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
			rep.EXPECT().
				GetByTemplate(context.Background(), part.Sentence).
				AnyTimes().
				Return(&part.Sentence, errors.New(""))
		}
	}

	c := NewConcepterAction(rep, client)
	givenSentence, err := c.Handle(context.Background(), &fullSentence)
	require.Nil(t, err)

	expectedSentence := []sentence.Sentence{getSentence("необходимо выполнить команду для перемещения")}
	require.Equal(t, true, reflect.DeepEqual(givenSentence, expectedSentence))
}

//func TestNewConcepterAction(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// givenSentence    перемести - глагол в повелительном наклонении
//
//	partSentence := getSentence("глагол в повелительном наклонении")
//	findSentence := getSentence("перемести глагол в повелительном наклонении")
//	fullSentence := getSentence("необходимо выполнить команду для глагола в повелительном наклонении")
//
//	parts := SplitSentence(fullSentence)
//	findCases(parts) // 1
//	parts = removeWithoutNouns(parts)
//	require.NotNil(t, parts)
//	NounsToNomn(parts) // 2
//
//	rep := NewMockRepository(ctrl)
//	client := NewMockMorphClient(ctrl)
//
//	rep.EXPECT().
//		GetByTemplate(context.Background(), partSentence).
//		AnyTimes().
//		Return(&findSentence, nil)
//	client.EXPECT().
//		ChangePOS(context.Background(), findSentence.Words[0], "NOUN")
//	client.EXPECT().
//		Inflect(context.Background(), findSentence.Words[0], "gent")
//
//	c := NewConcepterAction(rep, client)
//	sent, _ := rep.GetByTemplate(context.Background(), partSentence)
//	Case := "gent"
//	part := sentence.Part{
//		Sentence: partSentence,
//		Case:     &Case,
//	}
//	c.client.ChangePOS(context.Background(), sent.Words[0], "NOUN") // 4
//	c.client.Inflect(context.Background(), sent.Words[0], *part.Case)
//	_ = part
//}

func getSentence(str string) sentence.Sentence {
	m := make(map[string]string)
	m["глагол в повелительном наклонении"] = `{
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
}`
	m["перемести глагол в повелительном наклонении"] = `{
	"id": 0,
	"count_words": 4,
	"words": [{
		"id": 0,
		"count_words": 1,
		"words": {
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
		}
	}, {
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
}`
	m["необходимо выполнить команду для глагола в повелительном наклонении"] = `{
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
}`
	m["необходимо выполнить команду для перемещения"] = `{
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
}`

	t := sentence.Sentence{}
	err := json.Unmarshal([]byte(m[str]), &t)
	if err != nil {
		log.Fatalln(err)
	}
	return t
}
