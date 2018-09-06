package plugin

var(
	plugins map[string]interface{} = make(map[string]interface{})
)

func PluginAddObject(key string, object interface{}) {
	if key == "" || object == nil {
		return
	}
	plugins[key] = object
}

func PluginGetObject(key string) interface{} {
	v, ok := plugins[key]
	if ok {
		return v
	} else {
		return nil
	}
}

func PluginGetObjects() map[string]interface{} {
	return plugins
}


