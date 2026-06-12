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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// -------------------------------------------------------------------------------------------

func getCPUhardware(cpuIndex int) (fyne.CanvasObject, uint64, uint64) {
	base := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/", cpuIndex)
	files := []struct {
		file  string
		label string
	}{
		{"cpuinfo_min_freq", "ความถี่ต่ำสุด"},
		{"cpuinfo_max_freq", "ความถี่สูงสุด"},
	}
	x := widget.NewLabel("กำลังโหลด...")
	x.Alignment = fyne.TextAlignCenter //ทำให้ตรงกลาง

	var x1 strings.Builder
	var val_cpuinfo_min_freq uint64
	var val_cpuinfo_max_freq uint64

	x1.WriteString("Default Kernel and Hardware")
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

		//เอาค่าออกมา
		if strings.Contains(item.file, "cpuinfo_min_freq") {
			val, _ := strconv.ParseFloat(value, 64)
			val_cpuinfo_min_freq = uint64(val)
		}
		if strings.Contains(item.file, "cpuinfo_max_freq") {
			val, _ := strconv.ParseFloat(value, 64)
			val_cpuinfo_max_freq = uint64(val)
		}

	}
	fyne.Do(func() {
		x.SetText(x1.String())
	})

	return x, val_cpuinfo_min_freq, val_cpuinfo_max_freq
}

// getCPUFreqInfo อ่านข้อมูลความถี่ของ CPU
func getCPUFreqUpdate(cpuIndex int) (fyne.CanvasObject, uint64, uint64) {

	base := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/", cpuIndex)
	files := []struct {
		file  string
		label string
	}{
		{"scaling_cur_freq", "ความถี่ปัจจุบัน"},
		{"scaling_min_freq", "ความถี่ต่ำสุด (จำกัด)"},
		{"scaling_max_freq", "ความถี่สูงสุด (จำกัด)"},
		{"cpuinfo_transition_latency", "เวลาในการเปลี่ยนความเร็ว"},
		{"scaling_governor", "Governor ที่ใช้อยู่"},
	}

	var val_scaling_min_freq uint64
	var val_scaling_mmax_freq uint64

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
			if strings.Contains(item.file, "latency") {
				val, _ := strconv.ParseFloat(value, 64)
				x1.WriteString(fmt.Sprintf(" nS // (%.f uS)", val/1e3))
			}
			//เอาค่าปัจจุบันออกมา
			if strings.Contains(item.file, "scaling_min_freq") {
				val, _ := strconv.ParseFloat(value, 64)
				val_scaling_min_freq = uint64(val)
			}
			if strings.Contains(item.file, "scaling_max_freq") {
				val, _ := strconv.ParseFloat(value, 64)
				val_scaling_mmax_freq = uint64(val)
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

	return x, val_scaling_min_freq, val_scaling_mmax_freq
}

// ============================================================================
// เพิ่ม label ตามจำนวนคอร์
// ============================================================================
func sysCPUFreqUpdate() fyne.CanvasObject {
	coreCount := CpuCoreCount()
	box := container.NewVBox()

	for i := 0; i < coreCount; i++ {
		coreInfo, _, _ := getCPUFreqUpdate(i)
		box.Add(coreInfo)
	}
	if coreCount == 0 {
		return widget.NewLabel("ไม่พบข้อมูลจำนวนคอร์ CPU")
	}
	return box
}

// ============================================================================
// เลือกทั้งหมด
// ============================================================================
func CheckAllBoxCpu(checkboxes []*widget.Check, selected []bool, updateLabel func()) {
	for idx, check := range checkboxes {
		check.SetChecked(true)
		selected[idx] = true
	}
	updateLabel()
	fmt.Println("เลือกทั้งหมด")
}

// ============================================================================
// ไม่เลือกทั้งหมด
// ============================================================================
func nonCheckBoxCpu(checkboxes []*widget.Check, selected []bool, updateLabel func()) {
	for idx, check := range checkboxes {
		check.SetChecked(false)
		selected[idx] = false
	}
	updateLabel()
	fmt.Println("ล้างทั้งหมด")
}

// ============================================================================
// เพิ่ม checkbox ตามจำนวนคอร์
// ============================================================================
func checkboxNumcpu() (fyne.CanvasObject, []bool, []*widget.Check, func()) {
	coreCount := CpuCoreCount()
	if coreCount == 0 {
		return widget.NewLabel("ไม่พบข้อมูลจำนวนคอร์ CPU"), nil, nil, nil
	}

	selected := make([]bool, coreCount)
	checkboxes := make([]*widget.Check, coreCount)
	for i := 0; i < coreCount; i++ {
		selected[i] = true
	}

	selectedGet, _ := getSelectedCoresText(selected)
	//fmt.Println("core เริ่ม", selectedGet)
	selectedLabel := widget.NewLabel(selectedGet)

	box := container.NewGridWithColumns(8)
	for i := 0; i < coreCount; i++ {
		idx := i
		coreName := strconv.Itoa(idx)
		x := widget.NewCheck("core "+coreName, func(checked bool) {
			selected[idx] = checked
			if checked {
				fmt.Println("core", idx, "เปิด")
			} else {
				fmt.Println("core", idx, "ปิด")
			}
			selectedGet, _ := getSelectedCoresText(selected)
			//fmt.Println("core ใน for", selectedGet)

			selectedLabel.SetText(selectedGet)
		})

		x.SetChecked(true)
		checkboxes[idx] = x
		box.Add(x)
	}

	// สร้าง function เพื่ออัปเดต label
	updateLabel := func() {
		selectedGet, _ := getSelectedCoresText(selected)
		selectedLabel.SetText(selectedGet)
	}

	return container.NewVBox(selectedLabel, box), selected, checkboxes, updateLabel
}

func getSelectedCoresText(selected []bool) (string, []int) {
	var cores []string
	var coresIndices []int
	for idx, checked := range selected {
		if checked {
			cores = append(cores, strconv.Itoa(idx))
			coresIndices = append(coresIndices, idx)
		}
	}

	if len(cores) == 0 {
		return "คอร์ที่เลือก : ไม่มี", nil
	}

	var lines []string

	for i := 0; i < len(cores); i += 40 { //40
		end := i + 40 //40
		if end > len(cores) {
			end = len(cores)
		}

		lines = append(lines,
			strings.Join(cores[i:end], ", "),
		)
	}
	//return "คอร์ที่เลือก : " + strings.Join(lines, "\n"), coresIndices
	return "คอร์ที่เลือก : " + lines[0] + "\n                     " + strings.Join(lines[1:], "\n                     "), coresIndices
} //21

// ============================================================================
// เรียกไฟล์ govenors
// ============================================================================
func GetGovernors() ([]string, error) {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_available_governors")
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(data)), nil
}

// สร้าง check Govern...
func GovernorscheckBox() (fyne.CanvasObject, *widget.RadioGroup) {

	governors, _ := GetGovernors()
	governorsST := widget.NewRadioGroup(governors, nil)
	governorsST.Horizontal = true

	if len(governors) > 0 {
		governorsST.SetSelected(governors[0])
	}

	return governorsST, governorsST
}

func onButtonMinN(min_freq_Slider *widget.Slider) { //ลดค่า min
	freq_min := min_freq_Slider.Value - 1
	if freq_min >= min_freq_Slider.Min {
		min_freq_Slider.SetValue(freq_min)
	}
}

func onButtonMinP(min_freq_Slider *widget.Slider) { //เพิ่มค่า min
	freq_min := min_freq_Slider.Value + 1
	if freq_min <= min_freq_Slider.Max {
		min_freq_Slider.SetValue(freq_min)
	}
}

func onButtonMaxN(max_freq_Slider *widget.Slider) { //ลดค่า max
	freq_max := max_freq_Slider.Value - 1
	if freq_max >= max_freq_Slider.Min {
		max_freq_Slider.SetValue(freq_max)
	}
}

func onButtonMaxP(max_freq_Slider *widget.Slider) { //เพิ่มค่า max
	freq_max := max_freq_Slider.Value + 1
	if freq_max <= max_freq_Slider.Max {
		max_freq_Slider.SetValue(freq_max)
	}
}

func slider() (*widget.Slider, *widget.Slider, *widget.Label, *widget.Label, *widget.Entry, *widget.Entry) {

	_, val_min, val_max := getCPUhardware(0)
	_, cur_min, cur_max := getCPUFreqUpdate(0)

	entry_min := widget.NewEntry()
	entry_min.SetText(strconv.FormatUint(cur_min, 10)) //10 คือแปลงเป็นเลขฐาน 10
	//entry_min.SetText(fmt.Sprintf("%d", cur_min/1000)) //หรือ

	entry_max := widget.NewEntry()
	entry_max.SetText(strconv.FormatUint(cur_max, 10)) //10 คือแปลงเป็นเลขฐาน 10
	//entry_max.SetText(fmt.Sprintf("%d", cur_max/1000)) //หรือ

	val_ch_min := val_min
	val_ch_max := val_max

	//label slider min และ max
	min_freq_Label := widget.NewLabel(fmt.Sprintf("[ จำกัด - ความถี่ต่ำสุด ] %d kHz [ %.2f Ghz ]", val_ch_min, float64(val_ch_min)/1e6))
	max_freq_Label := widget.NewLabel(fmt.Sprintf("[ จำกัด - ความถี่สูงสุด ]%d kHz [ %.2f Ghz ]", val_ch_max, float64(val_ch_max)/1e6))

	//ค่า slider ต่ำสุด-มากสุด ของ min และ max
	min_freq_Slider := widget.NewSlider(float64(val_min), float64(val_max)) //*min
	max_freq_Slider := widget.NewSlider(float64(val_min), float64(val_max)) //*max

	//*min
	entry_min.OnChanged = func(s string) {
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			min_freq_Slider.SetValue(v)
		}
	}
	//slider_min
	min_freq_Slider.Step = 1
	min_freq_Slider.Value = float64(cur_min) //ตั้งค่าเริ่มต้นของ slider
	min_freq_Slider.OnChanged = func(v float64) {
		//ตรวจสอบ min มากกว่า max ให้ เปลี่ยนค่า max ตาม min
		if v > max_freq_Slider.Value {
			max_freq_Slider.SetValue(v)
			//fmt.Println("max <= min")
		}

		entry_min.SetText(fmt.Sprintf("%.f", v))
		val_ch_min = uint64(v) //แปลงเป็น uint64
		min_freq_Label.SetText(fmt.Sprintf("[ จำกัด - ความถี่ต่ำสุด ] %d kHz [ %.2f Ghz ]", val_ch_min, float64(val_ch_min)/1e6))
	}

	//*max
	entry_max.OnChanged = func(s string) {
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			max_freq_Slider.SetValue(v)
		}
	}
	//slider_max
	max_freq_Slider.Step = 1
	max_freq_Slider.Value = float64(cur_max)
	max_freq_Slider.OnChanged = func(v float64) {
		//ตรวจสอบ max มากกว่า min ให้ เปลี่ยนค่า min ตาม max
		if v < min_freq_Slider.Value {
			min_freq_Slider.SetValue(v)
			//fmt.Println("min >= max")
		}
		val_ch_max = uint64(v)
		entry_max.SetText(fmt.Sprintf("%.f", v))
		max_freq_Label.SetText(fmt.Sprintf("[ จำกัด - ความถี่สูงสุด ] %d kHz [ %.2f Ghz ]", val_ch_max, float64(val_ch_max)/1e6))
	}

	return min_freq_Slider, max_freq_Slider, min_freq_Label, max_freq_Label, entry_min, entry_max
}

func onButtonClickApply(selected []bool, min_freq_Slider, max_freq_Slider *widget.Slider, governorsST *widget.RadioGroup) {

	// อ่านค่าจากวิดเจต slider โดยตรง
	freq_min := uint64(min_freq_Slider.Value)
	freq_max := uint64(max_freq_Slider.Value)
	governorsSt := governorsST.Selected

	go func() { // รันใน goroutine ไม่ให้ UI ค้าง
		var scriptLines []string
		for idx, sel := range selected {
			if !sel {
				continue
			}
			scriptLines = append(scriptLines, fmt.Sprintf("echo %d | tee /sys/devices/system/cpu/cpu%d/cpufreq/scaling_max_freq", freq_max, idx))
			scriptLines = append(scriptLines, fmt.Sprintf("echo %d | tee /sys/devices/system/cpu/cpu%d/cpufreq/scaling_min_freq", freq_min, idx))
			scriptLines = append(scriptLines, fmt.Sprintf("echo %s | tee /sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor", governorsSt, idx))
		}

		if len(scriptLines) == 0 {
			//ฟังชั้น popup++
			fmt.Println("ไม่พบคอร์ที่เลือกให้ปรับค่า")
			return
		}

		script := strings.Join(scriptLines, "\n")

		cmd := exec.Command("pkexec", "bash", "-c", script)
		err := cmd.Run()
		if err != nil {
			fmt.Println("ล้มเหลว:", err)
			return
		}
		fmt.Println("สำเร็จ", "[ min ]", freq_min, "kHz", "[ max ]", freq_max, "kHz")
	}()

}

func checkCoreCpu() {
	//coreCount := CpuCoreCount()
	//widget.NewCheck()
}

// ส่งออก
func CpuControl(w fyne.Window) fyne.CanvasObject {

	perCore := sysCPUFreqUpdate()
	info, _, _ := getCPUhardware(0)
	slider_min, slider_max, label_min, label_max, entry_min, entry_max := slider()

	chekCpu, selected, checkboxes, updateLabel := checkboxNumcpu()

	allCheck := widget.NewButton("เลือกทั้งหมด", func() {
		CheckAllBoxCpu(checkboxes, selected, updateLabel)
	})

	nonCheck := widget.NewButton("Reset", func() {
		nonCheckBoxCpu(checkboxes, selected, updateLabel)
	})

	//check govern...
	governors, governorsSt := GovernorscheckBox()

	apply := widget.NewButton("Apply", func() {
		onButtonClickApply(selected, slider_min, slider_max, governorsSt)
	})

	bt_min_n := widget.NewButton("-", func() {
		onButtonMinN(slider_min)
	})

	bt_min_p := widget.NewButton("+", func() {
		onButtonMinP(slider_min)
	})

	bt_max_n := widget.NewButton("-", func() {
		onButtonMaxN(slider_max)
	})

	bt_max_p := widget.NewButton("+", func() {
		onButtonMaxP(slider_max)
	})

	abbtn := widget.NewButton("โหมดการทำงาน", func() {
		dialog.ShowInformation("about_test", "test test", w)
	})

	//dialog.ShowInformation("warn", "Please add the folder first.", w)

	x := container.NewBorder(
		container.NewVBox(
			info,
			chekCpu,
			container.NewGridWithColumns(2, allCheck, nonCheck),
			//popup widget.NewLabel("โหมดการทำงาน\n#conservative- เพิ่มความเร็วแบบค่อยเป็นค่อยไป #####"),
			abbtn,
			governors,
			widget.NewSeparator(),
			container.NewHBox(label_min,
				container.NewGridWrap(fyne.NewSize(100, 35), entry_min),
				container.NewGridWrap(fyne.NewSize(35, 35), bt_min_n),
				container.NewGridWrap(fyne.NewSize(35, 35), bt_min_p)),
			slider_min,

			widget.NewSeparator(),
			container.NewHBox(label_max,
				container.NewGridWrap(fyne.NewSize(100, 35), entry_max),
				container.NewGridWrap(fyne.NewSize(35, 35), bt_max_n),
				container.NewGridWrap(fyne.NewSize(35, 35), bt_max_p)),
			slider_max,
			widget.NewSeparator(),
			apply,
			widget.NewSeparator(),
			perCore,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		nil,
	)

	return x
}
