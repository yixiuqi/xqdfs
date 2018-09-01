package errors

type Error int32

func (e Error) Error() string {
	return errorMsg[int32(e)]
}

var (
	errorMsg = map[int32]string{
		//common
		RetRpc:					"rpc error",
		RetImageData:				"image data error",
		RetImageSourceNotExist:	"image source host not exist",
		RetMissingParameter:		"missing parameter",
		RetParameterError:		"parameter error",
		RetOptUpload:				"upload opt error",
		RetOptGet:					"get opt error",
		RetOptDelete:				"delete opt error",
		RetNoSupport:				"no support",
		// block
		RetSuperBlockMagic:		"super block magic not match",
		RetSuperBlockVer:			"super block ver not match",
		RetSuperBlockPadding:	"super block padding not match",
		RetSuperBlockNoSpace:	"super block no left free space",
		RetSuperBlockRepairSize:	"super block repair size must equal original",
		RetSuperBlockClosed:		"super block closed",
		RetSuperBlockOffset:		"super block offset not consistency with size",
		// index
		RetIndexSize:   			"index size error",
		RetIndexClosed: 			"index closed",
		RetIndexOffset: 			"index offset",
		RetIndexEOF:    			"index eof",
		// needle
		RetNeedleExist:       	"needle already exist",
		RetNeedleNotExist:    	"needle not exist",
		RetNeedleChecksum:    	"needle data checksum not match",
		RetNeedleFlag:        	"needle flag not match",
		RetNeedleSize:        	"needle size error",
		RetNeedleHeaderMagic: 	"needle header magic not match",
		RetNeedleFooterMagic: 	"needle footer magic not match",
		RetNeedleKey:         	"needle key not match",
		RetNeedlePadding:     	"needle padding not match",
		RetNeedleCookie:      	"needle cookie not match",
		RetNeedleDeleted:     	"needle deleted",
		RetNeedleTooLarge:    	"needle has no left free space",
		RetNeedleHeaderSize:  	"needle header size",
		RetNeedleDataSize:    	"needle data size",
		RetNeedleFooterSize:  	"needle footer size",
		RetNeedlePaddingSize: 	"needle padding size",
		RetNeedleFull:        	"needle full",
		// store
		RetStoreVolumeIndex:  	"store volume index",
		RetStoreNoFreeVolume: 	"store has no free volume",
		RetStoreFileExist:    	"store rename file exist",
		RetStoreInitFailed:	 	"store init failed",
		RetStoreConfigure:		"store configure error",
		// volume
		RetVolumeExist:     		"volume exist",
		RetVolumeNotExist:  		"volume not exist",
		RetVolumeDel:       		"volume deleted",
		RetVolumeInCompact: 		"volume in compacting",
		RetVolumeClosed:    		"volume closed",
		RetVolumeAdd:     		"volume add failed",
		RetVolumeAddFree:			"volume addfree failed",
		RetVolumeClear:			"volume clear failed",
		RetVolumeTooManyCompact:"volume too many in compact",
		RetVolumeCompact:			"volume compact error",
		// configure
		RetGroupIsEmpty:			"group is empty",
		RetGroupGetAll:			"get all groups info error",
		RetGroupGet:				"get one group info error",
		RetGroupNotExist:			"group not exist",
		RetGroupAdd:				"group add error",
		RetGroupEdit:				"group edit error",
		RetGroupRemove:			"group remove error",
		RetGroupAddStorage:		"group add storage error",
		RetGroupEditStorage:		"group edit storage error",
		RetGroupRemoveStorage:	"group remove storage error",
		RetStorageVolumeGetAll:	"get all storage's volumes error",
		RetStorageNotExist:		"storage not exist",
		RetStorageExist:			"storage exist",
		RetStorageAdd:			"storage add error",
		RetStorageRemove:			"storage remove error",
		RetStorageGet:			"get storage info error",
		RetStorageGetAll:			"get all storages info error",
		RetParamNotExist:			"param not exist",
		RetParamSet:				"param set error",
		//binlog
		RetBinlogWrite:			"binlog write error",
		RetBinlogRead:			"binlog read error",
		RetBinlogLength:			"binlog length error",
	}
)
