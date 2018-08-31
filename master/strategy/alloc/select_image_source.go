package order

import (
	"xqdfs/discovery/defines"
	"xqdfs/master/strategy/tool"
	"xqdfs/errors"
	strategydef "xqdfs/master/strategy/defines"
)

func selectImageSource(groups []*defines.Group,location *tool.ImageLocation) (*strategydef.ImageSource,error) {
	var imageSource *strategydef.ImageSource

	for _,g:=range groups {
		if location.GroupId == g.Id {
			for _,s:=range g.Storage {
				if location.StorageId == s.Id && s.Online == true {
					imageSource=&strategydef.ImageSource{
						Host:s.Addr,
					}
					return imageSource,nil
				}
			}
		}
	}
	return nil,errors.ErrImageSourceNotExist
}
