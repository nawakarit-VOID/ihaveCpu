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

	//system
	sys_vendor := read("/sys/class/dmi/id/sys_vendor")
	product_name := read("/sys/class/dmi/id/product_name")
	product_family := read("/sys/class/dmi/id/product_family")
	product_version := read("/sys/class/dmi/id/product_version")
	//product_serial := read("/sys/class/dmi/id/product_serial")
	//product_uuid := read("/sys/class/dmi/id/product_uuid")
	product_sku := read("/sys/class/dmi/id/product_sku")

	//mainboard
	board_vendor := read("/sys/class/dmi/id/board_vendor")
	board_name := read("/sys/class/dmi/id/board_name")
	board_version := read("/sys/class/dmi/id/board_version")
	//board_serial := read("/sys/class/dmi/id/board_serial")
	board_asset_tag := read("/sys/class/dmi/id/board_asset_tag")

	//BIOS/UEFI
	bios_vendor := read("/sys/class/dmi/id/bios_vendor")
	bios_version := read("/sys/class/dmi/id/bios_version")
	bios_date := read("/sys/class/dmi/id/bios_date")
	bios_release := read("/sys/class/dmi/id/bios_release")

	//Chassis
	chassis_vendor := read("/sys/class/dmi/id/chassis_vendor")
	chassis_type := read("/sys/class/dmi/id/chassis_type")
	//chassis_serial := read("/sys/class/dmi/id/chassis_serial")
	chassis_version := read("/sys/class/dmi/id/chassis_version")
	chassis_asset_tag := read("/sys/class/dmi/id/chassis_asset_tag")

	modalias := read("/sys/class/dmi/id/modalias")

	//-----------------------------------------------------------------------//

	//System
	fmt.Println(sys_vendor)      //ผู้ผลิตเครื่องทั้งเครื่อง (OEM)
	fmt.Println(product_name)    //รุ่นของเครื่อง
	fmt.Println(product_family)  //ตระกูลของเครื่อง
	fmt.Println(product_version) //เวอร์ชันหรือ Revision ของรุ่นเครื่อง
	//fmt.Println(product_serial)//Serial Number ของเครื่องทั้งเครื่อง
	//fmt.Println(product_uuid)//UUID ของเครื่อง
	fmt.Println(product_sku) //SKU/Part Number ของเครื่อง

	//Mainboard
	fmt.Println(board_vendor)  //ผู้ผลิตเมนบอร์ด
	fmt.Println(board_name)    //รุ่นเมนบอร์ด
	fmt.Println(board_version) //Revision/Version ของเมนบอร์ด
	//fmt.Println(board_serial)//Serial Number ของเมนบอร์ด
	fmt.Println(board_asset_tag) //รหัสทรัพย์สิน (Asset Tag) ของเมนบอร์ด ใช้ในองค์กร

	//BIOS/UEFI
	fmt.Println(bios_vendor)  //ผู้ผลิต BIOS
	fmt.Println(bios_version) //เวอร์ชัน BIOS
	fmt.Println(bios_date)    //วันที่ออก BIOS
	fmt.Println(bios_release) //เวอร์ชัน Release ของ BIOS ตาม SMBIOS

	//Chassis
	fmt.Println(chassis_vendor) //ผู้ผลิตตัวเครื่อง/เคส
	fmt.Println(chassis_type)   //ประเภทของเครื่อง
	//fmt.Println(chassis_serial)//Serial Number ของตัวเครื่อง/เคส
	fmt.Println(chassis_version)   //รุ่นหรือ Revision ของตัวเครื่อง
	fmt.Println(chassis_asset_tag) //Asset Tag ของตัวเครื่อง

	fmt.Println(modalias) //Hardware ID สำหรับ kernel ใช้จับคู่ driver

	return map[string]interface{}{
		"Sys_vendor":      sys_vendor,      //ผู้ผลิตเครื่องทั้งเครื่อง (OEM)
		"Product_name":    product_name,    //รุ่นของเครื่อง
		"Product_family":  product_family,  //ตระกูลของเครื่อง
		"Product_version": product_version, //เวอร์ชันหรือ Revision ของรุ่นเครื่อง
		//"Product_serial":    product_serial,    //Serial Number ของเครื่องทั้งเครื่อง
		//"Product_uuid":      product_uuid,      //UUID ของเครื่อง
		"Product_sku":   product_sku,   //SKU/Part Number ของเครื่อง
		"Board_vendor":  board_vendor,  //ผู้ผลิตเมนบอร์ด
		"Board_name":    board_name,    //รุ่นเมนบอร์ด
		"Board_version": board_version, //Revision/Version ของเมนบอร์ด
		//"Board_serial":      board_serial,      //Serial Number ของเมนบอร์ด
		"Board_asset_tag": board_asset_tag, //รหัสทรัพย์สิน (Asset Tag) ของเมนบอร์ด ใช้ในองค์กร
		"Bios_vendor":     bios_vendor,     //ผู้ผลิต BIOS
		"Bios_version":    bios_version,    //เวอร์ชัน BIOS
		"Bios_date":       bios_date,       //วันที่ออก BIOS
		"Bios_release":    bios_release,    //เวอร์ชัน Release ของ BIOS ตาม SMBIOS
		"Chassis_vendor":  chassis_vendor,  //ผู้ผลิตตัวเครื่อง/เคส
		"Chassis_type":    chassis_type,    //ประเภทของเครื่อง
		//"Chassis_serial":    chassis_serial,    //Serial Number ของตัวเครื่อง/เคส
		"Chassis_version":   chassis_version,   //รุ่นหรือ Revision ของตัวเครื่อง
		"Chassis_asset_tag": chassis_asset_tag, //Asset Tag ของตัวเครื่อง
		"Modalias":          modalias,          //Hardware ID สำหรับ kernel ใช้จับคู่ driver
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
