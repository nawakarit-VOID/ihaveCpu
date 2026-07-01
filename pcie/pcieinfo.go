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

		dir := filepath.Join(base, entry.Name())

		//fmt.Println("กำลังอ่าน:", dir)

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

		fmt.Printf("%+v\n", device)
	}

	return devices, nil
}

func ClassName(class string) string {
	if name, ok := pciClass[class]; ok {
		return name
	}
	return class
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

func VendorName(id string) string {
	if name, ok := pciVendor[id]; ok {
		return name
	}
	return id
}

var pciVendor = map[string]string{
	"0x8086": "Intel",
	"0x10de": "NVIDIA",
	"0x1002": "AMD",
	"0x1022": "AMD",
	"0x144d": "Samsung",
	"0x1b21": "ASMedia",
	"0x14e4": "Broadcom",
	"0x10ec": "Realtek",
	"0x8087": "Intel",
	"0x106b": "Apple",
	"0x13B5": "ARM",
	"0x17CB": "Qualcomm",
	"0x1010": "Imagination Technologies",
	"0x15ad": "VMware",
	"0x1414": "Microsoft",
	"0x1014": "IBM",
	"0x1106": "VIA Technologies",
	"0x1000": "Broadcom",
	"0x1095": "Silicon Image",
	"0x105A": "Promise Technology",
	"0x1D0F": "Amazon",
	"0x1013": "Cirrus Logic",
	"0x1050": "Winbond",
	"0x1234": "Bochs (Emulator)",
	"0x1AF4": "Virtio (Paravirtualization)",
	"0x1B36": "QEMU (Emulator)",
}

func PCIeGeneration(speed string) string {
	switch speed {
	case "2.5 GT/s PCIe":
		return "Gen1"
	case "5.0 GT/s PCIe":
		return "Gen2"
	case "8.0 GT/s PCIe":
		return "Gen3"
	case "16.0 GT/s PCIe":
		return "Gen4"
	case "32.0 GT/s PCIe":
		return "Gen5"
	case "64.0 GT/s PCIe":
		return "Gen6"
	default:
		return "Unknown"
	}
}

func PCIeCanvas() fyne.CanvasObject {

	var pcieString string

	devices, err := GetPCIeDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range devices {

		//pcieX += fmt.Sprintln("===================================")
		pcieString += fmt.Sprintln("Address :", d.Address)

		pcieString += fmt.Sprintf("Vendor : %s (%s)\n",
			VendorName(d.VendorID),
			d.VendorID)

		pcieString += fmt.Sprintln("Device  :", d.DeviceID)

		pcieString += fmt.Sprintf("Class  : %s\n",
			ClassName(d.Class))

		if d.CurrentLinkSpeed != "" {

			pcieString += fmt.Sprintf(
				"PCIe   : %s x%s\n",
				PCIeGeneration(d.CurrentLinkSpeed),
				d.CurrentLinkWidth,
			)

			pcieString += fmt.Sprintf(
				"Speed  : %s\n",
				d.CurrentLinkSpeed,
			)
		} else {
			pcieString += "\n"
		}
	}

	pcie := widget.NewLabel(pcieString)

	return pcie
}

var PcieDetailLabel *widget.Label //ประกาศแบบ golbal
func PcieDetailLabelcmd(text string) {
	if PcieDetailLabel != nil {
		PcieDetailLabel.SetText(text)
	}
}

func PcieTabs() fyne.CanvasObject {

	PcieDetailLabel = widget.NewLabel("")

	pcie := PCIeCanvas()

	PcieX := container.NewVBox(
		//detail
		widget.NewCard("Pcie", "", pcie),
	)

	return container.NewAppTabs(
		container.NewTabItem("Pcie", container.NewScroll(PcieX)),
		//container.NewTabItem("Detail", container.NewScroll(pcie)),
	)
}
