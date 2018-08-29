package errors

const(
	RetBinlogWrite 			= 9000
	RetBinlogRead 			= 9001
	RetBinlogLength			= 9002
)

var(
	ErrBinlogWrite		= Error(RetBinlogWrite)
	ErrBinlogRead		= Error(RetBinlogRead)
	ErrBinlogLength		= Error(RetBinlogLength)
)