package luascripting

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func loadScript(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

var RegisterScript = redis.NewScript(
	loadScript("./internal/lua-scripting/scripts/script.register.lua"),
)

var GetClassScript = redis.NewScript(
	loadScript("./internal/lua-scripting/scripts/script.get-class.lua"),
)

var UnregisterScript = redis.NewScript(
	loadScript("./internal/lua-scripting/scripts/script.unregister.lua"),
)
