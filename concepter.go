package main

import (
	"context"
	"github.com/pkg/errors"
	"optimization/internal/pkg/morph"
	"optimization/internal/pkg/sentence"
)

type Repository interface {
	GetByTemplate(ctx context.Context, j sentence.Sentence) (*sentence.Sentence, error)
}

type MorphClient interface {
	// Склонение
	Inflect(ctx context.Context, word sentence.Form, wordCase string) (sentence.Form, error)
	// Изменения части речи // pos - part of speach
	ChangePOS(ctx context.Context, word sentence.Form, pos string) (sentence.Form, error)
}

type concepter struct {
	rep    Repository
	client MorphClient
}

func NewConcepterAction(rep Repository, client MorphClient) *concepter {
	return &concepter{rep, client}
}

func (m concepter) Handle(ctx context.Context, s *sentence.Sentence) (judgments []sentence.Sentence, err error) {
	parts := SplitSentence(*s)
	findCases(parts) // 1
	parts = removeWithoutNouns(parts)
	if parts == nil {
		return nil, errors.New("the sentence does not contain any nouns")
	}
	NounsToNomn(parts)                            // 2
	sent, part, err := m.findTemplate(ctx, parts) // 3
	if err != nil {
		return nil, err
	}
	sent.Words[0], err = m.client.ChangePOS(ctx, sent.Words[0], "NOUN") // 4
	if err != nil {
		return nil, err
	}
	m.client.Inflect(ctx, sent.Words[0], *part.Case)
	_ = part

	// TODO
	//необходимо выполнить команду для глагола в повелительном наклонении
	//разбиваем на части // по 1, по 2 ... по n // n - количество слов
	//затем для каждой части
	//
	//// 1. ищем первое с начала имя сущ., получаем падеж, для того чтобы потом склонить в эту форму слово
	//глагола в повелительном наклонении (глагола - родительный)
	//
	//// 2. (ищем первое существительное и переводим его в им. падеж)
	//глагола в повелительном наклонении -> глагол в повелительном наклонении
	//
	//// 3. поиск по шаблону используя начальную форму: {любое словосочетание} - глагол в повелительном наклонении
	//перемести - глагол в повелительном наклонении
	//
	//// 4. берем normalForm найденного словосочетания и переводим его в имя сущ
	//c помощью 	ChangePOS(ctx context.Context, word sentence.Form, pos string)
	//перемести -> переместить
	//переместить -> перемещение
	//
	//// 5. склоняем первое имя сущ. в падеж полученный из 1 (в данном случае родительный падеж)
	//Inflect(word sentence.Form, wordCase []string)
	//перемещение -> перемещения
	//
	//// 6. проводим замену
	//необходимо выполнить команду для глагола в повелительном наклонении
	//->
	//необходимо выполнить команду для перемещения

	return nil, nil
}

func SplitSentence(s sentence.Sentence) []*sentence.Part {
	var parts []*sentence.Part
	for i := 0; uint(i) < s.CountWord; i++ {
		for j := i + 1; uint(j) <= s.CountWord; j++ {
			var words []sentence.Form
			for _, word := range s.Words[i:j] {
				w := word
				w.Tag = word.Tag
				words = append(words, w)
			}
			sent := sentence.Sentence{
				ID:        s.ID,
				CountWord: uint(len(words)),
				Words:     words,
			}
			part := sentence.Part{Sentence: sent}
			parts = append(parts, &part)
		}
	}
	return parts
}

func (m concepter) findTemplate(ctx context.Context, parts []*sentence.Part) (*sentence.Sentence, *sentence.Part, error) {
	for _, part := range parts {
		s, err := m.rep.GetByTemplate(ctx, part.Sentence)
		if s != nil {
			return s, part, err
		}
	}
	return nil, nil, errors.New("template not found")
}

func NounsToNomn(parts []*sentence.Part) {
	for _, part := range parts {
		NounToNomn(part)
	}
}

func NounToNomn(part *sentence.Part) {
	for n, word := range part.Sentence.Words {
		if word.Tag.Case == part.Case {
			form := &part.Sentence.Words[n]
			changeCase(form, morph.CaseNomn)
		}
	}
}

func changeCase(form *sentence.Form, wordCase string) {
	form.Word = form.NormalForm
	*form.Tag.Case = wordCase
}

func findCases(parts []*sentence.Part) {
	for _, part := range parts {
		findCase(part)
	}
}

func findCase(part *sentence.Part) {
	for _, word := range part.Sentence.Words {
		if isNoun(word) {
			(*part).Case = word.Tag.Case
			break
		}
	}
}

func isNoun(word sentence.Form) bool {
	return *word.Tag.POS == morph.PartOfSpeachNOUN
}

func removeWithoutNouns(parts []*sentence.Part) []*sentence.Part {
	var result []*sentence.Part
	for _, part := range parts {
		if part.Case != nil {
			result = append(result, part)
		}
	}
	return result
}
