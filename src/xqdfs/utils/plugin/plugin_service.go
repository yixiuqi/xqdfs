package plugin

type HandlerFunc func(m map[string]interface{}) interface{}
var(
	services map[string]HandlerFunc = make(map[string]HandlerFunc)
)

func PluginAddService(path string,handler HandlerFunc) {
	if path == "" || handler ==nil {
		return
	}
	services[path] = handler
}

func PluginGetServices() map[string]HandlerFunc {
	return services
}

