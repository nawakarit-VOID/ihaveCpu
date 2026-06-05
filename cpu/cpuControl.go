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

// -------------------------------------------------------------------------------------------
type xx struct {
	Usage string //
	//Timesusage string

	//////////////////////
}

func getCPUhardware(cpuIndex int) fyne.CanvasObject {
	base := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/", cpuIndex)
	files := []struct {
		file  string
		label string
	}{
		{"cpuinfo_min_freq", "ความถี่ต่ำสุด"},
		{"cpuinfo_max_freq", "ความถี่สูงสุด"},
		{"cpuinfo_transition_latency", "เวลาในการเปลี่ยนความเร็ว"},
		{"scaling_governor", "Governor ที่ใช้อยู่"},
	}

	x := widget.NewLabel("กำลังโหลด...")

	update := func() {
		var x1 strings.Builder
		x1.WriteString("Min - Max Hardware")
		x1.WriteString("\n|")

		for _, item := range files {
			data, err := os.ReadFile(base + item.file)
			if err != nil {
				x1.WriteString(fmt.Sprintf("\n%s: ไม่สามารถอ่านได้", item.label))
				continue
			}

			value := strings.TrimSpace(string(data))
			x1.WriteString(fmt.Sprintf("\n%s : %s", item.label, value))

			if strings.Contains(item.file, "freq") {
				val, _ := strconv.ParseFloat(value, 64)
				x1.WriteString(fmt.Sprintf(" kHz // (%.2f GHz)", val/1e6))
			}
			if strings.Contains(item.file, "latency") {
				val, _ := strconv.ParseFloat(value, 64)
				x1.WriteString(fmt.Sprintf(" nS // (%.f uS)", val/1e3))
			}
		}
		fyne.Do(func() {
			x.SetText(x1.String())
		})
	}
	update()
	return x
}

// getCPUFreqInfo อ่านข้อมูลความถี่ของ CPU
func getCPUFreqUpdate(cpuIndex int) fyne.CanvasObject {

	base := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/", cpuIndex)
	files := []struct {
		file  string
		label string
	}{
		{"scaling_cur_freq", "ความถี่ปัจจุบัน"},
		{"scaling_max_freq", "ความถี่สูงสุด (จำกัด)"},
		{"scaling_min_freq", "ความถี่ต่ำสุด (จำกัด)"},
	}

	x := widget.NewLabel("กำลังโหลด...")

	update := func() {
		var x1 strings.Builder
		//x1.WriteString("ยังไม่รองรับหลาย cpu")
		x1.WriteString(fmt.Sprintf("core [ %d ]", cpuIndex))

		for _, item := range files {
			data, err := os.ReadFile(base + item.file)
			if err != nil {
				x1.WriteString(fmt.Sprintf("\n%s: ไม่สามารถอ่านได้", item.label))
				continue
			}

			value := strings.TrimSpace(string(data))
			x1.WriteString(fmt.Sprintf("\n%s : %s", item.label, value))

			if strings.Contains(item.file, "freq") {
				val, _ := strconv.ParseFloat(value, 64)
				x1.WriteString(fmt.Sprintf(" kHz // (%.2f GHz)", val/1e6))
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
func sysCPUFreqUpdate() fyne.CanvasObject {
	coreCount := CpuCoreCount()
	box := container.NewVBox()

	for i := 0; i < coreCount; i++ {
		coreInfo := getCPUFreqUpdate(i)
		box.Add(coreInfo)
	}
	if coreCount == 0 {
		return widget.NewLabel("ไม่พบข้อมูลจำนวนคอร์ CPU")
	}
	return box
}

func cpuSlider(cpuIndex int) uint64 {
	base := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/", cpuIndex)
	files := []struct {
		file  string
		label string
	}{
		{"cpuinfo_min_freq", "ความถี่ต่ำสุด"},
		{"cpuinfo_max_freq", "ความถี่สูงสุด"},
		{"cpuinfo_transition_latency", "เวลาในการเปลี่ยนความเร็ว"},
	}

	x := widget.NewLabel("กำลังโหลด...")

	update := func() {
		var x1 strings.Builder
		x1.WriteString("Min - Max Hardware")
		x1.WriteString("\n|")

		for _, item := range files {
			data, err := os.ReadFile(base + item.file)
			if err != nil {
				x1.WriteString(fmt.Sprintf("\n%s: ไม่สามารถอ่านได้", item.label))
				continue
			}

			value := strings.TrimSpace(string(data))
			x1.WriteString(fmt.Sprintf("\n%s : %s", item.label, value))

			if strings.Contains(item.file, "freq") {
				val, _ := strconv.ParseFloat(value, 64)
				x1.WriteString(fmt.Sprintf(" kHz // (%.2f GHz)", val/1e6))
			}
			if strings.Contains(item.file, "latency") {
				val, _ := strconv.ParseFloat(value, 64)
				x1.WriteString(fmt.Sprintf(" nS // (%.f uS)", val/1e3))
			}
		}
		fyne.Do(func() {
			x.SetText(x1.String())
		})
	}
	update()
	return x

	/*
		cpuSlider := widget.NewSlider(1, float64(maxCPU)) // สร้าง slider สำหรับเลือกจำนวน CPU ที่จะใช้ โดยมีค่าตั้งแต่ 1 ถึงจำนวน CPU สูงสุดของเครื่อง
		cpuSlider.Step = 1                                //ใช้เฉพาะจำนวนเต็ม เพราะ workers และ parallelism ต้องเป็นจำนวนเต็ม
		cpuSlider.Value = float64(maxCPU)                 //ตั้งค่าเริ่มต้นของ slider ให้เป็นจำนวน CPU สูงสุด (ใช้ทุก core ที่มี)
		cpuSlider.OnChanged = func(v float64) {           //เมื่อ slider ถูกเปลี่ยนค่า จะคำนวณเปอร์เซ็นต์การใช้ CPU ใหม่และอัปเดตข้อความใน cpuLabel ตามค่าที่เลือก
			pvcpus := pmcpu * v
			symbol := SpeedSymbol(pvcpus) //แสดงสัญลักษณ์ความเร็วตามเปอร์เซ็นต์การใช้ CPU เริ่มต้น
			cpuLabel.Text = fmt.Sprintf("CPU Speed x%.1f %s ( %.0f%% / cores ) %s", v, symbol, pvcpus, symbol)
			cpuLabel.Refresh()
		}


			if strings.Contains(item.file, "cpuinfo_max_freq") {
				val_cpuinfo_max_freq, _ := strconv.ParseFloat(value, 64)
			}




	return */

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

	perCore := sysCPUFreqUpdate()
	info := getCPUhardware(0)

	bt1 := widget.NewButton("TTT", func() {
		onButtonClick()
	})

	x := container.NewBorder(
		container.NewVBox(
			info,
			widget.NewSeparator(),
			bt1,
			widget.NewSeparator(),
			perCore),
		nil,
		nil,
		nil,
		nil,
	)

	return x
}
