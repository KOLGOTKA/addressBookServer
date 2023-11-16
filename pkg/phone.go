package pkg

import (
	"regexp"
	"strconv"
	"strings"
	// "github.com/dongri/phonenumber"
	// "github.com/pkg/errors"
)

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
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
		err = ErrorGenerate("Empty phone number")
		return "", err
	}
	if normalizedPhone[0] == '8' {
		normalizedPhone = "7" + normalizedPhone[1:]
	}

	if normalizedPhone[0] != '7' {
		err = ErrorGenerate("Incorrect phone number: " + normalizedPhone)
		return "", err
	}
	if len(normalizedPhone) != 11 {
		err = ErrorGenerate("Incorrect len of phone number: " + strconv.Itoa((len(normalizedPhone))))
		return "", err
	}
	return normalizedPhone, err
}

/// Ещё один вариант реализации данной функции

// / С помощью регулярных выражений
func PhoneNormalize3(phone string) (normalizedPhone string, err error) {
	normalizedPhone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	/// Проверка на пустой номер телефона
	if normalizedPhone == "" {
		err = ErrorGenerate("Empty phone number")
		return "", err
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
		err = ErrorGenerate("Incorrect phone number: " + phone)
		return "", err
	}
	return normalizedPhone, err
}
