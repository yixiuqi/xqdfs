package errors

const(
	RetIndexSize		= 4000
	RetIndexClosed	= 4001
	RetIndexOffset	= 4002
	RetIndexEOF		= 4003
)

var(
	ErrIndexSize   = Error(RetIndexSize)
	ErrIndexClosed = Error(RetIndexClosed)
	ErrIndexOffset = Error(RetIndexOffset)
	ErrIndexEOF    = Error(RetIndexEOF)
)

