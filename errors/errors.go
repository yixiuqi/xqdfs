package errors

type Error string

func (e Error) Error() string {
	return string(e)
}

const(
	ErrNeedleExist = Error("ErrNeedleExist")
	ErrSuperBlockMagic      = Error("ErrSuperBlockMagic")
	ErrSuperBlockVer        = Error("ErrSuperBlockVer")
	ErrSuperBlockNoSpace    = Error("ErrSuperBlockNoSpace")
	ErrSuperBlockClosed     = Error("ErrSuperBlockClosed")
	ErrIndexSize   = Error("ErrIndexSize")
	ErrIndexClosed = Error("ErrIndexClosed")
	ErrIndexOffset = Error("ErrIndexOffset")
	ErrIndexEOF    = Error("ErrIndexEOF")
	ErrNeedleNotExist    = Error("ErrNeedleNotExist")
	ErrNeedleChecksum    = Error("ErrNeedleChecksum")
	ErrNeedleFlag        = Error("ErrNeedleFlag")
	ErrNeedleSize        = Error("ErrNeedleSize")
	ErrNeedleHeaderMagic = Error("ErrNeedleHeaderMagic")
	ErrNeedleFooterMagic = Error("ErrNeedleFooterMagic")
	ErrNeedleKey         = Error("ErrNeedleKey")
	ErrNeedlePadding     = Error("ErrNeedlePadding")
	ErrNeedleCookie      = Error("ErrNeedleCookie")
	ErrNeedleDeleted     = Error("ErrNeedleDeleted")
	ErrNeedleHeaderSize  = Error("ErrNeedleHeaderSize")
	ErrNeedleDataSize    = Error("ErrNeedleDataSize")
	ErrNeedleFooterSize  = Error("ErrNeedleFooterSize")
	ErrStoreVolumeIndex  = Error("ErrStoreVolumeIndex")
	ErrStoreNoFreeVolume = Error("ErrStoreNoFreeVolume")
	ErrStoreNoVolume = Error("ErrStoreNoVolume")
	ErrVolumeExist     = Error("ErrVolumeExist")
	ErrVolumeNotExist		=	Error("ErrVolumeNotExist")
	ErrVolumeInCompact = Error("ErrVolumeInCompact")
)
