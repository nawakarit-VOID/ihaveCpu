package main

import (
	"fmt"
	"image/color"
	"math"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ──────────────────────────────────────────
// Colours
// ──────────────────────────────────────────
var (
	colGreen  = color.NRGBA{R: 29, G: 158, B: 117, A: 255}
	colAmber  = color.NRGBA{R: 239, G: 159, B: 39, A: 255}
	colRed    = color.NRGBA{R: 226, G: 75, B: 74, A: 255}
	colBg     = color.NRGBA{R: 30, G: 30, B: 36, A: 255}
	colCard   = color.NRGBA{R: 42, G: 42, B: 52, A: 255}
	colText   = color.NRGBA{R: 220, G: 220, B: 230, A: 255}
	colMuted  = color.NRGBA{R: 140, G: 140, B: 160, A: 255}
	colBorder = color.NRGBA{R: 60, G: 60, B: 80, A: 255}
)

// ──────────────────────────────────────────
// Metric bar widget
// ──────────────────────────────────────────
type MetricBar struct {
	widget.BaseWidget
	label   string
	percent float64
	detail  string
}

func NewMetricBar(label string) *MetricBar {
	m := &MetricBar{label: label}
	m.ExtendBaseWidget(m)
	return m
}

func (m *MetricBar) Update(pct float64, detail string) {
	m.percent = pct
	m.detail = detail
	m.Refresh()
}

func (m *MetricBar) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(colCard)
	bg.CornerRadius = 8

	barBg := canvas.NewRectangle(color.NRGBA{R: 55, G: 55, B: 70, A: 255})
	barBg.CornerRadius = 4

	barFill := canvas.NewRectangle(colGreen)
	barFill.CornerRadius = 4

	lblName := canvas.NewText(m.label, colText)
	lblName.TextStyle = fyne.TextStyle{Bold: true}
	lblName.TextSize = 14

	lblPct := canvas.NewText("0%", colGreen)
	lblPct.TextStyle = fyne.TextStyle{Bold: true}
	lblPct.TextSize = 14

	lblDetail := canvas.NewText("", colMuted)
	lblDetail.TextSize = 12

	lblStatus := canvas.NewText("ปกติ", colGreen)
	lblStatus.TextSize = 12
	lblStatus.TextStyle = fyne.TextStyle{Bold: true}

	return &metricBarRenderer{
		widget: m, bg: bg, barBg: barBg, barFill: barFill,
		lblName: lblName, lblPct: lblPct, lblDetail: lblDetail, lblStatus: lblStatus,
	}
}

type metricBarRenderer struct {
	widget    *MetricBar
	bg        *canvas.Rectangle
	barBg     *canvas.Rectangle
	barFill   *canvas.Rectangle
	lblName   *canvas.Text
	lblPct    *canvas.Text
	lblDetail *canvas.Text
	lblStatus *canvas.Text
}

func (r *metricBarRenderer) barColor() color.Color {
	p := r.widget.percent
	if p >= 85 {
		return colRed
	} else if p >= 65 {
		return colAmber
	}
	return colGreen
}

func (r *metricBarRenderer) statusText() string {
	p := r.widget.percent
	if p >= 85 {
		return "⚠ คอขวด"
	} else if p >= 65 {
		return "↑ สูง"
	}
	return "✓ ปกติ"
}

func (r *metricBarRenderer) Layout(size fyne.Size) {
	pad := float32(12)
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	topY := pad
	r.lblName.Move(fyne.NewPos(pad, topY))
	r.lblName.Resize(fyne.NewSize(size.Width*0.5, 20))

	r.lblStatus.Move(fyne.NewPos(size.Width-100, topY))
	r.lblStatus.Resize(fyne.NewSize(90, 20))

	barY := topY + 26
	barH := float32(14)
	totalW := size.Width - pad*2

	r.barBg.Move(fyne.NewPos(pad, barY))
	r.barBg.Resize(fyne.NewSize(totalW, barH))

	fillW := float32(r.widget.percent/100) * totalW
	if fillW < 0 {
		fillW = 0
	}
	r.barFill.Move(fyne.NewPos(pad, barY))
	r.barFill.Resize(fyne.NewSize(fillW, barH))

	bottomY := barY + barH + 6
	r.lblDetail.Move(fyne.NewPos(pad, bottomY))
	r.lblDetail.Resize(fyne.NewSize(size.Width*0.6, 18))

	r.lblPct.Move(fyne.NewPos(size.Width-70, bottomY))
	r.lblPct.Resize(fyne.NewSize(58, 18))
}

func (r *metricBarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 76)
}

func (r *metricBarRenderer) Refresh() {
	c := r.barColor()
	r.barFill.FillColor = c
	r.lblPct.Color = c
	r.lblStatus.Color = c
	r.lblPct.Text = fmt.Sprintf("%.1f%%", r.widget.percent)
	r.lblDetail.Text = r.widget.detail
	r.lblStatus.Text = r.statusText()
	r.lblName.Text = r.widget.label
	canvas.Refresh(r.widget)
}

func (r *metricBarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.barBg, r.barFill, r.lblName, r.lblStatus, r.lblDetail, r.lblPct}
}
func (r *metricBarRenderer) Destroy() {}

// ──────────────────────────────────────────
// Sparkline widget (mini history graph)
// ──────────────────────────────────────────
type Sparkline struct {
	widget.BaseWidget
	data   []float64
	maxPts int
	title  string
	barCol color.Color
}

func NewSparkline(title string, maxPts int) *Sparkline {
	s := &Sparkline{title: title, maxPts: maxPts, data: make([]float64, maxPts), barCol: colGreen}
	s.ExtendBaseWidget(s)
	return s
}

func (s *Sparkline) Push(v float64) {
	s.data = append(s.data[1:], v)
	if v >= 85 {
		s.barCol = colRed
	} else if v >= 65 {
		s.barCol = colAmber
	} else {
		s.barCol = colGreen
	}
	s.Refresh()
}

func (s *Sparkline) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(colCard)
	bg.CornerRadius = 8
	lbl := canvas.NewText(s.title, colMuted)
	lbl.TextSize = 12
	return &sparklineRenderer{w: s, bg: bg, lbl: lbl, lines: []*canvas.Line{}}
}

type sparklineRenderer struct {
	w     *Sparkline
	bg    *canvas.Rectangle
	lbl   *canvas.Text
	lines []*canvas.Line
}

func (r *sparklineRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))
	r.lbl.Move(fyne.NewPos(10, 6))
	r.lbl.Resize(fyne.NewSize(size.Width-20, 16))
}

func (r *sparklineRenderer) MinSize() fyne.Size { return fyne.NewSize(300, 80) }

func (r *sparklineRenderer) Refresh() {
	data := r.w.data
	n := len(data)
	if n < 2 {
		return
	}
	size := r.w.Size()
	pad := float32(10)
	topPad := float32(24)
	botPad := float32(10)
	plotH := size.Height - topPad - botPad
	plotW := size.Width - pad*2

	// build lines
	for len(r.lines) < n-1 {
		l := canvas.NewLine(r.w.barCol)
		l.StrokeWidth = 2
		r.lines = append(r.lines, l)
	}

	stepX := plotW / float32(n-1)
	for i := 0; i < n-1; i++ {
		x1 := pad + float32(i)*stepX
		x2 := pad + float32(i+1)*stepX
		y1 := topPad + plotH - float32(data[i]/100)*plotH
		y2 := topPad + plotH - float32(data[i+1]/100)*plotH
		r.lines[i].Position1 = fyne.NewPos(x1, y1)
		r.lines[i].Position2 = fyne.NewPos(x2, y2)
		r.lines[i].StrokeColor = r.w.barCol
		canvas.Refresh(r.lines[i])
	}
	canvas.Refresh(r.w)
}

func (r *sparklineRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.bg, r.lbl}
	for _, l := range r.lines {
		objs = append(objs, l)
	}
	return objs
}
func (r *sparklineRenderer) Destroy() {}

// ──────────────────────────────────────────
// Process row
// ──────────────────────────────────────────
type ProcessInfo struct {
	Name   string
	PID    int32
	CPU    float64
	Memory uint64
}

// ──────────────────────────────────────────
// App state
// ──────────────────────────────────────────
type Monitor struct {
	app    fyne.App
	window fyne.Window

	// metric bars
	barCPU  *MetricBar
	barRAM  *MetricBar
	barDisk *MetricBar
	barNet  *MetricBar

	// sparklines
	sparkCPU *Sparkline
	sparkRAM *Sparkline

	// summary labels
	lblOS     *widget.Label
	lblAlert  *widget.Label
	lblUptime *widget.Label

	// process table
	procTable *widget.Table
	procs     []ProcessInfo

	// previous net counters
	prevNetSent uint64
	prevNetRecv uint64
	prevTime    time.Time
}

func newMonitor() *Monitor {
	return &Monitor{
		prevTime: time.Now(),
	}
}

// ──────────────────────────────────────────
// Collect metrics
// ──────────────────────────────────────────
func (m *Monitor) collect() (cpuPct, ramPct, diskPct, netPct float64,
	cpuDetail, ramDetail, diskDetail, netDetail string,
	alerts []string) {

	// CPU
	cpuPercents, err := cpu.Percent(300*time.Millisecond, false)
	if err == nil && len(cpuPercents) > 0 {
		cpuPct = cpuPercents[0]
	}
	cpuInfo, _ := cpu.Info()
	cpuCount, _ := cpu.Counts(true)
	cpuName := "CPU"
	if len(cpuInfo) > 0 {
		cpuName = cpuInfo[0].ModelName
		if len(cpuName) > 32 {
			cpuName = cpuName[:32] + "…"
		}
	}
	cpuDetail = fmt.Sprintf("%s  |  %d threads", cpuName, cpuCount)

	// RAM
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		ramPct = vmStat.UsedPercent
		ramDetail = fmt.Sprintf("%.1f GB / %.1f GB  |  Swap: %.1f GB",
			float64(vmStat.Used)/1e9,
			float64(vmStat.Total)/1e9,
			float64(vmStat.SwapTotal)/1e9)
	}

	// Disk
	diskStat, err := disk.Usage("/")
	if err == nil {
		diskPct = diskStat.UsedPercent
		diskDetail = fmt.Sprintf("%.1f GB / %.1f GB  |  FS: %s",
			float64(diskStat.Used)/1e9,
			float64(diskStat.Total)/1e9,
			diskStat.Fstype)
	}

	// Network (bytes/sec)
	netStats, err := net.IOCounters(false)
	if err == nil && len(netStats) > 0 {
		now := time.Now()
		elapsed := now.Sub(m.prevTime).Seconds()
		if elapsed > 0 && m.prevNetSent > 0 {
			sent := float64(netStats[0].BytesSent-m.prevNetSent) / elapsed
			recv := float64(netStats[0].BytesRecv-m.prevNetRecv) / elapsed
			maxBps := 125_000_000.0 // assume 1 Gbps
			netPct = math.Min(100, (sent+recv)/maxBps*100)
			netDetail = fmt.Sprintf("↑ %.1f KB/s  ↓ %.1f KB/s", sent/1024, recv/1024)
		}
		m.prevNetSent = netStats[0].BytesSent
		m.prevNetRecv = netStats[0].BytesRecv
		m.prevTime = now
	}

	// Alerts
	if cpuPct >= 85 {
		alerts = append(alerts, fmt.Sprintf("⚠ CPU คอขวด: %.1f%%", cpuPct))
	}
	if ramPct >= 85 {
		alerts = append(alerts, fmt.Sprintf("⚠ RAM คอขวด: %.1f%%", ramPct))
	}
	if diskPct >= 90 {
		alerts = append(alerts, fmt.Sprintf("⚠ Disk เต็ม: %.1f%%", diskPct))
	}
	if netPct >= 85 {
		alerts = append(alerts, fmt.Sprintf("⚠ Network คับคั่ง: %.1f%%", netPct))
	}

	return
}

func (m *Monitor) collectProcesses() []ProcessInfo {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}
	var list []ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		cpuP, _ := p.CPUPercent()
		memI, _ := p.MemoryInfo()
		var memBytes uint64
		if memI != nil {
			memBytes = memI.RSS
		}
		list = append(list, ProcessInfo{
			Name: name, PID: p.Pid, CPU: cpuP, Memory: memBytes,
		})
	}
	// sort by CPU descending (simple bubble for top-15)
	for i := 0; i < len(list) && i < 30; i++ {
		for j := i + 1; j < len(list) && j < 30; j++ {
			if list[j].CPU > list[i].CPU {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	if len(list) > 15 {
		list = list[:15]
	}
	return list
}

// ──────────────────────────────────────────
// Build UI
// ──────────────────────────────────────────
func (m *Monitor) buildUI() fyne.CanvasObject {
	// OS info label
	m.lblOS = widget.NewLabel(fmt.Sprintf("OS: %s  |  Arch: %s  |  Go: %s",
		runtime.GOOS, runtime.GOARCH, runtime.Version()))
	m.lblOS.TextStyle = fyne.TextStyle{Italic: true}

	m.lblAlert = widget.NewLabel("✓ ระบบปกติ — ไม่พบคอขวด")
	m.lblAlert.TextStyle = fyne.TextStyle{Bold: true}

	m.lblUptime = widget.NewLabel("")

	// Metric bars
	m.barCPU = NewMetricBar("CPU")
	m.barRAM = NewMetricBar("RAM")
	m.barDisk = NewMetricBar("Disk")
	m.barNet = NewMetricBar("Network")

	barsSection := container.NewVBox(
		sectionHeader("  การใช้งานทรัพยากร"),
		m.barCPU, m.barRAM, m.barDisk, m.barNet,
	)

	// Sparklines
	m.sparkCPU = NewSparkline("CPU history (60s)", 60)
	m.sparkRAM = NewSparkline("RAM history (60s)", 60)

	sparksSection := container.NewVBox(
		sectionHeader("  กราฟประวัติการใช้งาน"),
		container.NewGridWithColumns(2, m.sparkCPU, m.sparkRAM),
	)

	// Process table
	headers := []string{"ชื่อโปรเซส", "PID", "CPU %", "RAM (MB)"}
	m.procTable = widget.NewTable(
		func() (int, int) { return len(m.procs) + 1, 4 },
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			lbl.Wrapping = fyne.TextTruncate
			return lbl
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			lbl := obj.(*widget.Label)
			if id.Row == 0 {
				lbl.SetText(headers[id.Col])
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			lbl.TextStyle = fyne.TextStyle{}
			idx := id.Row - 1
			if idx >= len(m.procs) {
				lbl.SetText("")
				return
			}
			p := m.procs[idx]
			switch id.Col {
			case 0:
				lbl.SetText(p.Name)
			case 1:
				lbl.SetText(fmt.Sprintf("%d", p.PID))
			case 2:
				lbl.SetText(fmt.Sprintf("%.1f%%", p.CPU))
			case 3:
				lbl.SetText(fmt.Sprintf("%.1f", float64(p.Memory)/1e6))
			}
		},
	)
	m.procTable.SetColumnWidth(0, 200)
	m.procTable.SetColumnWidth(1, 70)
	m.procTable.SetColumnWidth(2, 80)
	m.procTable.SetColumnWidth(3, 90)
	m.procTable.Resize(fyne.NewSize(460, 360))

	procSection := container.NewVBox(
		sectionHeader("  โปรเซสที่ใช้ทรัพยากรสูง (Top 15)"),
		container.NewScroll(m.procTable),
	)

	// Assemble
	topBar := container.NewHBox(m.lblOS, widget.NewSeparator(), m.lblUptime)
	alertBar := container.NewVBox(m.lblAlert)

	left := container.NewVBox(alertBar, barsSection, sparksSection)
	right := procSection

	content := container.NewGridWithColumns(2, left, right)

	return container.NewBorder(
		container.NewVBox(topBar, widget.NewSeparator()),
		statusFooter(),
		nil, nil,
		container.NewScroll(content),
	)
}

func sectionHeader(title string) *canvas.Text {
	t := canvas.NewText(title, colText)
	t.TextStyle = fyne.TextStyle{Bold: true}
	t.TextSize = 13
	return t
}

func statusFooter() fyne.CanvasObject {
	lbl := canvas.NewText("  System Bottleneck Monitor  |  อัปเดตทุก 1 วินาที  |  github.com/yourname/sysmonitor", colMuted)
	lbl.TextSize = 11
	return lbl
}

// ──────────────────────────────────────────
// Update loop
// ──────────────────────────────────────────
func (m *Monitor) startLoop() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		procTicker := time.NewTicker(3 * time.Second)
		start := time.Now()
		for {
			select {
			case <-ticker.C:
				cpuPct, ramPct, diskPct, netPct,
					cpuDetail, ramDetail, diskDetail, netDetail,
					alerts := m.collect()

				m.barCPU.Update(cpuPct, cpuDetail)
				m.barRAM.Update(ramPct, ramDetail)
				m.barDisk.Update(diskPct, diskDetail)
				m.barNet.Update(netPct, netDetail)

				m.sparkCPU.Push(cpuPct)
				m.sparkRAM.Push(ramPct)

				up := time.Since(start).Round(time.Second)
				m.lblUptime.SetText(fmt.Sprintf("Uptime: %s", up))

				if len(alerts) > 0 {
					alertText := ""
					for _, a := range alerts {
						alertText += a + "  "
					}
					m.lblAlert.SetText(alertText)
				} else {
					m.lblAlert.SetText("✓ ระบบปกติ — ไม่พบคอขวด")
				}

			case <-procTicker.C:
				m.procs = m.collectProcesses()
				m.procTable.Refresh()
			}
		}
	}()
}

// ──────────────────────────────────────────
// Main
// ──────────────────────────────────────────
func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("System Bottleneck Monitor")
	w.Resize(fyne.NewSize(1100, 720))
	w.SetFixedSize(false)

	mon := newMonitor()
	mon.app = a
	mon.window = w

	ui := mon.buildUI()
	w.SetContent(ui)

	mon.startLoop()

	w.ShowAndRun()
}
