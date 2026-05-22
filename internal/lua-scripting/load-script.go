package luascripting

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed scripts/script.register.lua
var registerScript string

//go:embed scripts/script.get-class.lua
var getClassScript string

//go:embed scripts/script.unregister.lua
var unregisterScript string

var RegisterScript = redis.NewScript(registerScript)
var GetClassScript = redis.NewScript(getClassScript)
var UnregisterScript = redis.NewScript(unregisterScript)
