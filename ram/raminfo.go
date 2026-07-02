// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package ram

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

var ramDetailLabel *widget.Label //ประกาศแบบ golbal

func RamDetailLabelcmd(text string) {
	if ramDetailLabel != nil {
		ramDetailLabel.SetText(text)
	}
}

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
	ramDetailLabel = widget.NewLabel("")

	// ============================================================================
	// ghw.Memory
	// ============================================================================

	var total string
	var Support string
	var status string
	var modulesS string

	//ป้องกัน Memory() คืนค่า nil ให้เป็นค่าเริ่มต้นก่อนเข้าถึงฟิลด์
	memInfo := Memory()
	if memInfo == nil {
		memInfo = &ghw.MemoryInfo{}
	}

	//info += fmt.Sprintf("Area: %v\n", &memInfo.Area) //ดูทั้งหมด

	//Total
	//TotalPhysicalBytes
	TotalPhysicalBytes, TotalPhysicalBytesString := newProcessValue(float64(memInfo.TotalPhysicalBytes))
	total += fmt.Sprintf("Ram ทั้งหมด : %.2f %s\n", TotalPhysicalBytes, TotalPhysicalBytesString)
	//TotalUsableBytes
	TotalUsableBytes, TotalUsableBytesString := newProcessValue(float64(memInfo.TotalUsableBytes))
	total += fmt.Sprintf("Ram ที่นำมาใช้งานได้ : %.2f %s\n", TotalUsableBytes, TotalUsableBytesString)

	//การรองรับ Hug Page
	//SupportedPageSizes
	for NoSupported, amount := range memInfo.SupportedPageSizes {
		SupportedPageSizes, SupportedPageSizesString := newProcessValue(float64(amount))
		Support += fmt.Sprintf("ลำดับ %d : ขนาด %.2f %s\n", NoSupported, SupportedPageSizes, SupportedPageSizesString)
	}
	//DefaultHugePageSize
	DefaultHugePageSize, DefaultHugePageSizeString := newProcessValue(float64(memInfo.DefaultHugePageSize))
	Support += fmt.Sprintf("\nค่าเริ่มต้นของ Hug Page คือ %.2f %s\n", DefaultHugePageSize, DefaultHugePageSizeString)

	//สถานะ HugPage
	//HugePageAmountsBySize
	for size, amount := range memInfo.HugePageAmountsBySize {
		HugePageAmountsBySize, HugePageAmountsBySizeString := newProcessValue(float64(size))
		status += fmt.Sprintf("HugePage %.2f %s : สถานะ %d\n", HugePageAmountsBySize, HugePageAmountsBySizeString, amount)
	}
	//TotalHugePageBytes
	TotalHugePageBytes, TotalHugePageBytesString := newProcessValue(float64(memInfo.TotalHugePageBytes))
	status += fmt.Sprintf("\nใช้ Hug Page ไปแล้ว %.2f %s\n", TotalHugePageBytes, TotalHugePageBytesString)

	//ยังไม่แสดง
	for i, m := range memInfo.Modules {
		modulesS += fmt.Sprintf("\nModule %d\n", i+1)
		modulesS += fmt.Sprintf("  Vendor : %s\n", m.Vendor)
		//modules += fmt.Sprintf("  Product: %s\n", m.Product)
		modulesS += fmt.Sprintf("  Size   : %d\n", m.SizeBytes)
		modulesS += fmt.Sprintf("  Serial : %s\n", m.SerialNumber)
		modulesS += fmt.Sprintf("  Serial : %s\n", m.Label)
		modulesS += fmt.Sprintf("  Serial : %s\n", m.Location)
	} //*ไม่ขึ้น

	//ghw
	physical_usable := container.NewVBox(
		widget.NewLabel(total),
	)

	SupportedPage_DefaultHugePage := container.NewVBox(
		/*	widget.NewLabel(
						`*ปกติแล้วระบบปฏิบัติการจะแบ่ง RAM ออกเป็นช่องเล็ก ๆ เรียกว่า [ Page ]
			ขนาดปกติคือ 4 KB แต่ถ้าโปรแกรมต้องใช้ RAM เยอะ ๆ
			เช่น Database หรือแอปพลิเคชันใหญ่ ๆ) การใช้ช่องเล็ก ๆ
			จะทำให้หาข้อมูลช้า ระบบจึงมีฟีเจอร์ Huge Pages
			เพื่อรวมเป็นช่องขนาดใหญ่ขึ้น ทำให้ทำงานเร็วขึ้น`),
		*/
		widget.NewLabel(Support),
	)

	TotalHugePage_HugePageAmounts := container.NewVBox(
		widget.NewLabel(status),
	)

	Modules := container.NewVBox(
		widget.NewLabel(modulesS),
	)

	//dmidecode
	sub_Detail_ram := container.NewVBox(
		ramDetailLabel,
	)

	//card
	Overview := container.NewVBox(
		widget.NewCard("Ram total", "", physical_usable),
		widget.NewCard("การรองรับ Huge Pages", "หน่วยความจำขนาดพิเศษ ", SupportedPage_DefaultHugePage),
		widget.NewCard("สถานะ Huge Pages", "", TotalHugePage_HugePageAmounts),
		widget.NewCard("Modules", "", Modules),
	)

	Detail := container.NewVBox(
		widget.NewCard("Detail", "", sub_Detail_ram),
	)

	return container.NewAppTabs(
		container.NewTabItem("Overview", container.NewScroll(Overview)),
		container.NewTabItem("Detail", container.NewScroll(Detail)),
		//container.NewTabItem("ram", container.NewScroll(BIOS_UEFI)),
		//container.NewTabItem("ram", container.NewScroll(Chassis)),
	)
}
