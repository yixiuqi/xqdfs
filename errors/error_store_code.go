package errors

const(
	RetStoreVolumeIndex  	= 6000
	RetStoreNoFreeVolume 	= 6001
	RetStoreFileExist    		= 6002
	RetStoreInitFailed		= 6003
	RetStoreConfigure			= 6004
)

var(
	ErrStoreVolumeIndex  	= Error(RetStoreVolumeIndex)
	ErrStoreNoFreeVolume 	= Error(RetStoreNoFreeVolume)
	ErrStoreFileExist    	= Error(RetStoreFileExist)
	ErrStoreInitFailed   	= Error(RetStoreInitFailed)
	ErrStoreConfigure   	= Error(RetStoreConfigure)
)