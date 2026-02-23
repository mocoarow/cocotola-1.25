package service

import (
	"errors"
)

var ErrPairOfUserAndGroupAlreadyExists = errors.New("pair of user and group already exists")
var ErrPairOfUserAndGroupNotFound = errors.New("pair of user and group not found")
