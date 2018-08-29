package errors

const(
	RetRpc						= 2000
	RetImageData				= 2001
	RetImageSourceNotExist	= 2002
	RetMissingParameter		= 2003
	RetParameterError			= 2004
	RetOptUpload				= 2005
	RetOptGet					= 2006
	RetOptDelete				= 2007
)

var(
	ErrRpc						= Error(RetRpc)
	ErrImageData				= Error(RetImageData)
	ErrImageSourceNotExist		= Error(RetImageSourceNotExist)
	ErrMissingParameter			= Error(RetMissingParameter)
	ErrParameterError			= Error(RetParameterError)
	ErrOptUpload				= Error(RetOptUpload)
	ErrOptGet					= Error(RetOptGet)
	ErrOptDelete				= Error(RetOptDelete)
)