package defines

import "xqdfs/utils/stat"

type Volume struct {
	Id int32 			`json:"id"`
	LastKey int64		`json:"lastKey"`
	Compact	bool		`json:"compact"`
	ImageCount uint64	`json:"imageCount"`
	Block *Block 		`json:"block,omitempty"`
	Index *Index		`json:"index,omitempty"`
	Stat *stat.Stats	`json:"stats,omitempty"`
}

type Block struct {
	File string		`json:"file"`
	Size int64		`json:"size"`
}

type Index struct {
	File string		`json:"file"`
}
