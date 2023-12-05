package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	ValidErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.ValidErrors) == 0
}

func (v *Validator) AddfieldErros(key, message string) {
	if v.ValidErrors == nil {
		v.ValidErrors = make(map[string]string)
	}
	if _, exist := v.ValidErrors[key]; !exist {
		v.ValidErrors[key] = message
	}
}

func (v *Validator) Checkvield(ok bool, key, message string) {
	if !ok {
		v.AddfieldErros(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, PermittedValue ...int) bool {
	for i := range PermittedValue {
		if value == PermittedValue[i] {
			return true
		}
	}
	return false
}
