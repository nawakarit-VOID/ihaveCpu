// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_mainboardinfo

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jaypipes/ghw"
)

func mainboard_info() map[string]interface{} {
	baseboard, err := ghw.Baseboard()
	if err != nil {
		panic(err)
	}

	//var vendor string

	vendor, vendorjson := fmt.Println(baseboard.Vendor)
	fmt.Println(baseboard.Product)
	fmt.Println(baseboard.Version)
	fmt.Println(baseboard.SerialNumber)

	fmt.Println(baseboard.AssetTag)
	fmt.Println(baseboard.Vendor)

	fmt.Println(baseboard)

	return map[string]interface{}{
		// gopsutil
		"Vendor":     vendor,
		"Vendorjson": vendorjson,
	}

}

func MainboardTabs() fyne.CanvasObject {
	x := mainboard_info()

	mainboard := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s", x["Vendor"])),
		widget.NewLabel(fmt.Sprintf("%s", x["Vendorjson"])),
		widget.NewSeparator(),
	)

	return container.NewAppTabs(
		container.NewTabItem("Overview", container.NewScroll(mainboard)),
	)
}
