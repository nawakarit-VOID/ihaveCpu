// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_cpuinfo

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//-------------------------------------------------------------------------------------------

// getCPUFreqInfo อ่านข้อมูลความถี่ของ CPU
func getCPUFreqInfo(cpuIndex int) fyne.CanvasObject {
	base := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/", cpuIndex)
	files := map[string]string{
		"scaling_cur_freq": "ความถี่ปัจจุบัน",
		"scaling_max_freq": "ความถี่สูงสุด (เพดาน)",
		"scaling_min_freq": "ความถี่ต่ำสุด",
		"cpuinfo_max_freq": "ความถี่สูงสุดของ hardware",
		"scaling_governor": "governor ที่ใช้อยู่",
	}

	x := widget.NewLabel("x...")
	x1 := ""
	for file, label := range files {
		data, err := os.ReadFile(base + file)
		if err != nil {
			fmt.Printf("  %s: ไม่สามารถอ่านได้\n", label)
			continue
		}

		//fmt.Printf("  %s: %s", label, strings.TrimSpace(string(data)))
		x1 += fmt.Sprintf("\n%s: %s", label, strings.TrimSpace(string(data)))

		if strings.Contains(file, "freq") {
			val, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			x1 += fmt.Sprintf(" kHz (%.2f GHz)", val/1e6)
		}
		x.SetText(x1)
	}
	return x
}

func onButtonClick() {
	freq := uint64(2000000) // อ่านจาก input field

	go func() { // รันใน goroutine ไม่ให้ UI ค้าง
		script := fmt.Sprintf(
			"echo %d | tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_max_freq",
			freq,
		)
		cmd := exec.Command("pkexec", "bash", "-c", script)
		err := cmd.Run()
		if err != nil {
			// แสดง error dialog
			fmt.Println("ล้มเหลว")
		}
		// แสดง success dialog
		fmt.Println("สำเร็จ 2GHz")
	}()
}

// ส่งออก
func CpuControl() fyne.CanvasObject {

	//bar := widget.NewProgressBar()
	//label := widget.NewLabel("0%")

	/*
		go func() {
			for {
				//v := CpuPercentAVG()
				//val := v[0] / 100.0

				fyne.Do(func() { //กันพัง'
					 อัปเดต UI
					bar.SetValue(val)
					label.SetText(fmt.Sprintf("%.0f%%", val*100))

				})

				time.Sleep(500 * time.Millisecond)
			}
		}()
	*/
	xu1 := getCPUFreqInfo(0)

	bt1 := widget.NewButton("TTT", func() {
		onButtonClick()
	})

	x := container.NewBorder(
		container.NewVBox(xu1),
		nil,
		nil,
		nil,
		bt1)

	return x
}
