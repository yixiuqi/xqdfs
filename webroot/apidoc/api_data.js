define({ "api": [
  {
    "type": "post",
    "url": "/opt/delete",
    "title": "[Opt]图片删除",
    "description": "<p>[Opt]图片删除</p>",
    "group": "Master",
    "version": "1.0.0",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "url",
            "description": "<p>图片url</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>0表示成功</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "string",
            "optional": false,
            "field": "img",
            "description": "<p>图片</p>"
          }
        ],
        "失败返回参数": [
          {
            "group": "失败返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>非0错误码</p>"
          },
          {
            "group": "失败返回参数",
            "type": "string",
            "optional": false,
            "field": "info",
            "description": "<p>信息</p>"
          }
        ]
      }
    },
    "filename": "./master/service/service_opt_delete.go",
    "groupTitle": "Master",
    "name": "PostOptDelete",
    "sampleRequest": [
      {
        "url": "http://ip:port/opt/delete"
      }
    ]
  },
  {
    "type": "post",
    "url": "/opt/get",
    "title": "[Opt]图片下载",
    "description": "<p>[Opt]图片下载</p>",
    "group": "Master",
    "version": "1.0.0",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "url",
            "description": "<p>图片url</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>0表示成功</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "string",
            "optional": false,
            "field": "img",
            "description": "<p>图片</p>"
          }
        ],
        "失败返回参数": [
          {
            "group": "失败返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>非0错误码</p>"
          },
          {
            "group": "失败返回参数",
            "type": "string",
            "optional": false,
            "field": "info",
            "description": "<p>信息</p>"
          }
        ]
      }
    },
    "filename": "./master/service/service_opt_get.go",
    "groupTitle": "Master",
    "name": "PostOptGet",
    "sampleRequest": [
      {
        "url": "http://ip:port/opt/get"
      }
    ]
  },
  {
    "type": "post",
    "url": "/opt/upload",
    "title": "[Opt]图片上传",
    "description": "<p>[Opt]图片上传</p>",
    "group": "Master",
    "version": "1.0.0",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "img",
            "description": "<p>图片数据</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>0表示成功</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "失败返回参数": [
          {
            "group": "失败返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>非0错误码</p>"
          },
          {
            "group": "失败返回参数",
            "type": "string",
            "optional": false,
            "field": "info",
            "description": "<p>信息</p>"
          }
        ]
      }
    },
    "filename": "./master/service/service_opt_upload.go",
    "groupTitle": "Master",
    "name": "PostOptUpload",
    "sampleRequest": [
      {
        "url": "http://ip:port/opt/upload"
      }
    ]
  },
  {
    "type": "post",
    "url": "/usage/groups",
    "title": "[Usage]所有组存储信息",
    "description": "<p>[Usage]所有组存储信息</p>",
    "group": "Master",
    "version": "1.0.0",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": true,
            "field": "type",
            "description": "<p>排序类型(非必填) <br>空:不排序 <br>sortById:按id从小到大排序 <br>sortBySize:按使用情况从小到大排序 <br>sortByWriteTps:按写TPS从小到大排序</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>0表示成功</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "失败返回参数": [
          {
            "group": "失败返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>非0错误码</p>"
          },
          {
            "group": "失败返回参数",
            "type": "string",
            "optional": false,
            "field": "info",
            "description": "<p>信息</p>"
          }
        ]
      }
    },
    "filename": "./master/service/service_usage_group.go",
    "groupTitle": "Master",
    "name": "PostUsageGroups",
    "sampleRequest": [
      {
        "url": "http://ip:port/usage/groups"
      }
    ]
  },
  {
    "type": "post",
    "url": "/store/stat",
    "title": "[Store]查询状态",
    "description": "<p>[Store]查询当前存储节点状态信息</p>",
    "group": "Storage",
    "version": "1.0.0",
    "success": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>0表示成功</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "失败返回参数": [
          {
            "group": "失败返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>非0错误码</p>"
          },
          {
            "group": "失败返回参数",
            "type": "string",
            "optional": false,
            "field": "info",
            "description": "<p>信息</p>"
          }
        ]
      }
    },
    "filename": "./storage/service/service_store_stat.go",
    "groupTitle": "Storage",
    "name": "PostStoreStat",
    "sampleRequest": [
      {
        "url": "http://ip:port/store/stat"
      }
    ]
  },
  {
    "type": "post",
    "url": "/volume/clear",
    "title": "[Volume]卷回收",
    "description": "<p>[Volume]卷回收</p>",
    "group": "Storage",
    "version": "1.0.0",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "int",
            "optional": false,
            "field": "vid",
            "description": "<p>volume id</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>0表示成功</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "失败返回参数": [
          {
            "group": "失败返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>非0错误码</p>"
          },
          {
            "group": "失败返回参数",
            "type": "string",
            "optional": false,
            "field": "info",
            "description": "<p>信息</p>"
          }
        ]
      }
    },
    "filename": "./storage/service/service_volume_clear.go",
    "groupTitle": "Storage",
    "name": "PostVolumeClear",
    "sampleRequest": [
      {
        "url": "http://ip:port/volume/clear"
      }
    ]
  },
  {
    "type": "post",
    "url": "/volume/compact",
    "title": "[Volume]卷压缩",
    "description": "<p>[Volume]卷压缩</p>",
    "group": "Storage",
    "version": "1.0.0",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "int32",
            "optional": false,
            "field": "vid",
            "description": "<p>volume id</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "成功返回参数": [
          {
            "group": "成功返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>0表示成功</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "失败返回参数": [
          {
            "group": "失败返回参数",
            "type": "int32",
            "optional": false,
            "field": "result",
            "description": "<p>非0错误码</p>"
          },
          {
            "group": "失败返回参数",
            "type": "string",
            "optional": false,
            "field": "info",
            "description": "<p>信息</p>"
          }
        ]
      }
    },
    "filename": "./storage/service/service_volume_compact.go",
    "groupTitle": "Storage",
    "name": "PostVolumeCompact",
    "sampleRequest": [
      {
        "url": "http://ip:port/volume/compact"
      }
    ]
  }
] });
