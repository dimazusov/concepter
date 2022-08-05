package morph

type Form struct {
	Word               string  `json:"word"`
	NormalForm         string  `json:"normalWord"`
	Score              float64 `json:"score"`
	PositionInSentence int     `json:"positionInSentence"`
	Tag                Tag     `json:"tag"`
}
