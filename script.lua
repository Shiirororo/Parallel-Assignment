local slot = tonumber(redis.call("GET", KEYS[1]))

if slot > 0 then
    redis.call("DECR", KEYS[1])
    return 1
end

return 0