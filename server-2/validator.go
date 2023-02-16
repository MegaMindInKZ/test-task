package main

import (
	"regexp"
)

type Validator struct {
	Errors []string
}

func (v *Validator) isValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Check(email string) {
	if _, err := GetUserByEmail(email); err == nil {
		v.Errors = append(v.Errors, ErrDuplicationEmail.Error())
	}
	if r := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`); !r.MatchString(email) {
		v.Errors = append(v.Errors, ErrNotAcceptableEmailFormat.Error())
	}
}

func NewValidator() Validator {
	return Validator{}
}
