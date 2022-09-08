package sentence

type Sentence struct {
	ID        uint   `json:"id" db:"id"`
	CountWord uint   `json:"count_words"`
	Words     []Form `json:"words" gorm:"foreignKey:JudgmentID"`
}

func (s Sentence) Sentence() string {
	var result string
	for _, word := range s.Words {
		result += word.Word + " "
	}
	return result
}

type Form struct {
	ID                 uint    `json:"id" db:"id"`
	JudgmentID         uint    `json:"judgmentID" db:"judgment_id"`
	Word               string  `json:"word" db:"word"`
	NormalForm         string  `json:"normalForm" db:"normal_form"`
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
	Case     *string
}

type Template struct {
	Sentence Sentence `json:"sentence" db:"sentence"`
	Left     bool     `json:"left" db:"left"`
	Right    bool     `json:"right" db:"right"`
}
