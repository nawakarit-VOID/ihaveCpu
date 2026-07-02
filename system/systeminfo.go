// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package system

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var systemsDetailLabel *widget.Label //ประกาศแบบ golbal
func SystemsDetailLabelcmd(text string) {
	if systemsDetailLabel != nil {
		systemsDetailLabel.SetText(text)
	}
}

func SystemTabs() fyne.CanvasObject {

	systemsDetailLabel = widget.NewLabel("")

	SysTemsDetail := container.NewVBox(
		//detail
		widget.NewCard("Systems", "", systemsDetailLabel),
	)

	return container.NewScroll(SysTemsDetail)

}
