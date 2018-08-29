package errors

const(
	RetSuperBlockMagic		= 3000
	RetSuperBlockVer			= 3001
	RetSuperBlockPadding		= 3002
	RetSuperBlockNoSpace		= 3003
	RetSuperBlockRepairSize = 3004
	RetSuperBlockClosed		= 3005
	RetSuperBlockOffset		= 3006
)

var(
	ErrSuperBlockMagic      = Error(RetSuperBlockMagic)
	ErrSuperBlockVer        = Error(RetSuperBlockVer)
	ErrSuperBlockPadding    = Error(RetSuperBlockPadding)
	ErrSuperBlockNoSpace    = Error(RetSuperBlockNoSpace)
	ErrSuperBlockRepairSize = Error(RetSuperBlockRepairSize)
	ErrSuperBlockClosed     = Error(RetSuperBlockClosed)
	ErrSuperBlockOffset     = Error(RetSuperBlockOffset)
)
