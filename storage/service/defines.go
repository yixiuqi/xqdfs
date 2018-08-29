package service

type HandlerFunc func(context *Context,m map[string]interface{}) interface{}

const(
	Success 	= 	0
	Failed 	=	1

	HttpVolumeAddFree 		= "/volume/addfree"
	HttpVolumeAdd	 			= "/volume/add"
	HttpVolumeClear 			= "/volume/clear"
	HttpVolumeUpload			= "/volume/upload"
	HttpVolumeGet				= "/volume/get"
	HttpVolumeDelete			= "/volume/delete"
	HttpVolumeNeedleInfo		= "/volume/needleinfo"
	HttpVolumeCompact			= "/volume/compact"

	HttpStoreStat 			= "/store/stat"
	HttpStoreConf 			= "/store/conf"
	HttpStoreUpload			= "/store/upload"
	HttpStoreGet				= "/store/get"

	Cookie						= int32(1982416)
)
