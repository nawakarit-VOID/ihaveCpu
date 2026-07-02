// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.

// - pciids.go - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
package pcie

import (
	"bufio"
	"bytes"
	"embed"
	"strings"
)

// - embed - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

//go:embed assets/pci.ids
var pciFS embed.FS

// - โครงสร้างข้อมูล - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
type Vendor struct {
	Name    string
	Devices map[string]string
}

type PCIIDs struct {
	Vendors map[string]*Vendor
}

// - โหลดฐานข้อมูล - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func LoadPCIIDs(data []byte) *PCIIDs {

	db := &PCIIDs{
		Vendors: make(map[string]*Vendor),
	}

	var current *Vendor

	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {

		line := scanner.Text()

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		// Vendor
		if line[0] != '\t' {

			fields := strings.Fields(line)

			if len(fields) < 2 {
				continue
			}

			id := strings.ToLower(fields[0])

			current = &Vendor{
				Name:    strings.Join(fields[1:], " "),
				Devices: make(map[string]string),
			}

			db.Vendors[id] = current

			continue
		}

		// Device
		if current != nil {

			line = strings.TrimLeft(line, "\t")

			fields := strings.Fields(line)

			if len(fields) < 2 {
				continue
			}

			id := strings.ToLower(fields[0])

			current.Devices[id] = strings.Join(fields[1:], " ")
		}
	}

	return db
}

// - ค้นหา - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func (db *PCIIDs) VendorName(id string) string {

	id = strings.TrimPrefix(strings.ToLower(id), "0x")

	if v, ok := db.Vendors[id]; ok {
		return v.Name
	}

	return "Unknown"
}

func (db *PCIIDs) DeviceName(vendorID, deviceID string) string {

	vendorID = strings.TrimPrefix(strings.ToLower(vendorID), "0x")
	deviceID = strings.TrimPrefix(strings.ToLower(deviceID), "0x")

	vendor, ok := db.Vendors[vendorID]
	if !ok {
		return "Unknown"
	}

	if device, ok := vendor.Devices[deviceID]; ok {
		return device
	}

	return "Unknown"
}

// - โหลดครั้งเดียว - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

var PCI *PCIIDs

func InitPCI() {

	data, _ := pciFS.ReadFile("assets/pci.ids")

	PCI = LoadPCIIDs(data)
}
