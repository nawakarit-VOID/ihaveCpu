// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package Ppackage_cpuinfo

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/v3/cpu"
)

/*
	func getCPU() []float64 {
		v, err := cpu.Percent(0, true)
		if err != nil {
			log.Println(err)
			return nil
		}
		return v
	}
*/
//จำนวนคอร์
func CpuCoreCount() int {
	physical, err := cpu.Counts(false) //core
	if err != nil {
		log.Println(err)
		return (0)
	}
	return physical
}

// จำนวนเทรด
func CpuThreadCount() int {
	logical, err := cpu.Counts(true) //thread
	if err != nil {
		log.Println(err)
		return (0)
	}
	return logical
}

func CpuPercentAVG() []float64 { //*[]float64
	// ดึง CPU usage รวม
	percentAVG, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil || len(percentAVG) == 0 {
		log.Println(err)
		return []float64{0.0}

	}

	return percentAVG
}

func CpuPercentPercore() []float64 {
	// ดึง CPU usage ต่อ core
	percentPerCore, err := cpu.Percent(100*time.Millisecond, true)
	if err != nil {
		//return []float64{0.0}
		return nil

	}
	return percentPerCore
}

//////////////////////////////////////////////////
// 📊 CPU
//////////////////////////////////////////////////
/*
func getCPU() []float64 {
	v, err := cpu.Percent(0, true)
	if err != nil {
		log.Println(err)
		return nil
	}
	return v
}
*/
func CPUdata() map[string]interface{} {
	// gopsutil
	info, _ := cpu.Info()
	//physical, _ := cpu.Counts(false) //core
	physical := CpuCoreCount()  //core
	logical := CpuThreadCount() //thread
	//times, _ := cpu.Times(true)

	// ============================================================================
	// Overview
	// ============================================================================

	//แยกส่วน
	var modelName string
	var vendorID string
	var core string
	var thread string
	var freqMax string
	var family string
	var modelid string
	var stepping string
	var cachet string
	var microcode string

	modelName = fmt.Sprintf("CPU : %s", info[0].ModelName)
	vendorID = fmt.Sprintf("Vendor : %s", info[0].VendorID)
	core = fmt.Sprintf("Cores : %d", physical)
	thread = fmt.Sprintf("Thread : %d", logical)
	freqMax = fmt.Sprintf("FreqMax : %.2f GHz", info[0].Mhz/1000)
	family = fmt.Sprintf("Family : %s", info[0].Family)
	modelid = fmt.Sprintf("Modelid : %s", info[0].Model)
	stepping = fmt.Sprintf("Stepping : %d", info[0].Stepping)
	cachet = fmt.Sprintf("Cache : %d MB", info[0].CacheSize/1024)
	microcode = fmt.Sprintf("Microcode : %s", info[0].Microcode)
	// ============================================================================
	// Detail
	// ============================================================================
	hyperthreading := fmt.Sprintf("Hyperthreading: [ %v ]", logical > physical)

	var cpuThreadCoreSocketresult string //จำนวน thread
	cpuThreadCoreSocketresult += ("[  Thread  ] : [ Core ] : [ Socket ]\n")
	for i, cpu := range info {
		cpuThreadCoreSocketresult += fmt.Sprintf("\nThread [%d] : Core [%s] : Socket [%s]",
			i, cpu.CoreID, cpu.PhysicalID)
	}

	// cpuid
	cpuInfo := cpuid.CPU

	c1d, xc1d := processValue(cpuInfo.Cache.L1D)
	c1i, xc1i := processValue(cpuInfo.Cache.L1I)
	c2, xc2 := processValue(cpuInfo.Cache.L2)
	c3, xc3 := processValue(cpuInfo.Cache.L3)

	var cache string //cpuid
	cache += "[ Cache ]\n"
	cache += fmt.Sprintf("\nL1d : %d %s", c1d, xc1d)
	cache += fmt.Sprintf("\nL1i : %d %s", c1i, xc1i)
	cache += fmt.Sprintf("\nL2 : %d %s", c2, xc2)
	cache += fmt.Sprintf("\nL3 : %d %s", c3, xc3)

	//"BrandName":          cpuInfo.BrandName, //ชื่อ cpu
	//"l1d_cache": cpuInfo.Cache.L1D,
	//"l1i_cache": cpuInfo.Cache.L1I,
	//"l2_cache":  cpuInfo.Cache.L2,
	//"l3_cache":  cpuInfo.Cache.L3,
	//"has_avx2": cpuInfo.Has(cpuid.AVX2),

	// ============================================================================
	// Flags Feature
	// ============================================================================
	var flagsfeature string
	//flagsfeature += "\n"
	for i, flag := range info[0].Flags {
		flagsfeature += flag
		if (i+1)%6 == 0 { // ทีละ 6 flags ต่อบรรทัด
			flagsfeature += "\n"
		} else {
			flagsfeature += " "
		}
	}

	return map[string]interface{}{
		// gopsutil
		"FlagsFeature":              flagsfeature,
		"ModelName":                 modelName,
		"VendorID":                  vendorID,
		"Core":                      core,
		"Thread":                    thread,
		"FreqMax":                   freqMax,
		"Family":                    family,
		"Modelid":                   modelid,
		"Stepping":                  stepping,
		"Cachet":                    cachet,
		"Microcode":                 microcode,
		"Hyperthreading":            hyperthreading,
		"CpuThreadCoreSocketresult": cpuThreadCoreSocketresult,
		"Cache":                     cache, //cpuid
	}
}

// ============================================================================
// monitor
// ============================================================================
type StCPUData struct {
	Usage string //
	//Timesusage string
	UsagepercentTotal         string
	UsagepercentPerCoreSTRING string
	TimesTotalAvg             string
	TimesSec                  string
	TimesHms                  string
	UsagePerCore              []float64 // CPU usage ต่อ core
	PercentPerCore            string
	Times                     []cpu.TimesStat
	//////////////////////
}
type CPUMonitor struct {
	ticker   *time.Ticker
	callback func(StCPUData)
}

// สร้าง instance ใหม่
func NewCPUMonitor(interval time.Duration, callback func(StCPUData)) *CPUMonitor {
	return &CPUMonitor{
		ticker:   time.NewTicker(interval),
		callback: callback,
	}
}

// เริ่ม monitoring
func (m *CPUMonitor) Start() {
	go func() {
		for range m.ticker.C {

			percentTotal := CpuPercentAVG()
			/*
			   			func CpuPercentAVG() []float64 {
			   	// ดึง CPU usage รวม
			   	percentAVG, err := cpu.Percent(100*time.Millisecond, false)
			   	if err != nil || len(percentAVG) == 0 {
			   		log.Println(err)
			   		return []float64{0.0}
			   	}
			   	return percentAVG
			   }
			*/

			percentPerCore := CpuPercentPercore()
			/*
			   	// ดึง CPU usage ต่อ core
			   	percentPerCore, err := cpu.Percent(100*time.Millisecond, true)
			   	if err != nil {
			   		return nil
			   	}
			   	return percentPerCore
			   }

			*/
			/*
				// ดึง CPU usage รวม
				percentTotal, err := cpu.Percent(100*time.Millisecond, false)
				if err != nil || len(percentTotal) == 0 {
					continue
				}
			*/
			/*
				// ดึง CPU usage ต่อ core
				percentPerCore, err := cpu.Percent(100*time.Millisecond, true)
				if err != nil {
					continue
				}
			*/
			//จัดเรียง usage
			usagepercentTotal := fmt.Sprintf("[ Usage Avg ] : %.2f%%", percentTotal[0])

			// แสดง usage ต่อ core
			var usagepercentPerCore string
			usagepercentPerCore += "[ Usage PerCore ]\n"
			for i, pc := range percentPerCore {
				usagepercentPerCore += fmt.Sprintf("\nCore [ %d ] : %.1f%%", i, pc)
			}

			//cpu.Times()
			times, err := cpu.Times(true)
			if err != nil || len(times) == 0 {
				continue
			}

			var timesTotalAvg string
			var timesSec string
			timesSec += "[ ข้อมูลดิบ ]"
			var timesHms string
			timesHms += "[ แปลงเป็นเวลาสากล ]"
			var totalUser float64
			var totalSystem float64
			var totalIdle float64
			var totalNice float64
			var totalIowait float64
			var totalIrq float64
			var totalSoftirq float64
			var totalSteal float64
			var totalGuest float64
			var totalGuestNice float64

			for _, d := range times {

				totalUser += d.User
				totalSystem += d.System
				totalIdle += d.Idle //รวม idle
				totalNice += d.Nice
				totalIowait += d.Iowait
				totalIrq += d.Irq
				totalSoftirq += d.Softirq
				totalSteal += d.Steal
				totalGuest += d.Guest
				totalGuestNice += d.GuestNice

				nCPU := d.CPU
				//วินาที *ดิบ
				timesSec += fmt.Sprintf(
					"\n[ %s ] | User: %.2f s | System: %.2f s | Idle: %.2f s | Nice: %.2f s | Iowait: %.2f s | Irq %.2f s | Softirq %.2f s | Steal %.2f s | Guest %.2f s | GuestNice %.2f s",
					nCPU, d.User, d.System, d.Idle, d.Nice, d.Iowait, d.Irq, d.Softirq, d.Steal, d.Guest, d.GuestNice)

				//แปลงเป็นเวลาสากล
				thUser, tmUser, tsUser := processTimeS(d.User)
				thSystem, tmSystem, tsSystem := processTimeS(d.System)
				thIdle, tmIdle, tsIdle := processTimeS(d.Idle)
				thNice, tmNice, tsNice := processTimeS(d.Nice)
				thIowait, tmIowait, tsIowait := processTimeS(d.Iowait)
				thIrq, tmIrq, tsIrq := processTimeS(d.Irq)
				thSoftirq, tmSoftirq, tsSoftirq := processTimeS(d.Softirq)
				thSteal, tmSteal, tsSteal := processTimeS(d.Steal)
				thGuest, tmGuest, tsGuest := processTimeS(d.Guest)
				thGuestNice, tmGuestNice, tsGuestNice := processTimeS(d.GuestNice)

				//จัดเรียงเวลาสากล
				timesHms += fmt.Sprintf(
					"\n[ %s ] | User [ %d : %d : %d ] | System [ %d : %d : %d ] | Idle [ %d : %d : %d ] | Nice [ %d : %d : %d ] | Iowait [ %d : %d : %d ] | Irq [ %d : %d : %d ] | Softirq [ %d : %d : %d ] | Steal [ %d : %d : %d ] | Guest [ %d : %d : %d ] | GuestNice [ %d : %d : %d ]",
					nCPU, thUser, tmUser, tsUser, thSystem, tmSystem, tsSystem, thIdle, tmIdle, tsIdle, thNice, tmNice, tsNice, thIowait, tmIowait, tsIowait, thIrq, tmIrq, tsIrq, thSoftirq, tmSoftirq, tsSoftirq, thSteal, tmSteal, tsSteal, thGuest, tmGuest, tsGuest, thGuestNice, tmGuestNice, tsGuestNice)

				//AVG//เวลาโดยเฉลี่ย
				// ***แยก system กับ idle
				hUser, mUser, sUser := Avg(totalUser)
				hSystem, mSystem, sSysteme := Avg(totalSystem)
				hIdle, mIdle, sIdle := Avg(totalIdle)
				hNice, mNice, sNice := Avg(totalNice)
				hIowait, mIowait, sIowait := Avg(totalIowait)
				hIrq, mIrq, sIrq := Avg(totalIrq)
				hSoftirq, mSoftirq, sSoftirq := Avg(totalSoftirq)
				hSteal, mSteal, sSteal := Avg(totalSteal)
				hGuest, mGuest, sGuest := Avg(totalGuest)
				hGuestNice, mGuestNice, sGuestNice := Avg(totalGuestNice)
				//จัดเรียงเวลาโดยเฉลี่ย
				timesTotalAvg = fmt.Sprintf(
					"[ เฉลี่ย ]\n[ AVG ] | User [ %d : %d : %d ] | System [ %d : %d : %d ] | Idle [ %d : %d : %d ] | Nice [ %d : %d : %d ] | Iowait [ %d : %d : %d ] | Irq [ %d : %d : %d ] | Softirq [ %d : %d : %d ] | Steal [ %d : %d : %d ] | Guest [ %d : %d : %d ] | GuestNice [ %d : %d : %d ]",
					hUser, mUser, sUser, hSystem, mSystem, sSysteme, hIdle, mIdle, sIdle, hNice, mNice, sNice, hIowait, mIowait, sIowait, hIrq, mIrq, sIrq, hSoftirq, mSoftirq, sSoftirq, hSteal, mSteal, sSteal, hGuest, mGuest, sGuest, hGuestNice, mGuestNice, sGuestNice)
			}

			if len(percentTotal) > 0 {

				data := StCPUData{
					//Usage: usage,
					//Timesusage: timesusage,
					UsagepercentTotal:         usagepercentTotal,
					UsagepercentPerCoreSTRING: usagepercentPerCore,
					TimesTotalAvg:             timesTotalAvg,
					TimesSec:                  timesSec,
					TimesHms:                  timesHms,
				}
				m.callback(data)
			}
		}
	}()
}

// ============================================================================
// cache
// ============================================================================
// ฟังก์ชันประมวลผลค่าด้วย switch case
func processValue(value int) (int, string) {
	// ตัวอักษร flag ที่สัมผัส
	var x string = ""
	// ตรวจสอบเงื่อนไข
	switch {
	case value >= 1099511627776:
		value = value / 1099511627776
		x = "TB"
	case value >= 1073741824:
		value = value / 1073741824
		x = "GB"
	case value >= 1048576:
		value = value / 1048576
		x = "MB"
	case value >= 1000:
		value = value / 1024
		x = "KB"
	default:
		x = "B"
	}
	return value, x
}

// ============================================================================
// เวลา
// ============================================================================
func processTimeS(value float64) (int, int, int) {

	hours := int(value) / 3600            // หาชั่วโมง  (int หาร int จะเป็นการหารไม่เอาเศษโดยอัตโนมัติ) *หารไม่เอาเศษ
	remainingSeconds := int(value) % 3600 //หาเศษวินาทีที่เหลือ *% หารเพื่อเอาเศษ
	minutes := remainingSeconds / 60      //  นำเศษที่เหลือมาหาหน่วยนาที *แบบไม่เอาเศษและวินาทีสุดท้าย
	seconds := remainingSeconds % 60      //และวินาทีสุดท้าย *หารเอาเศษ

	return hours, minutes, seconds
}

// ============================================================================
// หาค่าเฉลี่ย
// ============================================================================
func numSumAndCount(value []int) (int, int) {
	sum := 0
	count := 0
	for _, x := range value {
		sum += x
		if x > 0 { // ถ้ามากกว่า 0 ให้นับเพิ่ม
			count++
		}
	}
	return sum, count
}

func Avg(value float64) (int, int, int) {
	Core := CpuCoreCount()
	a1 := int(value) / int(Core)
	a2 := float64(a1)
	b1, b2, b3 := processTimeS(a2)
	return b1, b2, b3
}

// ============================================================================
// กราฟ
// ============================================================================

func grid() fyne.CanvasObject {
	coreCount := CpuCoreCount()

	colors := []color.RGBA{
		{0, 255, 0, 255},
		{0, 128, 255, 255},
		{255, 0, 0, 255},
		{255, 255, 0, 255},
		{255, 0, 255, 255},
		{0, 255, 255, 255},
		{255, 128, 0, 255},
		{128, 255, 0, 255},
	}

	cards := make([]*CoreCard, coreCount)
	items := make([]fyne.CanvasObject, coreCount)

	for i := 0; i < coreCount; i++ {
		c := NewCoreCard(i, colors[i%len(colors)])
		cards[i] = c
		items[i] = c.root
	}

	grid := container.NewGridWithColumns(1, items...)

	go func() {
		for {
			//values := getCPU()
			values := CpuPercentPercore()

			for i, c := range cards {
				if i < len(values) {
					v := values[i]
					c.graph.Update(v)
					c.val.Set(v)
				}
			}

			fyne.Do(func() {
				for _, c := range cards {
					c.raster.Refresh()
				}
			})
			//time.Sleep(100 * time.Millisecond)
			time.Sleep(100 * time.Millisecond)
		}

	}()
	return grid
}

// ============================================================================
// รวม + เอาออก CpuTabs
// ============================================================================
func CpuTabs() fyne.CanvasObject {
	grid := grid()
	dataCPUInfo := CPUdata()

	cpuOverviewPage := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["ModelName"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["VendorID"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Core"])),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Thread"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["FreqMax"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Family"])),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Modelid"])),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Stepping"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Cachet"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Microcode"])),
		widget.NewSeparator(),
	)

	cpuDetailPage := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Hyperthreading"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["CpuThreadCoreSocketresult"])),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["Cache"])), //cpuid
		widget.NewSeparator(),
	)

	cpuFlagsFeaturePage := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s", dataCPUInfo["FlagsFeature"])),
		widget.NewSeparator(),
	)
	//cpuUsagePage//
	usagepercentTotalLabel := widget.NewLabel("usagepercentTotalLabel...")
	usagepercentTotalLabel.Alignment = fyne.TextAlignCenter
	usagePerCoreSTRINGLabel := widget.NewLabel("usagePerCoreSTRINGLabel...")
	usagePerCoreSTRINGLabel.Alignment = fyne.TextAlignCenter
	//cpuTimesusagePage//
	timesTotalAvg := widget.NewLabel("timesTotalAvg...")
	timesSec := widget.NewLabel("timesSec...")
	timesHms := widget.NewLabel("timesHms...")

	// สร้าง monitor cpu
	monitor := NewCPUMonitor(1*time.Second, func(data StCPUData) {
		fyne.Do(func() {
			usagepercentTotalLabel.SetText(fmt.Sprintf("%s", data.UsagepercentTotal))          //4 // แสดง usage รวม
			usagePerCoreSTRINGLabel.SetText(fmt.Sprintf("%s", data.UsagepercentPerCoreSTRING)) //4 // แสดง usage รวม
			timesTotalAvg.SetText(fmt.Sprintf("%s", data.TimesTotalAvg))
			timesSec.SetText(fmt.Sprintf("%s", data.TimesSec))
			timesHms.SetText(fmt.Sprintf("%s", data.TimesHms))
		})
	})
	monitor.Start() // เริ่ม monitoring
	//********************************
	Grid := container.NewBorder(nil, nil, nil, nil, grid)

	cpuUsagePage := container.NewVBox(container.NewBorder(
		Grid,
		nil,
		nil,
		nil,
		container.NewCenter(
			container.NewVBox(
				usagepercentTotalLabel,
				widget.NewSeparator(),
				usagePerCoreSTRINGLabel,
				widget.NewSeparator())),
	))

	//cpuTimesusagePage
	cpuTimesusagePage := container.NewVBox(
		timesTotalAvg,
		widget.NewSeparator(),
		timesSec,
		widget.NewSeparator(),
		timesHms,
		widget.NewSeparator(),
		widget.NewLabel("[ ความหมาย ]\n[ User : CPU กำลังรันโปรแกรมทั่วไปของผู้ใช้ (โหมด user space) ]\n[ System : CPU กำลังทำงานในโหมดเคอร์เนล (เช่น ระบบเรียกไฟล์, จัดการหน่วยความจำ) ]\n[ Idle : CPU ไม่ได้ทำอะไร ไม่มีงานรอทำ ]\n[ Nice : เหมือน user mode แต่เป็นกระบวนการที่ถูกลด priority (nice value > 0) ]\n[ Iowait : ว่าง แต่มีกระบวนการรอ I/O (disk/SSD) อยู่ ถ้าค่าสูงแสดงว่า ssd ช้า ]\n[ Irq : กำลังทำงานตาม hardware interrupt (เช่น เมาส์, การ์ดเน็ตเวิร์ก) ]\n[ Softirq : งานต่อเนื่องจาก interrupt (มักเป็นงานเครือข่ายหรือ task scheduling) ]\n[ Steal : บน VM: CPU ถูก hypervisor แย่งไปให้ VM อื่น ]\n[ Guest : กำลังรัน VM อื่น (CPU ทำงานในโหมด guest OS) ]\n[ GuestNice : VM ใช้แบบ nice priority *ทำงาน แต่โดนลด priority บน host ]\n      **VM = Virtual Machine"),
		widget.NewSeparator(),
	)

	return container.NewAppTabs(
		container.NewTabItem("Overview", container.NewScroll(cpuOverviewPage)),
		container.NewTabItem("Detail", container.NewScroll(cpuDetailPage)),
		container.NewTabItem("Flags Feature", container.NewScroll(cpuFlagsFeaturePage)),
		container.NewTabItem("Usage", container.NewScroll(cpuUsagePage)),
		container.NewTabItem("TimeUsage", container.NewScroll(cpuTimesusagePage)),
	)
}
