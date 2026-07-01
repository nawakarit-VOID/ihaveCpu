// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Package_ui

import (
	"embed"
	biosinfo "ihavecpu/bios"
	cpuinfo "ihavecpu/cpu"
	mainboardinfo "ihavecpu/mainboard"
	pcieinfo "ihavecpu/pcie"
	raminfo "ihavecpu/ram"
	systeminfo "ihavecpu/system"

	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// โหลด icon
func loadIcon(size int) fyne.Resource {
	var file string

	switch {
	case size >= 512:
		file = "assets/icons/icon-512.png" ///ที่อยู่
	case size >= 256:
		file = "assets/icons/icon-256.png"
	case size >= 128:
		file = "assets/icons/icon-128.png"
	default:
		file = "assets/icons/icon-64.png"
	}

	data, _ := iconFS.ReadFile(file)
	return fyne.NewStaticResource(file, data)
}

//go:embed assets/icons/*
var iconFS embed.FS

//go:embed assets/font/Itim-Regular.ttf
var fontItim []byte
var myFont = fyne.NewStaticResource("Itim-Regular.ttf", fontItim)

func GetDataIn() (string, string, string, string, string, string, error) {
	// เปลี่ยนจาก "sudo" เป็น "pkexec"
	cmd := exec.Command("pkexec", "sh", "-c",
		`dmidecode -t 1 && dmidecode -t 3 && dmidecode -t 12 && dmidecode -t  15 && dmidecode -t 23 && dmidecode -t 24 && dmidecode -t 32
echo '(-@_@-)' && dmidecode -t 0 && dmidecode -t 13 && dmidecode -t 40  && dmidecode -t 45
echo '(-@_@-)' && dmidecode -t 4
echo '(-@_@-)' && dmidecode -t 7 
echo '(-@_@-)' && dmidecode -t 5 && dmidecode -t 6 && dmidecode -t 16 && dmidecode -t 17 && dmidecode -t 18 && dmidecode -t 19 && dmidecode -t 20 && dmidecode -t 33 && dmidecode -t 37
echo '(-@_@-)' && dmidecode -t 2 && dmidecode -t 10 && dmidecode -t 41
`)
	//dmidecode -t memory

	out, err := cmd.Output()
	if err != nil {
		return "", "", "", "", "", "", err
	}

	parts := strings.Split(string(out), "(-@_@-)")

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	//ย้าย bios มาไว้อันแรก แล้วก็เรียกมาแสดงทุก type
	sys := parts[0]
	bios := parts[1]
	cpu := parts[2]
	cache := parts[3]
	ram := parts[4]
	board := parts[5]

	return sys, bios, cpu, cache, ram, board, nil
}

func CreateWindow() {

	a := app.NewWithID("com.nawakarit.iHaveCPU")
	a.Settings().SetTheme(&MyTheme{})
	icon := loadIcon(64)
	w := a.NewWindow("iHaveCPU")
	w.SetIcon(icon)

	system := systeminfo.SystemTabs()
	biOsTabs := biosinfo.BiosTabs()
	cpuTabs := cpuinfo.CpuTabs(w)
	mainboardTabs := mainboardinfo.MainboardTabs()
	ramTabs := raminfo.RamTabs()
	pcieTabs := pcieinfo.PcieTabs()

	sys, bios, cpu, chsche, ram, boardInfo,
		err := GetDataIn()

	if err != nil {
		return
	}

	fyne.Do(func() {
		systeminfo.SystemsDetailLabelcmd(sys)
		biosinfo.BiosDetailLabelcmd(bios)
		cpuinfo.CPUDetailLabelcmd(cpu)
		cpuinfo.CacheLabelcmd(chsche)
		raminfo.RamDetailLabelcmd(ram)
		mainboardinfo.MainboardDetailLabelcmd(boardInfo)

	})

	/*
		teXt, err := raminfo.GetMemoryInfo()

		if err != nil {
			teXt = err.Error()
		}
		fyne.Do(func() {
			raminfo.RamDetailLabelcmd(teXt)
		})
	*/

	//MemoryPkexec.SetText(teXt) //ให้มันอัพเดท
	/*
		cmd := exec.Command("pkexec", "bash", "-t", script)
		output, err := cmd.CombinedOutput()
		text := string(output)
		if err != nil {
			text = fmt.Sprintf("%s\n%s", text, err.Error())
			fmt.Printf("failed to run pkexec: %v\n%s\n", err, string(output))
		}
		mainboardinfo.SetMainboardPkexecAllText(text)
	*/
	tabs := container.NewAppTabs(
		container.NewTabItem("System", system),
		container.NewTabItem("Bios", biOsTabs),
		container.NewTabItem("CPU", container.NewScroll(cpuTabs)),
		container.NewTabItem("Ram", ramTabs),
		container.NewTabItem("MainBoard", container.NewScroll(mainboardTabs)),
		container.NewTabItem("Pcie", container.NewScroll(pcieTabs)),
		//container.NewTabItem("Security", container.NewScroll(nil)),
		//container.NewTabItem("Virtualization", container.NewScroll(nil)),
	)

	//w.SetContent(container.NewBorder(nil, nil, nil, nil, cpu))
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(720, 800))
	w.ShowAndRun()
}
