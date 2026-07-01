// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package package_pcieinfo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PCIeDevice struct {
	Address          string
	VendorID         string
	DeviceID         string
	Class            string
	CurrentLinkSpeed string
	CurrentLinkWidth string
	MaxLinkSpeed     string
	MaxLinkWidth     string
}

func read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func GetPCIeDevices() ([]PCIeDevice, error) {

	const base = "/sys/bus/pci/devices"

	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}

	devices := make([]PCIeDevice, 0)

	for _, entry := range entries {

		if !entry.IsDir() {
			continue
		}

		dir := filepath.Join(base, entry.Name())

		device := PCIeDevice{
			Address:          entry.Name(),
			VendorID:         read(filepath.Join(dir, "vendor")),
			DeviceID:         read(filepath.Join(dir, "device")),
			Class:            read(filepath.Join(dir, "class")),
			CurrentLinkSpeed: read(filepath.Join(dir, "current_link_speed")),
			CurrentLinkWidth: read(filepath.Join(dir, "current_link_width")),
			MaxLinkSpeed:     read(filepath.Join(dir, "max_link_speed")),
			MaxLinkWidth:     read(filepath.Join(dir, "max_link_width")),
		}

		devices = append(devices, device)
	}

	return devices, nil
}

/*
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
*/

var systemsDetailLabel *widget.Label //ประกาศแบบ golbal
func SystemsDetailLabelcmd(text string) {
	if systemsDetailLabel != nil {
		systemsDetailLabel.SetText(text)
	}
}

var pciClass = map[string]string{
	"0x010802": "NVMe",
	"0x010601": "SATA",
	"0x020000": "Ethernet",
	"0x028000": "Wi-Fi",
	"0x030000": "VGA",
	"0x030200": "3D",
	"0x040300": "Audio",
	"0x060400": "PCI Bridge",
	"0x0c0330": "USB xHCI",
	"0x0c0500": "SMBus",
}

func PcieTabs() fyne.CanvasObject {

	systemsDetailLabel = widget.NewLabel("")

	var pcieString string

	devices, err := GetPCIeDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range devices {

		//pcieX += fmt.Sprintln("===================================")
		pcieString += fmt.Sprintln("Address :", d.Address)
		pcieString += fmt.Sprintln("Vendor  :", d.VendorID)
		pcieString += fmt.Sprintln("Device  :", d.DeviceID)
		pcieString += fmt.Sprintln("Class   :", d.Class)
		pcieString += fmt.Sprintln("Current :", d.CurrentLinkSpeed, d.CurrentLinkWidth)
		pcieString += fmt.Sprintln("Max     :", d.MaxLinkSpeed, d.MaxLinkWidth)

		fmt.Println("Address :", d.Address)
		fmt.Println("Vendor  :", d.VendorID)
		fmt.Println("Device  :", d.DeviceID)
		fmt.Println("Class   :", d.Class)
		fmt.Println("Current :", d.CurrentLinkSpeed, d.CurrentLinkWidth)
		fmt.Println("Max     :", d.MaxLinkSpeed, d.MaxLinkWidth)
	}

	pcie := container.NewVBox(widget.NewLabel(pcieString))

	PcieDetail := container.NewVBox(
		//detail
		widget.NewCard("Pcie", "", pcie),
	)

	return container.NewAppTabs(
		//container.NewTabItem("BIOS/UEFI", container.NewScroll(PcieDetail)),
		container.NewTabItem("Detail", container.NewScroll(PcieDetail)),
	)
}
