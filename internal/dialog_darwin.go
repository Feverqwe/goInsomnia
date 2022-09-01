//go:build darwin

package internal

import (
	"context"
	"errors"

	"github.com/gabyx/githooks/githooks/apps/dialog/gui"
	"github.com/gabyx/githooks/githooks/apps/dialog/settings"
)

func ShowEntry(title string, text string, defaultValue string) (string, error) {
	props := settings.Entry{}
	props.DefaultEntry = defaultValue
	props.Title = title
	props.Text = text

	result, err := gui.ShowEntry(context.TODO(), &props)
	if err != nil {
		return "", err
	} else if result.IsOk() {
		return result.Text, nil
	}
	return "", errors.New("Canceled")
}
