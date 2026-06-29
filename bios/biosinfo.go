// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_biosinfo

import (
	"fmt"
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

func bios_info() map[string]interface{} {

	//BIOS/UEFI
	bios_vendor := read("/sys/class/dmi/id/bios_vendor")
	bios_version := read("/sys/class/dmi/id/bios_version")
	bios_date := read("/sys/class/dmi/id/bios_date")
	bios_release := read("/sys/class/dmi/id/bios_release")

	return map[string]interface{}{
		//BIOS/UEFI
		"Bios_vendor":  bios_vendor,  //ผู้ผลิต BIOS
		"Bios_version": bios_version, //เวอร์ชัน BIOS
		"Bios_date":    bios_date,    //วันที่ออก BIOS
		"Bios_release": bios_release, //เวอร์ชัน Release ของ BIOS ตาม SMBIOS
	}
}

var biosDetailLabel *widget.Label //ประกาศแบบ golbal
func BiosDetailLabelcmd(text string) {
	if biosDetailLabel != nil {
		biosDetailLabel.SetText(text)
	}
}

func BiosTabs() fyne.CanvasObject {
	b := bios_info()

	biosDetailLabel = widget.NewLabel("")

	subBIOS_UEFI_label := container.NewVBox(
		//BIOS/UEFI
		widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", b["Bios_vendor"])),   //ผู้ผลิต BIOS
		widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", b["Bios_version"])), //เวอร์ชัน BIOS
		widget.NewLabel(fmt.Sprintf("วันที่ออก : %s", b["Bios_date"])),   //วันที่ออก BIOS
		widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", b["Bios_release"])), //เวอร์ชัน Release ของ BIOS ตาม SMBIOS
	)

	subBIOS_UEFI := container.NewVBox(
		//BIOS/UEFI
		widget.NewCard("BIOS_UEFI", "", subBIOS_UEFI_label),
	)

	sub_Detail_BIOS_UEFI := container.NewVBox(
		//detail
		widget.NewCard("Detail", "", biosDetailLabel),
	)

	return container.NewAppTabs(
		container.NewTabItem("BIOS/UEFI", container.NewScroll(subBIOS_UEFI)),
		container.NewTabItem("Detail", container.NewScroll(sub_Detail_BIOS_UEFI)),
	)
}
