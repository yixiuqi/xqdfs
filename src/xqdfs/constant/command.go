package constant

const(
	Success 	= 	0
	Cookie		=	11223344

	//------------------------------------------------------------------------ storage
	CmdVolumeAddFree 				= "/volume/addfree"
	CmdVolumeAdd	 				= "/volume/add"
	CmdVolumeClear 				= "/volume/clear"
	CmdVolumeUpload				= "/volume/upload"
	CmdVolumeGet					= "/volume/get"
	CmdVolumeDelete				= "/volume/delete"
	CmdVolumeNeedleInfo			= "/volume/needleinfo"
	CmdVolumeCompact				= "/volume/compact"
	CmdVolumeCompactStatus		= "/volume/compact/status"

	CmdStoreInit 					= "/store/init"
	CmdStoreStat 					= "/store/stat"
	CmdStoreConf 					= "/store/conf"

	//------------------------------------------------------------------------ master
	CmdGroupGetAll				= "/group/getall"
	CmdGroupReadOnly				= "/group/readonly"
	CmdGroupAdd					= "/group/add"
	CmdGroupRemove				= "/group/remove"
	CmdGroupAddStorage			= "/group/storage/add"
	CmdGroupRemoveStorage		= "/group/storage/remove"

	CmdStorageVolumeGetAll		= "/storage/volume/getall"
	CmdStorageVolumeCompact		= "/storage/volume/compact"
	CmdStorageVolumeClear		= "/storage/volume/clear"
	CmdStorageInit				= "/storage/init"
	CmdStorageAdd					= "/storage/add"
	CmdStorageRemove				= "/storage/remove"
	CmdStorageGetAll				= "/storage/getall"
	CmdStorageGetConfigure		= "/storage/getconfigure"

	CmdUsageGroups 				= "/usage/groups"

	CmdOptUpload					= "/opt/upload"
	CmdOptGet						= "/opt/get"
	CmdOptDelete					= "/opt/delete"
)
