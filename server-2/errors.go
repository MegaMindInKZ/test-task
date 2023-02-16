package main

import "errors"

var ErrNoDocument = errors.New("mongo: no documents in result")
var ErrDuplicationEmail = errors.New("duplication email")
var ErrNotAcceptableEmailFormat = errors.New("not acceptable email format")
