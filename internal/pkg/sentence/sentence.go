package sentence

import (
	"optimization/internal/pkg/morph"
)

type Sentence struct {
	ID        uint   `json:"id" db:"id"`
	CountWord uint   `json:"count_words"`
	Words     []Form `json:"words" gorm:"foreignKey:JudgmentID"`
}

type Form struct {
	ID                 uint    `json:"id" db:"id"`
	JudgmentID         uint    `json:"judgmentID" db:"judgment_id"`
	Word               string  `json:"word" db:"word"`
	NormalForm         string  `json:"normalWord" db:"normal_word"`
	Score              float64 `json:"score" db:"score"`
	PositionInSentence int     `json:"positionInSentence" db:"position_in_sentence"`
	Tag                Tag     `json:"tag" db:"tag" gorm:"embedded;embeddedPrefix:tag_"`
}

type Tag struct {
	POS          *string `json:"pos" db:"pos"`
	Animacy      *string `json:"animacy" db:"animacy"`
	Aspect       *string `json:"aspect" db:"aspect"`
	Case         *string `json:"case" db:"case"`
	Gender       *string `json:"gender" db:"gender"`
	Involvment   *string `json:"involvment" db:"involvment"`
	Mood         *string `json:"mood" db:"mood"`
	Number       *string `json:"number" db:"number"`
	Person       *string `json:"person" db:"person"`
	Tense        *string `json:"tense" db:"tense"`
	Transitivity *string `json:"transitivity" db:"transitivity"`
	Voice        *string `json:"voice" db:"voice"`
}

type Part struct {
	Sentence Sentence
	Word     *Form
}

func (s *Sentence) SplitSentence() []*Part {
	var (
		parts []*Part
		id    = 0
	)
	for i := 0; uint(i) < s.CountWord; i++ {
		for j := i + 1; uint(j) <= s.CountWord; j++ {
			parts = append(parts, &Part{Sentence{
				ID:        uint(id),
				CountWord: uint(id),
				Words:     s.Words[i:j],
			}, nil})
			id++
		}
	}
	return parts
}

func (f *Form) ToNomn() { // скорее всего это не все
	f.Word = f.NormalForm
	*f.Tag.POS = morph.CaseNomn
}
