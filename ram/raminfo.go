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
)

func Memory() *ghw.MemoryInfo {
	//ข้อมูล RAM
	memInfo, err := ghw.Memory()
	if err != nil {
		return nil
	}

	return memInfo
}

/*
func g() *ghw.GPUInfo {
	gpuInfo, err := ghw.GPU()
	if err != nil {
		return nil
	}
	return gpuInfo

}
*/

// dmidecode
func RamTabs() fyne.CanvasObject {
	var info string
	memInfo := Memory()

	info += fmt.Sprintf("Area: %v\n", memInfo.Area)
	info += fmt.Sprintf("TotalPhysicalBytes: %d\n", memInfo.TotalPhysicalBytes)
	info += fmt.Sprintf("TotalUsableBytes: %d\n", memInfo.TotalUsableBytes)
	info += fmt.Sprintf("DefaultHugePageSize: %d\n", memInfo.DefaultHugePageSize)

	memName := fmt.Sprintf("1 memName : %s", &memInfo.Area)
	//c := memInfo.Area
	memDefaulSize := fmt.Sprintf("2 memDefaulSize : %d", memInfo.DefaultHugePageSize) //เปลี่ยนแปลงตลอด

	//memModule := fmt.Println(len(memInfo.Modules))                //[]
	/*
		var memModule string
		for _, area := range memInfo.Modules {
			m := fmt.Printf("%+v\n", area)
			memModule += m
		}
	*/
	memHugeSize := fmt.Sprintf("4 memHugeSize : %v\n", memInfo.HugePageAmountsBySize) //map uint64

	memSupportSize := fmt.Sprintf("5 memSupportSize : %T", &memInfo.SupportedPageSizes) //[] uint64
	memTotalHugeBytes := fmt.Sprintf("6 memTotalHugeBytes : %d", memInfo.TotalHugePageBytes)
	memTotalPhysicalBytes := fmt.Sprintf("7 memTotalPhysicalBytes : %d", memInfo.TotalPhysicalBytes)
	memTotalUsableBytes := fmt.Sprintf("8 memTotalUsableBytes : %d", memInfo.TotalUsableBytes)

	//fmt.Println(c)
	//memDefaulSize1 := c / 1024 //byte
	//println("kbyte", memDefaulSize1)
	//	var xx string
	/*
		for _, area := range memInfo.Modules {
			fmt.Printf("%+v\n", area)
		}
	*/
	/*
		for _, m := range memInfo.Modules {
			b, _ := json.MarshalIndent(m, "", "  ")
			fmt.Println(string(b))
		}
	*/
	/*
		b, _ := json.MarshalIndent(memInfo, "", "  ")
		fmt.Println(string(b))
	*/
	//fmt.Println(memory.Area)
	//fmt.Println(memory.DefaultHugePageSize)
	var x string
	x += fmt.Sprintf("%#v\n", memInfo)

	subRam := container.NewVBox(
		//System
		widget.NewLabel(memName),
		widget.NewLabel(memDefaulSize),
		//widget.NewLabel(fmt.Sprintln("%s", memModule)),
		widget.NewLabel(memHugeSize),
		widget.NewLabel(memSupportSize),
		widget.NewLabel(memTotalHugeBytes),
		widget.NewLabel(memTotalPhysicalBytes),
		widget.NewLabel(memTotalUsableBytes),
		widget.NewLabel(x),
		widget.NewLabel(info),
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
