package main

import (
	"context"
	"github.com/pkg/errors"
	"optimization/internal/pkg/morph"
	"optimization/internal/pkg/sentence"
)

type Repository interface {
	GetByTemplate(ctx context.Context, j sentence.Template) (*sentence.Sentence, error)
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
	parts := splitSentence(*s)
	for i, part := range parts { // 1
		newPart := getFirstNounCase(*part)
		if newPart != nil && newPart.Case != nil {
			parts[i] = newPart
		}
	}
	parts = filterNounless(parts)
	if parts == nil {
		return nil, errors.New("the sentence does not contain any nouns")
	}
	for i, part := range parts { // 2
		parts[i] = changeFirstNoun(*part, morph.CaseNomn)
	}
	sent, part, err := m.findTemplate(ctx, parts) // 3
	if err != nil {
		return nil, err
	}
	// TODO - переделать с одного слова на словосочетание
	sent.Words[0], err = m.client.ChangePOS(ctx, sent.Words[0], "NOUN") // 4
	if err != nil {
		return nil, err
	}
	replacement, err := m.client.Inflect(ctx, sent.Words[0], *part.Case) // 5
	if err != nil {
		return nil, err
	}
	_ = replacement
	s.Words = append(s.Words[:part.Indexes.I], replacement) // 6
	s.CountWord = uint(len(s.Words))

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
	//перемести -> переместить // берем normalForm
	//переместить -> перемещение // используем changePOS
	//
	//// 5. склоняем первое имя сущ. в падеж полученный из 1 (в данном случае родительный падеж)
	//Inflect(word sentence.Form, wordCase []string)
	//перемещение -> перемещения
	//
	//// 6. проводим замену
	//необходимо выполнить команду для глагола в повелительном наклонении
	//->
	//необходимо выполнить команду для перемещения

	return []sentence.Sentence{*s}, nil
}

func splitSentence(s sentence.Sentence) []*sentence.Part {
	var parts []*sentence.Part
	for i := 0; uint(i) <= s.CountWord; i++ {
		for j := i + 1; uint(j) <= s.CountWord; j++ {
			words := make([]sentence.Form, j-i)
			for idx, word := range s.Words[i:j] { // фабричный метод?
				w := word
				w.Tag = word.Tag
				words[idx] = w
			}
			sent := sentence.Sentence{
				ID:        s.ID,
				CountWord: uint(len(words)),
				Words:     words,
			}
			part := sentence.Part{Sentence: sent, Indexes: sentence.Indexes{I: i, J: j}}
			parts = append(parts, &part)
		}
	}
	return parts
}

func (m concepter) findTemplate(ctx context.Context, parts []*sentence.Part) (*sentence.Sentence, *sentence.Part, error) {
	for _, part := range parts {
		template := sentence.Template{
			Sentence: part.Sentence,
			Left:     true,
			Right:    false,
		}
		s, err := m.rep.GetByTemplate(ctx, template)
		if s != nil && s.Words != nil {
			return s, part, err
		}
	}
	return nil, nil, errors.New("template not found")
}

func changeFirstNoun(part sentence.Part, wordCase string) *sentence.Part {
	newSentence := part.Sentence
	newCase := part.Case
	if newCase != nil {
		s := *newCase
		newCase = &s
	}
	newPart := sentence.Part{
		Sentence: newSentence,
		Case:     newCase,
		Indexes:  part.Indexes,
	}
	for n, word := range part.Sentence.Words {
		if word.Tag.Case == *part.Case {
			form := newPart.Sentence.Words[n]
			form = changeCase(form, wordCase)
			newPart.Sentence.Words[n] = form
			return &newPart
		}
	}
	return &part
}

func changeCase(form sentence.Form, wordCase string) sentence.Form {
	form.Word = form.NormalForm
	form.Tag.Case = wordCase
	return form
}

func getFirstNounCase(part sentence.Part) *sentence.Part {
	for _, word := range part.Sentence.Words {
		if word.Tag.POS == morph.PartOfSpeachNOUN {
			part.Case = &word.Tag.Case // должно меняться все слово, а не только падеж
			return &part
		}
	}
	return nil
}

func filterNounless(parts []*sentence.Part) []*sentence.Part {
	var result []*sentence.Part
	for _, part := range parts {
		if part.Case != nil {
			result = append(result, part)
		}
	}
	return result
}
