package pkg

import (
	// "errors"
	"regexp"
	"strconv"
	"strings"
	// "github.com/dongri/phonenumber"
	// "github.com/pkg/errors"
)

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
	myerr := NewMyError("psg: PhoneNormalize(phone string)")
	// errorer = &Errorer{Where: "psg: PhoneNormalize(phone string)"}
	var builder strings.Builder
	for i := range phone {
		l := phone[i]
		if l >= '0' && l <= '9' {
			builder.WriteByte(l)
		}
	}
	normalizedPhone = builder.String()
	/// Проверка на пустой номер телефона
	if builder.Len() == 0 {
		return "", myerr.Wrap(nil, "Empty phone number")
		// errorer.Add("Empty phone number")
		// return "", errorer
	}
	if normalizedPhone[0] == '8' {
		normalizedPhone = "7" + normalizedPhone[1:]
	}

	if normalizedPhone[0] != '7' {
		return "", myerr.Wrap(nil, "Incorrect phone number: " + normalizedPhone)
		// errorer.Add("Incorrect phone number: " + normalizedPhone)
		// return "", errorer
	}
	if len(normalizedPhone) != 11 {
		return "", myerr.Wrap(nil, "Incorrect len of phone number: " + strconv.Itoa(len(normalizedPhone)))
		// errorer.Add("Incorrect len of phone number: " + strconv.Itoa(len(normalizedPhone)))
		// return "", errorer
	}
	return normalizedPhone, nil
}

/// Ещё один вариант реализации данной функции

// / С помощью регулярных выражений
func PhoneNormalize3(phone string) (normalizedPhone string, err error) {
	myerr := NewMyError("psg: PhoneNormalize3(phone string)")
	// errorer = &Errorer{Where: "psg: PhoneNormalize3(phone string)"}
	normalizedPhone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	/// Проверка на пустой номер телефона
	if normalizedPhone == "" {
		return "", myerr.Wrap(nil, "Empty phone number")
		// errorer.Add("Empty phone number")
		// return "", errorer
	}
	/// Проверка на то, что телефон состоит из 10 цифр
	if matched, _ := regexp.MatchString("^\\d{10}$", normalizedPhone); matched {
		normalizedPhone = "7" + normalizedPhone
	}
	/// Проверка на то, что первыя цифра 8
	if string(normalizedPhone[0]) == "8" {
		normalizedPhone = "7" + normalizedPhone[1:]
	}
	// Проверяем, корректность номера телефона (состоит только из цифр, имеет правильную длину и начинается с 7)
	if matched, _ := regexp.MatchString("^7\\d{10}$", normalizedPhone); !matched {
		return "", myerr.Wrap(nil, "Incorrect phone number: " + phone)
		// errorer.Add("Incorrect phone number: " + phone)
		// return "", errorer
	}
	return normalizedPhone, nil
}
