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

func getCPUhardware(cpuIndex int) (fyne.CanvasObject, uint64, uint64, uint64) {
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

	var x1 strings.Builder
	var val_cpuinfo_min_freq uint64
	var val_cpuinfo_max_freq uint64
	var val_cpuinfo_transition_latency uint64

	x1.WriteString("Min - Max Hardware")
	x1.WriteString("\n|")

	for _, item := range files {
		data, err := os.ReadFile(base + item.file)
		if err != nil {
			x1.WriteString(fmt.Sprintf("\n%s: ไม่สามารถอ่านได้", item.label))
			continue
		}
		//string ไม่จำเป็นต้อง update()
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
		//เอาค่าออกมา
		if strings.Contains(item.file, "cpuinfo_min_freq") {
			val, _ := strconv.ParseFloat(value, 64)
			val_cpuinfo_min_freq = uint64(val)
		}
		if strings.Contains(item.file, "cpuinfo_max_freq") {
			val, _ := strconv.ParseFloat(value, 64)
			val_cpuinfo_max_freq = uint64(val)
		}
		if strings.Contains(item.file, "cpuinfo_transition_latency") {
			val, _ := strconv.ParseFloat(value, 64)
			val_cpuinfo_transition_latency = uint64(val)
		}
	}
	fyne.Do(func() {
		x.SetText(x1.String())
	})

	return x, val_cpuinfo_min_freq, val_cpuinfo_max_freq, val_cpuinfo_transition_latency
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
func slider() (fyne.CanvasObject, fyne.CanvasObject, fyne.CanvasObject, *widget.Label, *widget.Label, *widget.Label, uint64, uint64, uint64) {
	_, val_min, val_max, val_latency := getCPUhardware(0)
	val_ch_min := val_min
	val_ch_max := val_max
	val_ch_latency := val_latency

	min_freq_Label := widget.NewLabel(fmt.Sprintf("%d kHz [ %.2f Ghz ]", val_ch_min, float64(val_ch_min)/1e6))
	max_freq_Label := widget.NewLabel(fmt.Sprintf("%d kHz [ %.2f Ghz ]", val_ch_max, float64(val_ch_max)/1e6))
	latency_Label := widget.NewLabel(fmt.Sprintf("%d nS [ %.2f uS ]", val_ch_latency, float64(val_ch_latency)/1e3))

	min_freq_Slider := widget.NewSlider(float64(val_min), float64(val_max))
	min_freq_Slider.Step = 1
	min_freq_Slider.Value = float64(val_min) //ตั้งค่าเริ่มต้นของ slider
	min_freq_Slider.OnChanged = func(v float64) {
		val_ch_min = uint64(v) //แปลงเป็น uint64
		min_freq_Label.SetText(fmt.Sprintf("%d kHz [ %.2f Ghz ]", val_ch_min, float64(val_ch_min)/1e6))
	}

	max_freq_Slider := widget.NewSlider(float64(val_min), float64(val_max))
	max_freq_Slider.Step = 1
	max_freq_Slider.Value = float64(val_max)
	max_freq_Slider.OnChanged = func(v float64) {
		val_ch_max = uint64(v)
		max_freq_Label.SetText(fmt.Sprintf("%d kHz [ %.2f Ghz ]", val_ch_max, float64(val_ch_max)/1e6))
	}
	latency_Slider := widget.NewSlider(500, 50000)
	latency_Slider.Step = 1
	latency_Slider.Value = float64(val_latency)
	latency_Slider.OnChanged = func(v float64) {
		val_ch_latency = uint64(v)
		latency_Label.SetText(fmt.Sprintf("%d nS [ %.2f uS ]", val_ch_latency, float64(val_ch_latency)/1e3))
	}

	return min_freq_Slider, max_freq_Slider, latency_Slider, min_freq_Label, max_freq_Label, latency_Label, val_ch_min, val_ch_max, val_ch_latency
}

func onButtonClick() {

	_, _, _, _, _, _, freq_min, freq_max, latency := slider()
	//freq := uint64(2000000) // อ่านจาก input field

	go func() { // รันใน goroutine ไม่ให้ UI ค้าง
		//		script := fmt.Sprintf(
		//		"echo %d | tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_max_freq",
		//		freq,
		//	)

		script := fmt.Sprintf(
			`echo %d | tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_max_freq
echo %d | tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_min_freq
echo %d | tee /sys/devices/system/cpu/cpu*/cpufreq/cpuinfo_transition_latency`,
			freq_max,
			freq_min,
			latency,
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
	info, _, _, _ := getCPUhardware(0)
	slider_min, slider_max, slider_latency, label_min, label_max, label_latency, _, _, _ := slider()

	bt1 := widget.NewButton("TTT", func() {
		onButtonClick()
	})

	x := container.NewBorder(
		container.NewVBox(
			info,
			widget.NewSeparator(),
			bt1,
			widget.NewSeparator(),
			perCore,
			widget.NewSeparator(),
			slider_min,
			label_min,
			widget.NewSeparator(),
			slider_max,
			label_max,
			widget.NewSeparator(),
			slider_latency,
			label_latency,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		nil,
	)

	return x
}
