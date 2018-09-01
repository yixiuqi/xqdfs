package errors

const(
	RetGroupIsEmpty 			= 8000
	RetGroupGetAll 			= 8001
	RetGroupGet 				= 8002
	RetGroupNotExist			= 8003
	RetGroupAdd				= 8004
	RetGroupEdit				= 8005
	RetGroupRemove			= 8006
	RetGroupAddStorage		= 8007
	RetGroupEditStorage		= 8008
	RetGroupRemoveStorage	= 8009
	RetStorageVolumeGetAll	= 8010
	RetStorageNotExist		= 8011
	RetStorageExist			= 8012
	RetStorageAdd				= 8013
	RetStorageRemove			= 8014
	RetStorageGet				= 8015
	RetStorageGetAll			= 8016
	RetParamNotExist			= 8017
	RetParamSet				= 8018
)

var(
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
)
