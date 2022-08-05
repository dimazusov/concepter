package morph

// Часть Речи POS
const PartOfSpeachNOUN = "NOUN" //	имя существительное	хомяк
const PartOfSpeachADJF = "ADJF" //	имя прилагательное (полное)	хороший
const PartOfSpeachADJS = "ADJS" //	имя прилагательное (краткое)	хорош
const PartOfSpeachCOMP = "COMP" //	компаратив	лучше, получше, выше
const PartOfSpeachVERB = "VERB" //	глагол (личная форма)	говорю, говорит, говорил
const PartOfSpeachINFN = "INFN" //	глагол (инфинитив)	говорить, сказать
const PartOfSpeachPRTF = "PRTF" //	причастие (полное)	прочитавший, прочитанная
const PartOfSpeachPRTS = "PRTS" //	причастие (краткое)	прочитана
const PartOfSpeachGRND = "GRND" //	деепричастие	прочитав, рассказывая
const PartOfSpeachNUMR = "NUMR" //	числительное	три, пятьдесят
const PartOfSpeachADVB = "ADVB" //	наречие	круто
const PartOfSpeachNPRO = "NPRO" //	местоимение-существительное	он
const PartOfSpeachPRED = "PRED" //	предикатив	некогда
const PartOfSpeachPREP = "PREP" //	предлог	в
const PartOfSpeachCONJ = "CONJ" //	союз	и
const PartOfSpeachPRCL = "PRCL" //	частица	бы, же, лишь
const PartOfSpeachINTJ = "INTJ" //	междометие	ой

// Падеж
const CaseNomn = "nomn" //	именительный	Кто? Что?	хомяк ест
const CaseGent = "gent" //	родительный	Кого? Чего?	у нас нет хомяка
const CaseDatv = "datv" //	дательный	Кому? Чему?	сказать хомяку спасибо
const CaseAccs = "accs" //	винительный	Кого? Что?	хомяк читает книгу
const CaseAblt = "ablt" //	творительный	Кем? Чем?	зерно съедено хомяком
const CaseLoct = "loct" //	предложный	О ком? О чём? и т.п.	хомяка несут в корзинке
const CaseVoct = "voct" //	звательный	Его формы используются при обращении к человеку.	Саш, пойдем в кино.
const CaseGen2 = "gen2" //	второй родительный (частичный)	 	ложка сахару (gent - производство сахара); стакан яду (gent - нет яда)
const CaseAcc2 = "acc2" //	второй винительный	 	записался в солдаты
const CaseLoc2 = "loc2" //	второй предложный (местный)	 	я у него в долгу (loct - напоминать о долге); висит в шкафу (loct - монолог о шкафе); весь в снегу (loct - писать о снеге)

// Число
const NumberSing = "sing"
const NumberPlug = "sing"

// Род
const GenderMasc = "masc" // мужской род	хомяк, говорил
const GenderFemn = "femn" // женский род	хомячиха, говорила
const GenderNeut = "neut" // средний род	зерно, говорило

// LATN	Токен состоит из латинских букв (например, “foo-bar” или “Maßstab”)
// PNCT	Пунктуация (например, , или !? или …)
// NUMB	Число (например, “204” или “3.14”)
// intg	целое число (например, “204”)
// real	вещественное число (например, “3.14”)
// ROMN	Римское число (например, XI)
// UNKN	Токен не удалось разобрать

type Tag struct {
	POS          string `json:"pos"`
	Animacy      string `json:"animacy"`
	Aspect       string `json:"aspect"`
	Case         string `json:"case"`
	Gender       string `json:"gender"`
	Involvment   string `json:"involvment"`
	Mood         string `json:"mood"`
	Number       string `json:"number"`
	Person       string `json:"person"`
	Tense        string `json:"tense"`
	Transitivity string `json:"transitivity"`
	Voice        string `json:"voice"`
}
