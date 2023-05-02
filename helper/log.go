package helper

import "fyne.io/fyne/v2/widget"

var LogLabel *widget.Label

func Log(msg string) {
	LogLabel.SetText("2006-01-02 15:04:05" + " - " + msg + LogLabel.Text)
}
