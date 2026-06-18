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

/*
	func g() *ghw.GPUInfo {
		gpuInfo, err := ghw.GPU()
		if err != nil {
			return nil
		}
		return gpuInfo

}
*/
func processValue(value float64) (float64, string) {
	// ตัวอักษร flag ที่สัมผัส
	var x string = ""
	// ตรวจสอบเงื่อนไข //แบบ บนลงล่าง
	switch {
	case value >= 1099511627776:
		value = value / 1099511627776
		x = "TB"
	case value >= 1073741824:
		value = value / 1073741824
		x = "GB"
	case value >= 1048576:
		value = value / 1048576
		x = "MB"
	case value >= 1000:
		value = value / 1024
		x = "KB"
	default:
		x = "B"
	}
	return value, x
}

func newProcessValue(value float64) (float64, string) {
	// ตัวอักษร flag ที่สัมผัส
	var x string = "B" //8Bit = 1Byte
	// ตรวจสอบเงื่อนไข //แบบ บนลงล่าง
	if value > 1024 {
		value = value / 1024
		x = "KB"
		if value > 1024 {
			value = value / 1024
			x = "MB"
			if value > 1024 {
				value = value / 1024
				x = "GB"
				if value > 1024 {
					value = value / 1024
					x = "TB"
					if value > 1024 {
						value = value / 1024
						x = "PB"
						if value > 1024 {
							value = value / 1024
							x = "EB"
							if value > 1024 {
								value = value / 1024
								x = "ZB"
								if value > 1024 {
									value = value / 1024
									x = "YB"
									if value > 1024 {
										value = value / 1024
										x = "Bronto Byte"
										if value > 1024 {
											value = value / 1024
											x = "Geop Byte"

										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	/*
	   case value >= 1099511627776:
	   		value = value / 1099511627776
	   		x = "TB"
	   	case value >= 1073741824:
	   		value = value / 1073741824
	   		x = "GB"
	   	case value >= 1048576:
	   		value = value / 1048576
	   		x = "MB"
	   	case value >= 1000:
	   		value = value / 1024
	   		x = "KB"
	   	default:
	   		x = "B"
	   	}*/
	return value, x
}

// dmidecode
func RamTabs() fyne.CanvasObject {
	var info string
	memInfo := Memory()

	info += fmt.Sprintf("Area: %v\n", memInfo.Area)

	q1, q2 := newProcessValue(float64(memInfo.TotalPhysicalBytes))

	info += fmt.Sprintf("newProcessValue : %.2f  %s\n", q1, q2)

	x1 := float64(memInfo.TotalPhysicalBytes)
	x2 := float64(x1) / 1024
	x3 := float64(x2) / 1024
	x4 := float64(x3) / 1024

	a, a1 := processValue(x1)

	//info += fmt.Sprintf("TotalPhysicalBytes: %d\n", memInfo.TotalPhysicalBytes)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f byte\n", x1)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f Kb\n", x2)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f Mb\n", x3)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f Gb\n", x4)

	info += fmt.Sprintf("processValue: %.2f %s\n", a, a1)

	xx1 := float64(memInfo.TotalUsableBytes) //b
	xx2 := float64(xx1) / 1024               //Kb
	xx3 := float64(xx2) / 1024               //Mb
	xx4 := float64(xx3) / 1024               //Gb

	b, b1 := processValue(x1)

	//info += fmt.Sprintf("TotalUsableBytes: %d\n", memInfo.TotalUsableBytes)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f byte\n", xx1)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f Kb\n", xx2)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f Mb\n", xx3)
	info += fmt.Sprintf("TotalPhysicalBytes: %.2f Gb\n", xx4) //มีจุดทศนิยม
	info += fmt.Sprintf("processValue: %.2f %s\n", b, b1)     //ถูกปัดเศษขึ้น

	xz1 := float64(memInfo.DefaultHugePageSize)
	xz2 := float64(xz1) / 1024
	xz3 := float64(xz2) / 1024
	//info += fmt.Sprintf("DefaultHugePageSize: %d\n", memInfo.DefaultHugePageSize)
	info += fmt.Sprintf("DefaultHugePageSize: %.2f Mb\n", xz3)

	for i, m := range memInfo.Modules {
		info += fmt.Sprintf("\nModule %d\n", i+1)
		info += fmt.Sprintf("  Vendor : %s\n", m.Vendor)
		//info += fmt.Sprintf("  Product: %s\n", m.Product)
		info += fmt.Sprintf("  Size   : %d\n", m.SizeBytes)
		info += fmt.Sprintf("  Serial : %s\n", m.SerialNumber)
		info += fmt.Sprintf("  Serial : %s\n", m.Label)
		info += fmt.Sprintf("  Serial : %s\n", m.Location)
	} //*ไม่ขึ้น

	for size, amount := range memInfo.HugePageAmountsBySize {
		var xx1 float64
		var xx2 float64
		var xx3 float64
		var xx4 float64

		xx1 = float64(size)
		if xx1 > 1024 {
			xx2 = float64(xx1) / 1024 //Kb
		}
		if xx2 > 1024 {
			xx3 = float64(xx1) / 1024 //Mb
		}
		if xx3 > 1024 {
			xx4 = float64(xx1) / 1024 //Gb
		}

		//xx3 = float64(xx2) / 1024 //Mb
		//xx4 = float64(xx3) / 1024 //Gb

		info += fmt.Sprintf("HugePage %.2f Gb : %d\n", xx4, amount)
	}

	for _, m := range memInfo.Modules {
		fmt.Printf("%+v\n", m)
	} //*ไม่ขึ้น

	c, _ := json.MarshalIndent(memInfo.Modules, "", "  ")
	fmt.Println(string(c)) //*ไม่ขึ้น

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
		container.NewTabItem("Ram", container.NewScroll(ram)),
		//container.NewTabItem("ram", container.NewScroll(Mainboard)),
		//container.NewTabItem("ram", container.NewScroll(BIOS_UEFI)),
		//container.NewTabItem("ram", container.NewScroll(Chassis)),
	)
}
