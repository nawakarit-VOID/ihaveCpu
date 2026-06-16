// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_raminfo

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/memory"
)

func Memory() *memory.Info {
	//ข้อมูล RAM
	memory, err := ghw.Memory()
	if err != nil {
		return nil
	}
	return memory
}

func RamTabs() fyne.CanvasObject {

	memory := Memory()

	memName := fmt.Sprintf("%s", &memory.Area)
	memDefaulSize := fmt.Sprintf("%d", &memory.DefaultHugePageSize) //เปลี่ยนแปลงตลอด
	memModule := fmt.Sprintf("%T", &memory.Modules)                 //[]
	memHugeSize := fmt.Sprintf("%T", &memory.HugePageAmountsBySize) //map uint64
	memSupportSize := fmt.Sprintf("%T", &memory.SupportedPageSizes) //[] uint64
	memTotalHugeBytes := fmt.Sprintf("%d", &memory.TotalHugePageBytes)
	memTotalPhysicalBytes := fmt.Sprintf("%d", &memory.TotalPhysicalBytes)
	memTotalUsableBytes := fmt.Sprintf("%d", &memory.TotalUsableBytes)

	//	var xx string
	for _, area := range memory.Modules {
		fmt.Printf("%+v\n", area)
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
