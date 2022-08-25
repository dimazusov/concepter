package main

import (
	"context"
	"github.com/pkg/errors"
	"optimization/internal/pkg/morph"
	"optimization/internal/pkg/sentence"
)

type Repository interface {
	GetByTemplate(ctx context.Context, j sentence.Sentence) ([]sentence.Sentence, error)
}

type MorphClient interface {
	// Склонение
	Inflect(ctx context.Context, word sentence.Form, wordCase []string)
	// Изменения части речи // pos - part of speach
	ChangePOS(ctx context.Context, word sentence.Form, pos string)
}

type template struct {
	part     *sentence.Part
	sentence []sentence.Sentence
}

type concepter struct {
	rep Repository
}

func NewConcepterAction(rep Repository) *concepter {
	return &concepter{rep}
}

func (m concepter) Handle(ctx context.Context, s *sentence.Sentence) (judgments []sentence.Sentence, err error) {
	parts := s.SplitSentence()
	findCases(parts) // 1
	parts = removeUnnecessary(parts)
	if parts == nil {
		return nil, errors.New("the sentence does not contain any nouns")
	}
	for _, part := range parts { // 2
		for n, word := range part.Sentence.Words {
			if word.Word == part.Word.Word { // возможно неверно
				part.Sentence.Words[n].ToNomn()
			}
		}
	}
	template, err := m.findTemplate(ctx, parts)
	if err != nil {
		return nil, err
	}
	_ = template

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

func (m concepter) findTemplate(ctx context.Context, parts []*sentence.Part) (*template, error) {
	defer func() {
		if r := recover(); r != nil {
			parts = parts[1:]
		}
	}()
	if len(parts) == 1 {
		return &template{
			part:     parts[0],
			sentence: []sentence.Sentence{parts[0].Sentence},
		}, nil
	}
	s, err := m.rep.GetByTemplate(ctx, parts[0].Sentence)
	return &template{
		part:     parts[0],
		sentence: s,
	}, err
}

func findCases(parts []*sentence.Part) {
	for _, part := range parts {
		for _, word := range part.Sentence.Words {
			if *word.Tag.POS == morph.PartOfSpeachNOUN {
				(*part).Word = &word
				break
			}
		}
	}
}

func removeUnnecessary(parts []*sentence.Part) []*sentence.Part {
	var result []*sentence.Part
	for _, part := range parts {
		if part.Word != nil {
			result = append(result, part)
		}
	}
	return result
}
