//go:build windows || linux

package internal

import (
	"errors"

	"github.com/ncruces/zenity"
)

func ShowEntry(title string, text string, defaultValue string) (string, error) {
	text, err := zenity.Entry(text,
		zenity.Title(title),
		zenity.EntryText(defaultValue),
	)
	if err != nil {
		if err.Error() == "dialog canceled" {
			return "", errors.New("Canceled")
		}
		return "", err
	} else {
		return text, nil
	}
}
