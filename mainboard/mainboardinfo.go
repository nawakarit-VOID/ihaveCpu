// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_mainboardinfo

import (
	"fmt"
	"os"
	"os/exec"
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

func GetMainboardInfo() (string, error) {
	// เปลี่ยนจาก "sudo" เป็น "pkexec"
	cmd := exec.Command("pkexec", "dmidecode", "-t", "2")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func GetInfoPkexecAll() (string, error) {
	types := []string{"2", "memory"}
	args := []string{"dmidecode"}
	for _, t := range types {
		args = append(args, "-t", t)
	}

	cmd := exec.Command("pkexec", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
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
	/*
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
	*/
	//-----------------------------------------------------------------------//

	return map[string]interface{}{
		//System
		"Sys_vendor":      sys_vendor,      //ผู้ผลิตเครื่องทั้งเครื่อง (OEM)
		"Product_name":    product_name,    //รุ่นของเครื่อง
		"Product_family":  product_family,  //ตระกูลของเครื่อง
		"Product_version": product_version, //เวอร์ชันหรือ Revision ของรุ่นเครื่อง
		//"Product_serial":    product_serial,    //Serial Number ของเครื่องทั้งเครื่อง
		//"Product_uuid":      product_uuid,      //UUID ของเครื่อง
		"Product_sku": product_sku, //SKU/Part Number ของเครื่อง
		//mainboard
		"Board_vendor":  board_vendor,  //ผู้ผลิตเมนบอร์ด
		"Board_name":    board_name,    //รุ่นเมนบอร์ด
		"Board_version": board_version, //Revision/Version ของเมนบอร์ด
		//"Board_serial":      board_serial,      //Serial Number ของเมนบอร์ด
		"Board_asset_tag": board_asset_tag, //รหัสทรัพย์สิน (Asset Tag) ของเมนบอร์ด ใช้ในองค์กร
		//BIOS/UEFI
		"Bios_vendor":  bios_vendor,  //ผู้ผลิต BIOS
		"Bios_version": bios_version, //เวอร์ชัน BIOS
		"Bios_date":    bios_date,    //วันที่ออก BIOS
		"Bios_release": bios_release, //เวอร์ชัน Release ของ BIOS ตาม SMBIOS
		//Chassis
		"Chassis_vendor": chassis_vendor, //ผู้ผลิตตัวเครื่อง/เคส
		"Chassis_type":   chassis_type,   //ประเภทของเครื่อง
		//"Chassis_serial":    chassis_serial,    //Serial Number ของตัวเครื่อง/เคส
		"Chassis_version":   chassis_version,   //รุ่นหรือ Revision ของตัวเครื่อง
		"Chassis_asset_tag": chassis_asset_tag, //Asset Tag ของตัวเครื่อง
		"Modalias":          modalias,          //Hardware ID สำหรับ kernel ใช้จับคู่ driver
	}

}

func MainboardTabs() fyne.CanvasObject {
	x := mainboard_info()

	MainboardPkexec := widget.NewMultiLineEntry()
	MainboardPkexec.SetPlaceHolder("ผลลัพธ์การเข้าถึงเมนบอร์ด")
	MainboardPkexec.SetMinRowsVisible(8)
	//MainboardPkexec.SetReadOnly(true)

	MainboardPkexecAll := widget.NewMultiLineEntry()
	MainboardPkexecAll.SetPlaceHolder("ผลลัพธ์การขอสิทธิ์ทั้งหมด")
	MainboardPkexecAll.SetMinRowsVisible(12)
	//MainboardPkexecAll.SetReadOnly(true)

	RequestingAccessToMainboard := widget.NewButton("ขอสิทธิ์เข้าถึงเมนบอร์ด", func() {
		go func() {
			teXt, err := GetMainboardInfo()
			if err != nil {
				teXt = err.Error()
			}
			MainboardPkexec.SetText(teXt) //ให้มันอัพเดท
		}()
	})

	RequestingFullAccessToMainboard := widget.NewButton("ขอสิทธิ์ทั้งหมด", func() {
		go func() {
			teXt, err := GetInfoPkexecAll()
			if err != nil {
				if teXt == "" {
					teXt = err.Error()
				} else {
					teXt = fmt.Sprintf("%s\n%s", teXt, err.Error())
				}
			}
			MainboardPkexecAll.SetText(teXt)
		}()
	})

	subdetail_mainboard := container.NewVBox(
		RequestingAccessToMainboard,
		container.NewScroll(MainboardPkexec),
		RequestingFullAccessToMainboard,
		container.NewScroll(MainboardPkexecAll),
	)

	detail_mainboard := container.NewVBox(
		widget.NewCard("Detail", "", subdetail_mainboard),
	)

	subSystem := container.NewVBox(
		//System
		widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", x["Sys_vendor"])),       //ผู้ผลิตเครื่องทั้งเครื่อง (OEM)
		widget.NewLabel(fmt.Sprintf("รุ่น : %s", x["Product_name"])),        //รุ่นของเครื่อง
		widget.NewLabel(fmt.Sprintf("ตระกูล : %s", x["Product_family"])),    //ตระกูลของเครื่อง
		widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", x["Product_version"])), //เวอร์ชันหรือ Revision ของรุ่นเครื่อง
		//widget.NewLabel(fmt.Sprintf("Serial Number : %s", x[//"Product_serial"])),    //Serial Number ของเครื่องทั้งเครื่อง
		//widget.NewLabel(fmt.Sprintf("UUID : %s", x[//"Product_uuid"])),      //UUID ของเครื่อง
		widget.NewLabel(fmt.Sprintf("รหัส : %s", x["Product_sku"])), //SKU/Part Number ของเครื่อง
	)

	System := container.NewVBox(
		//System
		widget.NewCard("System", "", subSystem),
	)

	subMainboard := container.NewVBox(
		//mainboard
		widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", x["Board_vendor"])),   //ผู้ผลิตเมนบอร์ด
		widget.NewLabel(fmt.Sprintf("รุ่น : %s", x["Board_name"])),        //รุ่นเมนบอร์ด
		widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", x["Board_version"])), //Revision/Version ของเมนบอร์ด
		//widget.NewLabel(fmt.Sprintf("Serial Number : %s", x["Board_serial"])),     //Serial Number ของเมนบอร์ด
		widget.NewLabel(fmt.Sprintf("รหัส : %s", x["Board_asset_tag"])), //รหัสทรัพย์สิน (Asset Tag) ของเมนบอร์ด ใช้ในองค์กร

	)

	Mainboard := container.NewVBox(
		//mainboard
		widget.NewCard("Mainboard", "", subMainboard),
	)

	subBIOS_UEFI := container.NewVBox(
		//BIOS/UEFI
		widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", x["Bios_vendor"])),   //ผู้ผลิต BIOS
		widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", x["Bios_version"])), //เวอร์ชัน BIOS
		widget.NewLabel(fmt.Sprintf("วันที่ออก : %s", x["Bios_date"])),   //วันที่ออก BIOS
		widget.NewLabel(fmt.Sprintf("เวอร์ชัน : %s", x["Bios_release"])), //เวอร์ชัน Release ของ BIOS ตาม SMBIOS

	)
	BIOS_UEFI := container.NewVBox(
		//BIOS/UEFI
		widget.NewCard("BIOS/UEFI", "", subBIOS_UEFI),
	)

	subChassis := container.NewVBox(
		//Chassis
		widget.NewLabel(fmt.Sprintf("ผู้ผลิต : %s", x["Chassis_vendor"])), //ผู้ผลิตตัวเครื่อง/เคส
		widget.NewLabel(fmt.Sprintf("ประเภท : %s", x["Chassis_type"])),    //ประเภทของเครื่อง
		//widget.NewLabel(fmt.Sprintf("รหัส : %s", x["Chassis_serial"])),    //Serial Number ของตัวเครื่อง/เคส
		widget.NewLabel(fmt.Sprintf("รุ่น : %s", x["Chassis_version"])),  //รุ่นหรือ Revision ของตัวเครื่อง
		widget.NewLabel(fmt.Sprintf("Tag : %s", x["Chassis_asset_tag"])), //Asset Tag ของตัวเครื่อง
		//widget.NewLabel(fmt.Sprintf("Hardware ID : %s", x["Modalias"])),          //Hardware ID สำหรับ kernel ใช้จับคู่ driver
		//widget.NewSeparator(),

	)

	Chassis := container.NewVBox(
		//Chassis
		widget.NewCard("Chassis", "", subChassis),
	)

	return container.NewAppTabs(
		container.NewTabItem("System", container.NewScroll(System)),
		container.NewTabItem("Mainboard", container.NewScroll(Mainboard)),
		container.NewTabItem("BIOS / UEFI", container.NewScroll(BIOS_UEFI)),
		container.NewTabItem("Chassis", container.NewScroll(Chassis)),
		container.NewTabItem("Detail", container.NewScroll(detail_mainboard)),
	)
}
