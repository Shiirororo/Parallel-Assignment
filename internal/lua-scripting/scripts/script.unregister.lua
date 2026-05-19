-- KEYS[1] = class id
local exists = redis.call("EXISTS", KEYS[1])

if exists == 0 then
    return 0
end

redis.call("INCR", KEYS[1])
return 1
