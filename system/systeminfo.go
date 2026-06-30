// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_systeminfo

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func read(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

var systemsDetailLabel *widget.Label //ประกาศแบบ golbal
func SystemsDetailLabelcmd(text string) {
	if systemsDetailLabel != nil {
		systemsDetailLabel.SetText(text)
	}
}

func SystemTabs() fyne.CanvasObject {

	systemsDetailLabel = widget.NewLabel("")
	/*
		subBIOS_UEFI_label := container.NewVBox(
			//BIOS/UEFI
			widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", b["Bios_vendor"])),   //ผู้ผลิต BIOS
			widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", b["Bios_version"])), //เวอร์ชัน BIOS
			widget.NewLabel(fmt.Sprintf("วันที่ออก : %s", b["Bios_date"])),   //วันที่ออก BIOS
			widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", b["Bios_release"])), //เวอร์ชัน Release ของ BIOS ตาม SMBIOS
		)
	*/

	/*
		subBIOS_UEFI := container.NewVBox(
			//BIOS/UEFI
			widget.NewCard("BIOS_UEFI", "", subBIOS_UEFI_label),
		)
	*/

	SysTemsDetail := container.NewVBox(
		//detail
		widget.NewCard("Systems", "", systemsDetailLabel),
	)

	return container.NewScroll(SysTemsDetail)

	/*
	   return container.NewAppTabs(
	   	container.NewTabItem("System", container.NewScroll(SysTemsDetail)),
	   	//container.NewTabItem("Detail", container.NewScroll(sub_Detail_BIOS_UEFI)),
	   )
	*/
}
