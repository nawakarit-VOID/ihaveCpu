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

func newProcessValue(value float64) (float64, string) {
	// ตัวอักษร flag ที่สัมผัส
	var x string = "B" //8Bit = 1Byte
	// ตรวจสอบเงื่อนไข //แบบ บนลงล่าง
	if value >= 1024 {
		value = value / 1024
		x = "KB"
		if value >= 1024 {
			value = value / 1024
			x = "MB"
			if value >= 1024 {
				value = value / 1024
				x = "GB"
				if value >= 1024 {
					value = value / 1024
					x = "TB"
					if value >= 1024 {
						value = value / 1024
						x = "PB"
						if value >= 1024 {
							value = value / 1024
							x = "EB"
							if value >= 1024 {
								value = value / 1024
								x = "ZB"
								if value >= 1024 {
									value = value / 1024
									x = "YB"
									if value >= 1024 {
										value = value / 1024
										x = "Bronto Byte"
										if value >= 1024 {
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
	return value, x
}

// dmidecode
func RamTabs() fyne.CanvasObject {

	var physical string
	var usable string
	var SupportedPage string
	var DefaultHugePage string
	var HugePageAmounts string
	var TotalHugePage string
	var modules string

	memInfo := Memory()

	//info += fmt.Sprintf("Area: %v\n", memInfo.Area)

	//info += fmt.Sprintf("TotalPhysicalBytes: %d\n", memInfo.TotalPhysicalBytes)
	TotalPhysicalBytes, TotalPhysicalBytesString := newProcessValue(float64(memInfo.TotalPhysicalBytes))
	physical += fmt.Sprintf("Ram physical total : %.2f %s", TotalPhysicalBytes, TotalPhysicalBytesString)

	//info += fmt.Sprintf("TotalUsableBytes: %d\n", memInfo.TotalUsableBytes)
	TotalUsableBytes, TotalUsableBytesString := newProcessValue(float64(memInfo.TotalUsableBytes))
	usable += fmt.Sprintf("Ram usable totsl : %.2f %s", TotalUsableBytes, TotalUsableBytesString)

	for NoSupported, amount := range memInfo.SupportedPageSizes {
		SupportedPageSizes, SupportedPageSizesString := newProcessValue(float64(amount))
		SupportedPage += fmt.Sprintf("ลำดับ %d : ขนาด %.2f %s\n", NoSupported, SupportedPageSizes, SupportedPageSizesString)
	}

	//info += fmt.Sprintf("DefaultHugePageSize: %d\n", memInfo.DefaultHugePageSize)
	DefaultHugePageSize, DefaultHugePageSizeString := newProcessValue(float64(memInfo.DefaultHugePageSize))
	DefaultHugePage += fmt.Sprintf("Default Huge Page size : %.2f %s", DefaultHugePageSize, DefaultHugePageSizeString)

	for size, amount := range memInfo.HugePageAmountsBySize {
		HugePageAmountsBySize, HugePageAmountsBySizeString := newProcessValue(float64(size))
		HugePageAmounts += fmt.Sprintf("HugePage %.2f %s : สถานะ %d\n", HugePageAmountsBySize, HugePageAmountsBySizeString, amount)
	}

	//info += "ระบบนี้รองรับการแบ่งหน้า RAM\n"

	TotalHugePageBytes, TotalHugePageBytesString := newProcessValue(float64(memInfo.TotalHugePageBytes))
	TotalHugePage += fmt.Sprintf("TotalHugePageBytes ขนาด %.2f %s\n", TotalHugePageBytes, TotalHugePageBytesString)

	/*
		info += `Ram physical total = แรมทั้งหมด
		Ram usable totsl = แรมที่ระบบสามารถใช้งานได้
		HugePage = หน้าหน่วยความจำขนาดพิเศษ

		Default Huge Page size = ระบบจะเลือกขนาดเริ่มต้น หากมีการเปิดใช้ Huge Pages

		`
	*/
	/*
		for _, m := range memInfo.Modules {
			fmt.Printf("%+v\n", m)
		} //*ไม่ขึ้น
	*/
	//หรือ
	for i, m := range memInfo.Modules {
		modules += fmt.Sprintf("\nModule %d\n", i+1)
		modules += fmt.Sprintf("  Vendor : %s\n", m.Vendor)
		//modules += fmt.Sprintf("  Product: %s\n", m.Product)
		modules += fmt.Sprintf("  Size   : %d\n", m.SizeBytes)
		modules += fmt.Sprintf("  Serial : %s\n", m.SerialNumber)
		modules += fmt.Sprintf("  Serial : %s\n", m.Label)
		modules += fmt.Sprintf("  Serial : %s\n", m.Location)
	} //*ไม่ขึ้น

	/*หรือ ดูทั้งหมด
	b, err := json.MarshalIndent(memInfo, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	x := string(b)
	fmt.Println(x)
	info += fmt.Sprintf("%s", x)
	*/

	physical_usable := container.NewVBox(
		widget.NewLabel(physical),
		widget.NewLabel(usable),
	)

	SupportedPage_DefaultHugePage := container.NewVBox(
		widget.NewLabel(SupportedPage),
		widget.NewLabel(DefaultHugePage),
	)

	TotalHugePage_HugePageAmounts_ := container.NewVBox(
		widget.NewLabel(TotalHugePage),
		widget.NewLabel(HugePageAmounts),
	)

	Modules := container.NewVBox(
		widget.NewLabel(modules),
		//widget.NewLabel(""),
	)
	detail := container.NewVBox(
		//
		widget.NewCard("Ram total", "", physical_usable),
		widget.NewCard("การรองรับ Huge Pages", "หน่วยความจำขนาดพิเศษ", SupportedPage_DefaultHugePage),
		widget.NewCard("สถานะ Huge Pages", "", TotalHugePage_HugePageAmounts_),
		widget.NewCard("Modules", "", Modules),
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
		container.NewTabItem("Detail", container.NewScroll(detail)),
		//container.NewTabItem("ram", container.NewScroll(Mainboard)),
		//container.NewTabItem("ram", container.NewScroll(BIOS_UEFI)),
		//container.NewTabItem("ram", container.NewScroll(Chassis)),
	)
}
