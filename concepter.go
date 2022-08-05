package main

import (
	"context"

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

type concepter struct {
	rep Repository
}

func NewConcepterAction(rep Repository) *concepter {
	return &concepter{rep: rep}
}

func (m concepter) Handle(ctx context.Context, s *sentence.Sentence) (judgments []sentence.Sentence, err error) {
	// TODO
	//необходимо выполнить команду для глагола в повелительном наклонении
	//разбиваем на части // по 1, по 2 ... по n
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
