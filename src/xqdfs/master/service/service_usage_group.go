package service

import (
	"context"
	"encoding/json"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/resource/usage"
)

func init() {
	plugin.PluginAddService(constant.CmdUsageGroups,ServiceUsageGroups)
}

/**
 * @api {post} /usage/groups [Usage]所有组存储信息
 * @apiDescription [Usage]所有组存储信息
 * @apiGroup Master
 * @apiVersion 1.0.0
 * @apiParam {string} [type] 排序类型(非必填)
						<br>空:不排序
						<br>sortById:按id从小到大排序
						<br>sortBySize:按使用情况从小到大排序
						<br>sortByWriteTps:按写TPS从小到大排序
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
type RequestUsageGroups struct {
	CmdType string	`json:"type"`
}
func ServiceUsageGroups(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestUsageGroups{"null"}
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var discoveryServer *discovery.DiscoveryServer
	if d:=plugin.PluginGetObject(plugin.PluginDiscoveryServer);d==nil {
		log.Errorf("%s no support",plugin.PluginDiscoveryServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		discoveryServer=d.(*discovery.DiscoveryServer)
	}

	groups:=discoveryServer.Groups()
	if groups==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,"discovery groups is empty")
	}

	if req.CmdType=="sortById"{
		groupsUsage:= usage.GetGroupsUsageSortById(groups)
		if groupsUsage==nil{
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,"discovery groups is empty")
		}

		json,err:=groupsUsage.Json()
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,err.Error())
		}else{
			return helper.ResultBuildWithBody(constant.Success,json)
		}
	}else if req.CmdType=="sortBySize"{
		groupsUsage:= usage.GetGroupsUsageSortBySize(groups)
		if groupsUsage==nil{
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,"discovery groups is empty")
		}

		json,err:=groupsUsage.Json()
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,err.Error())
		}else{
			return helper.ResultBuildWithBody(constant.Success,json)
		}
	}else if req.CmdType=="sortByWriteTps"{
		groupsUsage:= usage.GetGroupsUsageSortByWriteTps(groups)
		if groupsUsage==nil{
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,"discovery groups is empty")
		}

		json,err:=groupsUsage.Json()
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,err.Error())
		}else{
			return helper.ResultBuildWithBody(constant.Success,json)
		}
	}else{
		groupsUsage:= usage.GetGroupsUsage(groups)
		if groupsUsage==nil{
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,"discovery groups is empty")
		}

		json,err:=groupsUsage.Json()
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,err.Error())
		}else{
			return helper.ResultBuildWithBody(constant.Success,json)
		}
	}
}