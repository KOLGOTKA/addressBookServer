package pkg

import "github.com/pkg/errors"

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
	for i := 0; i < len(phone); i++ {
		l := phone[i]
		if l >= '0' && l <= '9' {
			normalizedPhone += string(l)
		}
	}
	if string(normalizedPhone[0]) == "8" {
		normalizedPhone = "7" + normalizedPhone[1:]
	}

	if string(normalizedPhone[0]) != "7" {
		err = errors.New("Incorrect phone number: " + normalizedPhone)
		return "", err
	}
	if len(normalizedPhone) != 11 {
		err = errors.New("Incorrect len of phone number: " + string(len(normalizedPhone)))
		return normalizedPhone, err
	}
	return "", err
}
