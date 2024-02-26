package version

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type System struct {
	Cpu  Cpu  `json:"cpu"`
	Mem  Mem  `json:"mem"`
	Host Host `json:"host"`
	Load Load `json:"load"`
}

type Cpu struct {
	Num int `json:"num"`
	//Info []cpu.InfoStat `json:"info"`
}

type Mem struct {
	Alloc         uint64  `json:"alloc"`
	TotalAlloc    uint64  `json:"total_alloc"`
	Frees         uint64  `json:"frees"`
	NumGC         uint32  `json:"num_gc"`
	NumForcedGC   uint32  `json:"num_forced_gc"`
	LastGC        uint64  `json:"last_gc"`
	NextGC        uint64  `json:"next_gc"`
	GCCPUFraction float64 `json:"gc_cpu_fraction"`
	Lookups       uint64  `json:"lookups"`
	MemStat       string  `json:"mem_stat"`
}

type Host struct {
	Disk string         `json:"disk"`
	Info *host.InfoStat `json:"info"`
}

type Load struct {
	Avg  *load.AvgStat  `json:"avg"`
	Misc *load.MiscStat `json:"misc"`
}

func Status() *System {
	//cpuInfo, _ := cpu.Info()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}

	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil
	}

	hostInfo, _ := host.Info()

	avg, _ := load.Avg()
	misc, _ := load.Misc()

	return &System{
		Cpu: Cpu{
			Num: runtime.NumCPU(),
		},
		Mem: Mem{
			Alloc:         m.Alloc,
			TotalAlloc:    m.TotalAlloc,
			Frees:         m.Frees,
			NumGC:         m.NumGC,
			NumForcedGC:   m.NumForcedGC,
			LastGC:        m.LastGC,
			GCCPUFraction: m.GCCPUFraction,
			Lookups:       m.Lookups,
			NextGC:        m.NextGC,
			MemStat:       memStat.String(),
		},
		Host: Host{
			Info: hostInfo,
			Disk: diskUsage.String(),
		},
		Load: Load{
			Avg:  avg,
			Misc: misc,
		},
	}
}
