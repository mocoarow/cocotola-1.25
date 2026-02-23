package service

import "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

var ListDecksAction = domain.NewRBACAction("ListDecks") //nolint:gochecknoglobals
var ReadDeckAction = domain.NewRBACAction("ReadDeck")   //nolint:gochecknoglobals
