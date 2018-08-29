package constant

const(
	Success 	= 	0
	Cookie		=	11223344

	//------------------------------------------------------------------------ storage
	HttpVolumeAddFree 			= "/volume/addfree"
	HttpVolumeAdd	 				= "/volume/add"
	HttpVolumeClear 				= "/volume/clear"
	HttpVolumeUpload				= "/volume/upload"
	HttpVolumeGet					= "/volume/get"
	HttpVolumeDelete				= "/volume/delete"
	HttpVolumeNeedleInfo			= "/volume/needleinfo"
	HttpVolumeCompact				= "/volume/compact"
	HttpVolumeCompactStatus		= "/volume/compact/status"

	HttpStoreInit 				= "/store/init"
	HttpStoreStat 				= "/store/stat"
	HttpStoreConf 				= "/store/conf"

	//------------------------------------------------------------------------ master
	HttpGroupGetAll				= "/group/getall"
	HttpGroupReadOnly				= "/group/readonly"
	HttpGroupAdd					= "/group/add"
	HttpGroupRemove				= "/group/remove"
	HttpGroupAddStorage			= "/group/storage/add"
	HttpGroupRemoveStorage		= "/group/storage/remove"

	HttpStorageVolumeGetAll		= "/storage/volume/getall"
	HttpStorageVolumeCompact	= "/storage/volume/compact"
	HttpStorageVolumeClear		= "/storage/volume/clear"
	HttpStorageInit				= "/storage/init"
	HttpStorageAdd				= "/storage/add"
	HttpStorageRemove				= "/storage/remove"
	HttpStorageGetAll				= "/storage/getall"
	HttpStorageGetConfigure		= "/storage/getconfigure"

	HttpUsageGroups 				= "/usage/groups"

	HttpOptUpload					= "/opt/upload"
	HttpOptGet						= "/opt/get"
	HttpOptDelete					= "/opt/delete"
)
