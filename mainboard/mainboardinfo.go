// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_mainboardinfo

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func read(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func mainboard_info() map[string]interface{} {

	//vendor := read("/sys/class/dmi/id/board_vendor")
	//name := read("/sys/class/dmi/id/board_name")
	//version := read("/sys/class/dmi/id/board_version")

	bios_date := read("/sys/class/dmi/id/bios_date")
	bios_release := read("/sys/class/dmi/id/bios_release")
	bios_vendor := read("/sys/class/dmi/id/bios_vendor")
	bios_version := read("/sys/class/dmi/id/bios_version")

	board_asset_tag := read("/sys/class/dmi/id/board_asset_tag")
	board_name := read("/sys/class/dmi/id/board_name")
	board_serial := read("/sys/class/dmi/id/board_serial")
	board_vendor := read("/sys/class/dmi/id/board_vendor")
	board_version := read("/sys/class/dmi/id/board_version")

	chassis_asset_tag := read("/sys/class/dmi/id/chassis_asset_tag")
	chassis_serial := read("/sys/class/dmi/id/chassis_serial")
	chassis_type := read("/sys/class/dmi/id/chassis_type")
	chassis_vendor := read("/sys/class/dmi/id/chassis_vendor")
	chassis_version := read("/sys/class/dmi/id/chassis_version")

	modalias := read("/sys/class/dmi/id/modalias")

	product_family := read("/sys/class/dmi/id/product_family")
	product_name := read("/sys/class/dmi/id/product_name")
	product_serial := read("/sys/class/dmi/id/product_serial")
	product_sku := read("/sys/class/dmi/id/product_sku")
	product_uuid := read("/sys/class/dmi/id/product_uuid")
	product_version := read("/sys/class/dmi/id/product_version")

	sys_vendor := read("/sys/class/dmi/id/sys_vendor")
	//uevent := read("/sys/class/dmi/id/uevent")

	/*	baseboard, err := ghw.Baseboard()
		if err != nil {
			panic(err)
		}
	*/
	//var vendor string
	//gp := widget.NewLabel(fmt.Println(baseboard.Vendor))
	//vendor, vendorjson := fmt.Println(baseboard.Vendor) //vendorjson

	//fmt.Println(baseboard.Vendor)
	//fmt.Println(baseboard.Product)
	//fmt.Println(baseboard.Version)
	//fmt.Println(baseboard.SerialNumber)

	//fmt.Println(baseboard.AssetTag)
	//fmt.Println(baseboard.Vendor)
	//fmt.Println(baseboard)
	//fmt.Println(vendor)
	//fmt.Println(name)
	//fmt.Println(version)

	fmt.Println(bios_date)
	fmt.Println(bios_release)
	fmt.Println(bios_vendor)
	fmt.Println(bios_version)

	fmt.Println(board_asset_tag)
	fmt.Println(board_name)
	fmt.Println(board_serial)
	fmt.Println(board_vendor)
	fmt.Println(board_version)

	fmt.Println(chassis_asset_tag)
	fmt.Println(chassis_serial)
	fmt.Println(chassis_type)
	fmt.Println(chassis_vendor)
	fmt.Println(chassis_version)
	fmt.Println(modalias)

	fmt.Println(product_family)
	fmt.Println(product_name)
	fmt.Println(product_serial)
	fmt.Println(product_sku)
	fmt.Println(product_uuid)
	fmt.Println(product_version)
	fmt.Println(sys_vendor)
	//fmt.Println(uevent)

	/*
		data, err := json.MarshalIndent(baseboard, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(data))
	*/

	return map[string]interface{}{

		//"Vendor":     vendor,
		//"Vendorjson": vendorjson,
		//"G":          g,
		//	"GP":         gp,
	}

}

func MainboardTabs() fyne.CanvasObject {
	x := mainboard_info()

	mainboard := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s", x["Vendor"])),
		widget.NewLabel(fmt.Sprintf("%s", x["Vendorjson"])),
		widget.NewLabel(fmt.Sprintf("%s", x["G"])),
		widget.NewSeparator(),
	)

	return container.NewAppTabs(
		container.NewTabItem("Overview", container.NewScroll(mainboard)),
	)
}
