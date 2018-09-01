package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/discovery"
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
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {string} [type] 排序类型(非必填)
						<br>空:不排序
						<br>sortById:按id从小到大排序
						<br>sortBySize:按使用情况从小到大排序
						<br>sortByWriteTps:按写TPS从小到大排序
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/usage/groups
{
	"seq": "f428b43e-5bc1-4bea-a93a-ed177ca4552e"
}
 * @apiErrorExample 失败返回示例
{
    "info": "查询错误",
    "result": 1
}

 * @apiSuccessExample 成功返回示例
{
    "info": "success",
    "result": 0
}
* */
func ServiceUsageGroups(m map[string]interface{}) interface{}{
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

	var cmdType string
	value,ok:=m["type"]
	if ok {
		cmdType=value.(string)
	}else{
		cmdType="null"
	}

	if cmdType=="sortById"{
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
	}else if cmdType=="sortBySize"{
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
	}else if cmdType=="sortByWriteTps"{
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