// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Package_ui

import (
	"embed"
	cpuinfo "ihavecpu/cpu"
	mainboardinfo "ihavecpu/mainboard"
	raminfo "ihavecpu/ram"
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

func GetDataIn() (string, string, error) {
	// เปลี่ยนจาก "sudo" เป็น "pkexec"
	cmd := exec.Command("pkexec", "sh", "-c",
		"dmidecode -t memory && echo '__SPLIT__' && dmidecode -t 2")

	out, err := cmd.Output()
	if err != nil {
		return "", "", err
	}

	parts := strings.SplitN(string(out), "__SPLIT__", 2)

	mem := strings.TrimSpace(parts[0])

	var board string
	if len(parts) > 1 {
		board = strings.TrimSpace(parts[1])
	}

	return mem, board, nil
}

func CreateWindow() {

	a := app.NewWithID("com.nawakarit.iHaveCPU")
	a.Settings().SetTheme(&MyTheme{})
	icon := loadIcon(64)
	w := a.NewWindow("iHaveCPU")
	w.SetIcon(icon)

	cpuTabs := cpuinfo.CpuTabs(w)
	mainboardTabs := mainboardinfo.MainboardTabs()
	ram := raminfo.RamTabs()

	memInfo, boardInfo, err := GetDataIn()

	if err != nil {
		return
	}

	fyne.Do(func() {

		//raminfo.TestDetailLabelcmd(testAll)
		raminfo.RamDetailLabelcmd(memInfo)
		mainboardinfo.TestDetailLabelcmd(boardInfo)
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
		container.NewTabItem("CPU", container.NewScroll(cpuTabs)),
		container.NewTabItem("MainBoard", container.NewScroll(mainboardTabs)),
		container.NewTabItem("Ram", ram),
		//container.NewTabItem("Security", container.NewScroll(nil)),
		//container.NewTabItem("Virtualization", container.NewScroll(nil)),
	)

	//w.SetContent(container.NewBorder(nil, nil, nil, nil, cpu))
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(720, 800))
	w.ShowAndRun()
}
