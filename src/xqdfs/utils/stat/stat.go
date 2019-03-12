package stat

import "encoding/json"

type Stats struct {
	TotalReadProcessed       	uint64 		`json:"-"`	//`json:"totalRead"`
	lastTotalReadProcessed   	uint64 		`json:"-"`
	ReadQPS                  	uint64 		`json:"readQPS,omitempty"`
	TotalReadBytes           	uint64 		`json:"-"`
	lastTotalReadBytes       	uint64 		`json:"-"`
	ReadFlow                 	uint64 		`json:"readFlow,omitempty"`

	TotalWriteProcessed     	uint64 		`json:"-"`	//`json:"totalWrite"`
	lastTotalWriteProcessed 	uint64 		`json:"-"`
	WriteTPS                	uint64 		`json:"writeTPS,omitempty"`
	TotalWriteBytes				uint64 		`json:"-"`
	lastTotalWriteBytes			uint64 		`json:"-"`
	WriteFlow					uint64 		`json:"writeFlow,omitempty"`
	LastWriteTime				int64		`json:"lastWriteTime"`

	TotalDelProcessed       	uint64 		`json:"totalDel,omitempty"`
}

func (s *Stats) Calc() {
	// qps & tps
	s.WriteTPS = s.TotalWriteProcessed - s.lastTotalWriteProcessed
	s.lastTotalWriteProcessed = s.TotalWriteProcessed
	s.ReadQPS = s.TotalReadProcessed - s.lastTotalReadProcessed
	s.lastTotalReadProcessed = s.TotalReadProcessed
	// bytes
	s.ReadFlow = s.TotalReadBytes - s.lastTotalReadBytes
	s.lastTotalReadBytes = s.TotalReadBytes
	s.WriteFlow = s.TotalWriteBytes - s.lastTotalWriteBytes
	s.lastTotalWriteBytes = s.TotalWriteBytes
	return
}

// Merge merge other stats.
func (s *Stats) Merge(s1 *Stats) {
	// qps & tps
	s.TotalWriteProcessed += s1.TotalWriteProcessed
	s.TotalReadProcessed += s1.TotalReadProcessed
	s.TotalDelProcessed += s1.TotalDelProcessed
	// bytes
	s.TotalReadBytes += s1.TotalReadBytes
	s.TotalWriteBytes += s1.TotalWriteBytes
}

// Reset reset the stat.
func (s *Stats) Reset() {
	// qps & tps
	s.TotalWriteProcessed = 0
	s.TotalReadProcessed = 0
	s.TotalDelProcessed = 0
	// bytes
	s.TotalReadBytes = 0
	s.TotalWriteBytes = 0
}

func (s *Stats) String() string {
	data,err:=json.Marshal(s)
	if err!=nil {
		return ""
	}else{
		return string(data)
	}
}
