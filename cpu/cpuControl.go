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

		x1 += fmt.Sprintf("  %s: %s\n", label, strings.TrimSpace(string(data)))

		if strings.Contains(file, "freq") {

			val, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			x1 += fmt.Printf(" kHz (%.2f GHz)", val/1e6)
		}
		x.SetText(x1)
		//y := fmt.Sprintf(x)

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

	bar := widget.NewProgressBar()
	label := widget.NewLabel("0%")

	go func() {
		for {
			v := CpuPercentAVG()
			val := v[0] / 100.0

			fyne.Do(func() { //กันพัง'
				// อัปเดต UI
				bar.SetValue(val)
				label.SetText(fmt.Sprintf("%.0f%%", val*100))

			})

			time.Sleep(500 * time.Millisecond)
		}
	}()

	//ProgressCpu0 := widget.NewProgressBar()
	fmt.Println("=== ข้อมูล CPU0 ===")

	//globalProgress.SetValue(float64(fi) / float64(totalFolders))
	xu1 := getCPUFreqInfo(0)

	/*
		// ตัวอย่าง: ตั้งเพดานที่ 2.0 GHz = 2,000,000 kHz
		targetFreq := uint64(2_000_000)
		fmt.Printf("\nตั้งเพดานความถี่ CPU0 เป็น %.1f GHz...\n", float64(targetFreq)/1e6)

		if err := setCPUMaxFreq(0, targetFreq); err != nil {
			fmt.Printf("เกิดข้อผิดพลาด: %v (ต้องรันด้วย root)\n", err)
			return
		}

		// governor ที่ใช้บ่อย: "powersave", "performance", "schedutil", "ondemand"
		setGovernor(0, "powersave")
		fmt.Println("สำเร็จ!")
	*/
	//x2 := widget.NewLabel(xu1)

	bt1 := widget.NewButton("TTT", func() {
		onButtonClick()
	})

	x := container.NewBorder(
		container.NewVBox(bar, label, xu1),
		nil,
		nil,
		nil,
		bt1)

	return x
}
