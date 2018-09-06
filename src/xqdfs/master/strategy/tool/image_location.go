package tool

import (
	"encoding/base64"
	"strings"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
)

type ImageLocation struct {
	GroupId int32
	StorageId int32
	VolumeId int32
	Key int64
	Cookie int32
}

func UrlToLocation(url string) (*ImageLocation,error) {
	v,err:=base64.StdEncoding.DecodeString(url)
	if err!=nil{
		log.Error(err)
		return nil,errors.ErrParameterError
	}

	items:=strings.Split(string(v),",")
	if len(items)!=5 {
		log.Error("url items!=5")
		return nil,errors.ErrParameterError
	}

	groupId,err:=helper.StringToInt32(items[0])
	if err!=nil{
		log.Error(err)
		return nil,errors.ErrParameterError
	}

	storageId,err:=helper.StringToInt32(items[1])
	if err!=nil{
		log.Error(err)
		return nil,errors.ErrParameterError
	}

	volumeId,err:=helper.StringToInt32(items[2])
	if err!=nil{
		log.Error(err)
		return nil,errors.ErrParameterError
	}

	key,err:=helper.StringToInt64(items[3])
	if err!=nil{
		log.Error(err)
		return nil,errors.ErrParameterError
	}

	cookie,err:=helper.StringToInt32(items[4])
	if err!=nil{
		log.Error(err)
		return nil,errors.ErrParameterError
	}

	addr:=&ImageLocation{
		GroupId:groupId,
		StorageId:storageId,
		VolumeId:volumeId,
		Key:key,
		Cookie:cookie,
	}

	return addr,nil
}
