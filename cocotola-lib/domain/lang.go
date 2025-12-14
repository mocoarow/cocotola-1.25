package domain

import "fmt"

const Lang2Len = 2
const Lang3Len = 3
const Lang5Len = 5

type Lang2 struct {
	value string
}

func NewLang2(lang string) (*Lang2, error) {
	if len(lang) != Lang2Len {
		return nil, fmt.Errorf("invalid parameter. Lang2: %s. err: %w", lang, ErrInvalidArgument)
	}

	return &Lang2{
		value: lang,
	}, nil
}

func (l *Lang2) String() string {
	return l.value
}

func (l *Lang2) ToLang3() *Lang3 {
	switch l.value {
	case "en":
		return Lang3ENG
	case "es":
		return Lang3ESP
	case "ja":
		return Lang3JPN
	case "ko":
		return Lang3KOR
	default:
		return Lang3Unknown
	}
}

type Lang3 struct {
	value string
}

func NewLang3(lang string) (*Lang3, error) {
	if len(lang) != Lang3Len {
		return nil, fmt.Errorf("invalid parameter. Lang3: %s, err: %w", lang, ErrInvalidArgument)
	}

	return &Lang3{
		value: lang,
	}, nil
}

func (l *Lang3) String() string {
	return l.value
}

func (l *Lang3) ToLang2() *Lang2 {
	switch l.value {
	case "eng":
		return Lang2EN
	case "esp":
		return Lang2ES
	case "jpn":
		return Lang2JA
	case "kor":
		return Lang2KO
	default:
		return Lang2Unknown
	}
}

type Lang5 struct {
	value string
}

func NewLang5(lang string) (*Lang5, error) {
	if len(lang) != Lang5Len {
		return nil, fmt.Errorf("invalid parameter. Lang5: %s", lang)
	}

	return &Lang5{
		value: lang,
	}, nil
}

func (l *Lang5) String() string {
	return l.value
}

func (l *Lang5) ToLang2() *Lang2 {
	switch l.value {
	case "en-US":
		return Lang2EN
	case "es-ES":
		return Lang2ES
	case "ja-JP":
		return Lang2JA
	case "ko-KR":
		return Lang2KO
	default:
		return Lang2Unknown
	}
}
