package domain

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

const (
	LoggerNameKey = "logger_name"
)

var (
	Validator = validator.New() //nolint:gochecknoglobals

	ErrInvalidArgument  = errors.New("invalid argument")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidField     = errors.New("invalid field")

	Lang2EN      *Lang2 //nolint:gochecknoglobals
	Lang2ES      *Lang2 //nolint:gochecknoglobals
	Lang2JA      *Lang2 //nolint:gochecknoglobals
	Lang2KO      *Lang2 //nolint:gochecknoglobals
	Lang2Unknown *Lang2 //nolint:gochecknoglobals

	Lang3ENG     *Lang3 //nolint:gochecknoglobals
	Lang3ESP     *Lang3 //nolint:gochecknoglobals
	Lang3JPN     *Lang3 //nolint:gochecknoglobals
	Lang3KOR     *Lang3 //nolint:gochecknoglobals
	Lang3Unknown *Lang3 //nolint:gochecknoglobals

	Lang5ENUS    *Lang5 //nolint:gochecknoglobals
	Lang5JAJP    *Lang5 //nolint:gochecknoglobals
	Lang5Unknown *Lang5 //nolint:gochecknoglobals
)

func init() {
	initLang2()
	initLang3()
	initLang5()
}

func initLang3() {
	var err error
	Lang3ENG, err = NewLang3("eng")
	if err != nil {
		panic(err)
	}
	Lang3ESP, err = NewLang3("esp")
	if err != nil {
		panic(err)
	}
	Lang3JPN, err = NewLang3("jpn")
	if err != nil {
		panic(err)
	}
	Lang3KOR, err = NewLang3("kor")
	if err != nil {
		panic(err)
	}
	Lang3Unknown, err = NewLang3("___")
	if err != nil {
		panic(err)
	}
}

func initLang2() {
	var err error
	Lang2EN, err = NewLang2("en")
	if err != nil {
		panic(err)
	}
	Lang2ES, err = NewLang2("es")
	if err != nil {
		panic(err)
	}
	Lang2JA, err = NewLang2("ja")
	if err != nil {
		panic(err)
	}
	Lang2KO, err = NewLang2("ko")
	if err != nil {
		panic(err)
	}
	Lang2Unknown, err = NewLang2("__")
	if err != nil {
		panic(err)
	}
}

func initLang5() {
	var err error
	Lang5ENUS, err = NewLang5("en-US")
	if err != nil {
		panic(err)
	}
	Lang5JAJP, err = NewLang5("ja-JP")
	if err != nil {
		panic(err)
	}
	Lang5Unknown, err = NewLang5("_____")
	if err != nil {
		panic(err)
	}
}
