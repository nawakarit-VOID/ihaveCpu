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
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//-------------------------------------------------------------------------------------------

// getCPUFreqInfo อ่านข้อมูลความถี่ของ CPU
func getCPUFreqInfo(cpuIndex int) fyne.CanvasObject {

	base := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/", cpuIndex)
	files := []struct {
		file  string
		label string
	}{
		{"scaling_cur_freq", "ความถี่ปัจจุบัน"},
		{"scaling_max_freq", "ความถี่สูงสุด (เพดาน)"},
		{"scaling_min_freq", "ความถี่ต่ำสุด"},
		{"cpuinfo_max_freq", "ความถี่สูงสุดของ hardware"},
		{"scaling_governor", "Governor ที่ใช้อยู่"},
	}

	x := widget.NewLabel("กำลังโหลด...")

	update := func() {
		var x1 strings.Builder
		x1.WriteString("ยังไม่รองรับหลาย cpu")
		x1.WriteString(fmt.Sprintf("\ncore %d ", cpuIndex))

		for _, item := range files {
			data, err := os.ReadFile(base + item.file)
			if err != nil {
				x1.WriteString(fmt.Sprintf("\n%s: ไม่สามารถอ่านได้", item.label))
				continue
			}

			value := strings.TrimSpace(string(data))
			x1.WriteString(fmt.Sprintf("\n%s: %s", item.label, value))

			if strings.Contains(item.file, "freq") {
				val, _ := strconv.ParseFloat(value, 64)
				x1.WriteString(fmt.Sprintf(" kHz (%.2f GHz)", val/1e6))
			}
		}

		fyne.Do(func() {
			x.SetText(x1.String())
		})
	}

	update()
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			update()
		}
	}()

	return x
}

// ============================================================================
// แบ่งตามจำนวนคอร์
// ============================================================================
func sysCPUFreqInfo() fyne.CanvasObject {
	coreCount := CpuCoreCount()
	box := container.NewVBox()

	for i := 0; i < coreCount; i++ {
		coreInfo := getCPUFreqInfo(i)
		box.Add(coreInfo)
	}
	if coreCount == 0 {
		return widget.NewLabel("ไม่พบข้อมูลจำนวนคอร์ CPU")
	}
	return box
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
	//xu0 := widget.NewLabel("ยังไม่รองรับหลาย cpu")
	//xu1 := getCPUFreqInfo(0) //เลือก คอร์ 0

	xu2 := sysCPUFreqInfo()

	bt1 := widget.NewButton("TTT", func() {
		onButtonClick()
	})

	x := container.NewBorder(
		container.NewVBox(bt1, xu2),
		nil,
		nil,
		nil,
	)

	return x
}
