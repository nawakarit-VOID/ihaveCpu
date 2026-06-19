// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_raminfo

import (
	"fmt"
	"os/exec"
	"strings"

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

type RAMModule struct {
	Size         string
	FormFactor   string
	Locator      string
	BankLocator  string
	Type         string
	Speed        string
	ConfigSpeed  string
	Manufacturer string
	PartNumber   string
	SerialNumber string
}

// เรียก dmidecode  แบบ 2
func readDMI() (string, error) {
	out, err := exec.Command("sudo", "dmidecode", "-t", "memory").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// เอาผล readDMI มาแยก
func ParseMemory(text string) []RAMModule {
	var modules []RAMModule
	var ram RAMModule

	lines := strings.Split(text, "\n")

	for _, line := range lines {

		line = strings.TrimSpace(line)

		switch {

		case line == "Memory Device":
			// เก็บตัวก่อนหน้า
			if ram.Size != "" {
				modules = append(modules, ram)
			}
			ram = RAMModule{}

		case strings.HasPrefix(line, "Size:"):
			ram.Size = strings.TrimSpace(strings.TrimPrefix(line, "Size:"))

		case strings.HasPrefix(line, "Form Factor:"):
			ram.FormFactor = strings.TrimSpace(strings.TrimPrefix(line, "Form Factor:"))

		case strings.HasPrefix(line, "Locator:"):
			ram.Locator = strings.TrimSpace(strings.TrimPrefix(line, "Locator:"))

		case strings.HasPrefix(line, "Bank Locator:"):
			ram.BankLocator = strings.TrimSpace(strings.TrimPrefix(line, "Bank Locator:"))

		case strings.HasPrefix(line, "Type:"):
			ram.Type = strings.TrimSpace(strings.TrimPrefix(line, "Type:"))

		case strings.HasPrefix(line, "Speed:"):
			ram.Speed = strings.TrimSpace(strings.TrimPrefix(line, "Speed:"))

		case strings.HasPrefix(line, "Configured Memory Speed:"):
			ram.ConfigSpeed = strings.TrimSpace(strings.TrimPrefix(line, "Configured Memory Speed:"))

		case strings.HasPrefix(line, "Manufacturer:"):
			ram.Manufacturer = strings.TrimSpace(strings.TrimPrefix(line, "Manufacturer:"))

		case strings.HasPrefix(line, "Part Number:"):
			ram.PartNumber = strings.TrimSpace(strings.TrimPrefix(line, "Part Number:"))

		case strings.HasPrefix(line, "Serial Number:"):
			ram.SerialNumber = strings.TrimSpace(strings.TrimPrefix(line, "Serial Number:"))
		}
	}

	// ตัวสุดท้าย
	if ram.Size != "" {
		modules = append(modules, ram)
	}

	return modules
}

// เรียก dmidecode  แบบ 1
func GetMemoryInfo() (string, error) {
	cmd := exec.Command("sudo", "dmidecode", "-t", "memory")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
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

	// ============================================================================
	// SECTION_NAME
	// ============================================================================
	//แบบ 1

	teXt, err := GetMemoryInfo()
	if err != nil {
		teXt = err.Error()
	}

	entry := widget.NewLabel("")
	entry.SetText(teXt)
	//entry.Disable()

	//แบบ 2
	var ram_dmidecode string

	text, err := readDMI()
	if err != nil {
		panic(err)
	}

	modules := ParseMemory(text)

	for _, m := range modules {
		ram_dmidecode += fmt.Sprintf("%+v\n", m)
	}

	// ============================================================================
	// ghw.Memory
	// ============================================================================

	var total string
	var Support string
	var status string
	var modulesS string

	memInfo := Memory()

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

	/*
		for _, m := range memInfo.Modules {
			fmt.Printf("%+v\n", m)
		} //*ไม่ขึ้น
	*/
	//หรือ
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

	/*หรือ ดูทั้งหมด
	b, err := json.MarshalIndent(memInfo, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	x := string(b)
	fmt.Println(x)
	info += fmt.Sprintf("%s", x)
	*/
	//

	//dmidecode

	overview := container.NewVBox(
		widget.NewLabel(ram_dmidecode),
		entry,
	)

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
	detail := container.NewVBox(
		//
		widget.NewCard("Ram total", "", physical_usable),
		widget.NewCard("การรองรับ Huge Pages", "หน่วยความจำขนาดพิเศษ ", SupportedPage_DefaultHugePage),
		widget.NewCard("สถานะ Huge Pages", "", TotalHugePage_HugePageAmounts),
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
		container.NewTabItem("Overview", container.NewScroll(overview)),
		container.NewTabItem("Detail", container.NewScroll(detail)),
		//container.NewTabItem("ram", container.NewScroll(Mainboard)),
		//container.NewTabItem("ram", container.NewScroll(BIOS_UEFI)),
		//container.NewTabItem("ram", container.NewScroll(Chassis)),
	)
}
