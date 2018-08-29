package errors

const(
	RetNeedleExist			= 5000
	RetNeedleNotExist			= 5001
	RetNeedleChecksum    		= 5002
	RetNeedleFlag        		= 5003
	RetNeedleSize        		= 5004
	RetNeedleHeaderMagic 	= 5005
	RetNeedleFooterMagic		= 5006
	RetNeedleKey         		= 5007
	RetNeedlePadding     		= 5008
	RetNeedleCookie      		= 5009
	RetNeedleDeleted     		= 5010
	RetNeedleTooLarge    		= 5011
	RetNeedleHeaderSize  	= 5012
	RetNeedleDataSize    		= 5013
	RetNeedleFooterSize  	= 5014
	RetNeedlePaddingSize 	= 5015
	RetNeedleFull        		= 5016
)

var(
	ErrNeedleExist    	 = Error(RetNeedleExist)
	ErrNeedleNotExist    = Error(RetNeedleNotExist)
	ErrNeedleChecksum    = Error(RetNeedleChecksum)
	ErrNeedleFlag        = Error(RetNeedleFlag)
	ErrNeedleSize        = Error(RetNeedleSize)
	ErrNeedleHeaderMagic = Error(RetNeedleHeaderMagic)
	ErrNeedleFooterMagic = Error(RetNeedleFooterMagic)
	ErrNeedleKey         = Error(RetNeedleKey)
	ErrNeedlePadding     = Error(RetNeedlePadding)
	ErrNeedleCookie      = Error(RetNeedleCookie)
	ErrNeedleDeleted     = Error(RetNeedleDeleted)
	ErrNeedleTooLarge    = Error(RetNeedleTooLarge)
	ErrNeedleHeaderSize  = Error(RetNeedleHeaderSize)
	ErrNeedleDataSize    = Error(RetNeedleDataSize)
	ErrNeedleFooterSize  = Error(RetNeedleFooterSize)
	ErrNeedlePaddingSize = Error(RetNeedlePaddingSize)
	ErrNeedleFull        = Error(RetNeedleFull)
)
