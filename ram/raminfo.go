// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_raminfo

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jaypipes/ghw"
)

func Memory() *ghw.MemoryInfo {
	//ข้อมูล RAM
	memInfo, err := ghw.Memory()
	if err != nil {
		return nil
	}

	return memInfo
}

func RamTabs() fyne.CanvasObject {

	memInfo := Memory()

	memName := fmt.Sprintf("%s", &memInfo.Area)
	memDefaulSize := fmt.Sprintf("%d", &memInfo.DefaultHugePageSize) //เปลี่ยนแปลงตลอด
	memModule := fmt.Sprintf("%T", &memInfo.Modules)                 //[]
	memHugeSize := fmt.Sprintf("%T", &memInfo.HugePageAmountsBySize) //map uint64
	memSupportSize := fmt.Sprintf("%T", &memInfo.SupportedPageSizes) //[] uint64
	memTotalHugeBytes := fmt.Sprintf("%d", &memInfo.TotalHugePageBytes)
	memTotalPhysicalBytes := fmt.Sprintf("%d", &memInfo.TotalPhysicalBytes)
	memTotalUsableBytes := fmt.Sprintf("%d", &memInfo.TotalUsableBytes)

	//	var xx string
	for _, area := range memInfo.Modules {
		fmt.Printf("%+v\n", area)
	}

	b, _ := json.MarshalIndent(memInfo, "", "  ")
	fmt.Println(string(b))

	for _, m := range memInfo.Modules {
		b, _ := json.MarshalIndent(m, "", "  ")
		fmt.Println(string(b))
	}
	//fmt.Println(memory.Area)
	//fmt.Println(memory.DefaultHugePageSize)

	subRam := container.NewVBox(
		//System
		widget.NewLabel(memName),
		widget.NewLabel(memDefaulSize),
		widget.NewLabel(memModule),
		widget.NewLabel(memHugeSize),
		widget.NewLabel(memSupportSize),
		widget.NewLabel(memTotalHugeBytes),
		widget.NewLabel(memTotalPhysicalBytes),
		widget.NewLabel(memTotalUsableBytes),
	)
	ram := container.NewVBox(
		//System
		widget.NewCard("Ram", "", subRam),
		widget.NewLabel("***"),
	)
	/*
		subSystem := container.NewVBox(
			//System
			widget.NewLabel("ผู้ผลิต"),
		)

		System := container.NewVBox(
			//System
			widget.NewCard("System", "", subSystem),
			widget.NewLabel("0000"),
		)

		subMainboard := container.NewVBox(
			//mainboard
			widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", x)),
		)

		Mainboard := container.NewVBox(
			//mainboard
			widget.NewCard("Mainboard", "", subMainboard),
		)

		subBIOS_UEFI := container.NewVBox(
			//BIOS/UEFI
			widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", x)),
		)
		BIOS_UEFI := container.NewVBox(
			//BIOS/UEFI
			widget.NewCard("BIOS/UEFI", "", subBIOS_UEFI),
		)

		subChassis := container.NewVBox(
			//Chassis
			widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s")),
		)

		Chassis := container.NewVBox(
			//Chassis
			widget.NewCard("Chassis", "", subChassis),
		)
	*/
	return container.NewAppTabs(
		container.NewTabItem("ram", container.NewScroll(ram)),
		//container.NewTabItem("ram", container.NewScroll(Mainboard)),
		//container.NewTabItem("ram", container.NewScroll(BIOS_UEFI)),
		//container.NewTabItem("ram", container.NewScroll(Chassis)),
	)
}
