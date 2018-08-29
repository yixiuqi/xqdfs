package defines

import "xqdfs/utils/stat"

type Volume struct {
	Id int32 			`json:"id"`
	LastKey int64		`json:"last_key"`
	Compact	bool		`json:"compact"`
	ImageCount int64	`json:"image_count"`
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
