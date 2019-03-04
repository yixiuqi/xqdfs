package errors

const(
	RetErrorBase  =  20060000

	//common
	RetRpc						= 0 + RetErrorBase
	RetImageData				= 1 + RetErrorBase
	RetImageSourceNotExist	= 2 + RetErrorBase
	RetMissingParameter		= 3 + RetErrorBase
	RetParameterError			= 4 + RetErrorBase
	RetOptUpload				= 5 + RetErrorBase
	RetOptGet					= 6 + RetErrorBase
	RetOptDelete				= 7 + RetErrorBase
	RetNoSupport				= 8 + RetErrorBase

	//block
	RetSuperBlockMagic		= 100 + RetErrorBase
	RetSuperBlockVer			= 101 + RetErrorBase
	RetSuperBlockPadding		= 102 + RetErrorBase
	RetSuperBlockNoSpace		= 103 + RetErrorBase
	RetSuperBlockRepairSize = 104 + RetErrorBase
	RetSuperBlockClosed		= 105 + RetErrorBase
	RetSuperBlockOffset		= 106 + RetErrorBase

	//index
	RetIndexSize				= 200 + RetErrorBase
	RetIndexClosed			= 201 + RetErrorBase
	RetIndexOffset			= 202 + RetErrorBase
	RetIndexEOF				= 203 + RetErrorBase

	//needle
	RetNeedleExist			= 300 + RetErrorBase
	RetNeedleNotExist			= 301 + RetErrorBase
	RetNeedleChecksum    		= 302 + RetErrorBase
	RetNeedleFlag        		= 303 + RetErrorBase
	RetNeedleSize        		= 304 + RetErrorBase
	RetNeedleHeaderMagic 	= 305 + RetErrorBase
	RetNeedleFooterMagic		= 306 + RetErrorBase
	RetNeedleKey         		= 307 + RetErrorBase
	RetNeedlePadding     		= 308 + RetErrorBase
	RetNeedleCookie      		= 309 + RetErrorBase
	RetNeedleDeleted     		= 310 + RetErrorBase
	RetNeedleTooLarge    		= 311 + RetErrorBase
	RetNeedleHeaderSize  	= 312 + RetErrorBase
	RetNeedleDataSize    		= 313 + RetErrorBase
	RetNeedleFooterSize  	= 314 + RetErrorBase
	RetNeedlePaddingSize 	= 315 + RetErrorBase
	RetNeedleFull        		= 316 + RetErrorBase

	//store
	RetStoreVolumeIndex  	= 400 + RetErrorBase
	RetStoreNoFreeVolume 	= 401 + RetErrorBase
	RetStoreFileExist    		= 402 + RetErrorBase
	RetStoreInitFailed		= 403 + RetErrorBase
	RetStoreConfigure			= 404 + RetErrorBase

	//volume
	RetVolumeExist     		= 500 + RetErrorBase
	RetVolumeNotExist  		= 501 + RetErrorBase
	RetVolumeDel       		= 502 + RetErrorBase
	RetVolumeInCompact 		= 503 + RetErrorBase
	RetVolumeClosed    		= 504 + RetErrorBase
	RetVolumeAdd     			= 505 + RetErrorBase
	RetVolumeAddFree			= 506 + RetErrorBase
	RetVolumeClear			= 507 + RetErrorBase
	RetVolumeTooManyCompact = 508 + RetErrorBase
	RetVolumeCompact 			= 509 + RetErrorBase

	//configure
	RetGroupIsEmpty 			= 600 + RetErrorBase
	RetGroupGetAll 			= 601 + RetErrorBase
	RetGroupGet 				= 602 + RetErrorBase
	RetGroupNotExist			= 603 + RetErrorBase
	RetGroupAdd				= 604 + RetErrorBase
	RetGroupEdit				= 605 + RetErrorBase
	RetGroupRemove			= 606 + RetErrorBase
	RetGroupAddStorage		= 607 + RetErrorBase
	RetGroupEditStorage		= 608 + RetErrorBase
	RetGroupRemoveStorage	= 609 + RetErrorBase
	RetStorageVolumeGetAll	= 610 + RetErrorBase
	RetStorageNotExist		= 611 + RetErrorBase
	RetStorageExist			= 612 + RetErrorBase
	RetStorageAdd				= 613 + RetErrorBase
	RetStorageRemove			= 614 + RetErrorBase
	RetStorageGet				= 615 + RetErrorBase
	RetStorageGetAll			= 616 + RetErrorBase
	RetParamNotExist			= 617 + RetErrorBase
	RetParamSet				= 618 + RetErrorBase

	//binlog
	RetBinlogWrite 			= 700 + RetErrorBase
	RetBinlogRead 			= 701 + RetErrorBase
	RetBinlogLength			= 702 + RetErrorBase
)

var(
	//common
	ErrRpc						= Error(RetRpc)
	ErrImageData				= Error(RetImageData)
	ErrImageSourceNotExist		= Error(RetImageSourceNotExist)
	ErrMissingParameter			= Error(RetMissingParameter)
	ErrParameterError			= Error(RetParameterError)
	ErrOptUpload				= Error(RetOptUpload)
	ErrOptGet					= Error(RetOptGet)
	ErrOptDelete				= Error(RetOptDelete)
	ErrNoSupport				= Error(RetNoSupport)

	//block
	ErrSuperBlockMagic      	= Error(RetSuperBlockMagic)
	ErrSuperBlockVer        	= Error(RetSuperBlockVer)
	ErrSuperBlockPadding    	= Error(RetSuperBlockPadding)
	ErrSuperBlockNoSpace    	= Error(RetSuperBlockNoSpace)
	ErrSuperBlockRepairSize 	= Error(RetSuperBlockRepairSize)
	ErrSuperBlockClosed     	= Error(RetSuperBlockClosed)
	ErrSuperBlockOffset     	= Error(RetSuperBlockOffset)

	//index
	ErrIndexSize  				= Error(RetIndexSize)
	ErrIndexClosed 				= Error(RetIndexClosed)
	ErrIndexOffset 				= Error(RetIndexOffset)
	ErrIndexEOF    				= Error(RetIndexEOF)

	//needle
	ErrNeedleExist    	 		= Error(RetNeedleExist)
	ErrNeedleNotExist    		= Error(RetNeedleNotExist)
	ErrNeedleChecksum    		= Error(RetNeedleChecksum)
	ErrNeedleFlag        		= Error(RetNeedleFlag)
	ErrNeedleSize        		= Error(RetNeedleSize)
	ErrNeedleHeaderMagic 		= Error(RetNeedleHeaderMagic)
	ErrNeedleFooterMagic 		= Error(RetNeedleFooterMagic)
	ErrNeedleKey         		= Error(RetNeedleKey)
	ErrNeedlePadding     		= Error(RetNeedlePadding)
	ErrNeedleCookie      		= Error(RetNeedleCookie)
	ErrNeedleDeleted     		= Error(RetNeedleDeleted)
	ErrNeedleTooLarge    		= Error(RetNeedleTooLarge)
	ErrNeedleHeaderSize  		= Error(RetNeedleHeaderSize)
	ErrNeedleDataSize    		= Error(RetNeedleDataSize)
	ErrNeedleFooterSize  		= Error(RetNeedleFooterSize)
	ErrNeedlePaddingSize		= Error(RetNeedlePaddingSize)
	ErrNeedleFull       		= Error(RetNeedleFull)

	//store
	ErrStoreVolumeIndex  		= Error(RetStoreVolumeIndex)
	ErrStoreNoFreeVolume 		= Error(RetStoreNoFreeVolume)
	ErrStoreFileExist    		= Error(RetStoreFileExist)
	ErrStoreInitFailed   		= Error(RetStoreInitFailed)
	ErrStoreConfigure   		= Error(RetStoreConfigure)

	//volume
	ErrVolumeExist     			= Error(RetVolumeExist)
	ErrVolumeNotExist  			= Error(RetVolumeNotExist)
	ErrVolumeDel       			= Error(RetVolumeDel)
	ErrVolumeInCompact 			= Error(RetVolumeInCompact)
	ErrVolumeClosed    			= Error(RetVolumeClosed)
	ErrVolumeAdd				= Error(RetVolumeAdd)
	ErrVolumeAddFree			= Error(RetVolumeAddFree)
	ErrVolumeClear				= Error(RetVolumeClear)
	ErrVolumeTooManyCompact 	= Error(RetVolumeTooManyCompact)
	ErrVolumeCompact 			= Error(RetVolumeCompact)

	//configure
	ErrGroupIsEmpty	 			= Error(RetGroupIsEmpty)
	ErrGroupGetAll 				= Error(RetGroupGetAll)
	ErrGroupGet					= Error(RetGroupGet)
	ErrGroupNotExist 			= Error(RetGroupNotExist)
	ErrGroupAdd					= Error(RetGroupAdd)
	ErrGroupEdit				= Error(RetGroupEdit)
	ErrGroupRemove 				= Error(RetGroupRemove)
	ErrGroupAddStorage			= Error(RetGroupAddStorage)
	ErrGroupEditStorage			= Error(RetGroupEditStorage)
	ErrGroupRemoveStorage 		= Error(RetGroupRemoveStorage)
	ErrStorageVolumeGetAll 		= Error(RetStorageVolumeGetAll)
	ErrStorageNotExist 			= Error(RetStorageNotExist)
	ErrStorageExist 			= Error(RetStorageExist)
	ErrStorageAdd 				= Error(RetStorageAdd)
	ErrStorageRemove 			= Error(RetStorageRemove)
	ErrStorageGet 				= Error(RetStorageGet)
	ErrStorageGetAll 			= Error(RetStorageGetAll)
	ErrParamNotExist 			= Error(RetParamNotExist)
	ErrParamSet					= Error(RetParamSet)

	//binlog
	ErrBinlogWrite				= Error(RetBinlogWrite)
	ErrBinlogRead				= Error(RetBinlogRead)
	ErrBinlogLength				= Error(RetBinlogLength)
)