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
	var info string
	memInfo := Memory()

	//info += fmt.Sprintf("Area: %v\n", memInfo.Area) //**

	//info += fmt.Sprintf("TotalPhysicalBytes: %d\n", memInfo.TotalPhysicalBytes)
	TotalPhysicalBytes, TotalPhysicalBytesString := newProcessValue(float64(memInfo.TotalPhysicalBytes))
	info += fmt.Sprintf("[ RAM ทั้งหมดที่มีในเครื่อง ] : %.2f %s\n", TotalPhysicalBytes, TotalPhysicalBytesString)
	//info += fmt.Sprintf("TotalUsableBytes: %d\n", memInfo.TotalUsableBytes)
	TotalUsableBytes, TotalUsableBytesString := newProcessValue(float64(memInfo.TotalUsableBytes))
	info += fmt.Sprintf("[ RAM ที่ระบบสามารถนำไปใช้งานได้จริง ] : %.2f %s\n", TotalUsableBytes, TotalUsableBytesString)

	//info += fmt.Sprintf("DefaultHugePageSize: %d\n", memInfo.DefaultHugePageSize)
	DefaultHugePageSize, DefaultHugePageSizeString := newProcessValue(float64(memInfo.DefaultHugePageSize))
	info += fmt.Sprintf("[ หากเปิดใช้ Huge Pages ขนาดเริ่มต้น คือ ] : %.2f %s\n", DefaultHugePageSize, DefaultHugePageSizeString)

	//

	for i, m := range memInfo.Modules {
		info += fmt.Sprintf("\nModule %d\n", i+1)
		info += fmt.Sprintf("  Vendor : %s\n", m.Vendor)
		//info += fmt.Sprintf("  Product: %s\n", m.Product)
		info += fmt.Sprintf("  Size   : %d\n", m.SizeBytes)
		info += fmt.Sprintf("  Serial : %s\n", m.SerialNumber)
		info += fmt.Sprintf("  Serial : %s\n", m.Label)
		info += fmt.Sprintf("  Serial : %s\n", m.Location)
	} //*ไม่ขึ้น

	//
	info += "การรองรับ Huge Pages (หน้าหน่วยความจำขนาดพิเศษ) อธิบายเสริม: ปกติแล้วระบบปฏิบัติการจะแบ่ง RAM ออกเป็นช่องเล็ก ๆ เรียกว่า [ Page ] ขนาดปกติคือ 4 KB แต่ถ้าโปรแกรมต้องใช้ RAM เยอะ ๆ (เช่น Database หรือแอปพลิเคชันใหญ่ ๆ) การใช้ช่องเล็ก ๆ จะทำให้หาข้อมูลช้า ระบบจึงมีฟีเจอร์ Huge Pages เพื่อรวมเป็นช่องขนาดใหญ่ขึ้น ทำให้ทำงานเร็วขึ้น\n"
	info += "ระบบนี้รองรับการแบ่งหน้า RAM ขนาดใหญ่ 2 ขนาด คือ\n"
	for size, amount := range memInfo.HugePageAmountsBySize {
		HugePageAmountsBySize, HugePageAmountsBySizeString := newProcessValue(float64(size))
		info += fmt.Sprintf("HugePage %.2f %s : สถานะ %d\n", HugePageAmountsBySize, HugePageAmountsBySizeString, amount)
	}

	info += "ระบบนี้รองรับการแบ่งหน้า RAM\n"

	for no, amount := range memInfo.SupportedPageSizes {
		HugePageAmountsBySize, HugePageAmountsBySizeString := newProcessValue(float64(amount))
		info += fmt.Sprintf("ลำดับ %d : ขนาด %.2f %s\n", no, HugePageAmountsBySize, HugePageAmountsBySizeString)
	}

	TotalHugePageBytes, TotalHugePageBytesString := newProcessValue(float64(memInfo.TotalHugePageBytes))
	info += fmt.Sprintf("TotalHugePageBytes ขนาด %.2f %s\n", TotalHugePageBytes, TotalHugePageBytesString)

	for _, m := range memInfo.Modules {
		fmt.Printf("%+v\n", m)
	} //*ไม่ขึ้น

	//

	subRam := container.NewVBox(
		//System
		//widget.NewLabel(memName),
		//widget.NewLabel(memDefaulSize),
		//widget.NewLabel(fmt.Sprintln("%s", memModule)),
		//widget.NewLabel(memHugeSize),
		//widget.NewLabel(memSupportSize),
		//widget.NewLabel(memTotalHugeBytes),
		//widget.NewLabel(memTotalPhysicalBytes),
		//widget.NewLabel(memTotalUsableBytes),
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
