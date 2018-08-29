package errors

const(
	RetVolumeExist     		= 7000
	RetVolumeNotExist  		= 7001
	RetVolumeDel       		= 7002
	RetVolumeInCompact 		= 7003
	RetVolumeClosed    		= 7004
	RetVolumeAdd     			= 7005
	RetVolumeAddFree			= 7006
	RetVolumeClear			= 7007
	RetVolumeTooManyCompact = 7008
	RetVolumeCompact 			= 7009
)

var(
	ErrVolumeExist     		= Error(RetVolumeExist)
	ErrVolumeNotExist  		= Error(RetVolumeNotExist)
	ErrVolumeDel       		= Error(RetVolumeDel)
	ErrVolumeInCompact 		= Error(RetVolumeInCompact)
	ErrVolumeClosed    		= Error(RetVolumeClosed)
	ErrVolumeAdd			= Error(RetVolumeAdd)
	ErrVolumeAddFree		= Error(RetVolumeAddFree)
	ErrVolumeClear			= Error(RetVolumeClear)
	ErrVolumeTooManyCompact = Error(RetVolumeTooManyCompact)
	ErrVolumeCompact 		= Error(RetVolumeCompact)
)
