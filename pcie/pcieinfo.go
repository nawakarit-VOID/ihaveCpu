// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package pcie

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

func vendorName(id string) string {
	if PCI == nil {
		return id
	}

	name := PCI.VendorName(id)
	if name == "Unknown" {
		return id
	}

	return name
}

func deviceName(vendorID, deviceID string) string {
	if PCI == nil {
		return deviceID
	}

	name := PCI.DeviceName(vendorID, deviceID)
	if name == "Unknown" {
		return deviceID
	}

	return name
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
			vendorName(d.VendorID),
			d.VendorID)

		pcieString += fmt.Sprintf("Device  : %s (%s)\n",
			deviceName(d.VendorID, d.DeviceID),
			d.DeviceID)

		pcieString += fmt.Sprintf("Class  : %s\n",
			ClassName(d.Class))

		if d.CurrentLinkSpeed != "" {

			pcieString += fmt.Sprintf(
				"PCIe   : %s x%s\n",
				PCIeGeneration(d.CurrentLinkSpeed),
				d.CurrentLinkWidth,
			)

			pcieString += fmt.Sprintf(
				"Speed  : %s\n\n",
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

	InitPCI()
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
