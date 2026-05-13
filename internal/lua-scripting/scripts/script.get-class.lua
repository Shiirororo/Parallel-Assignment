local result = {}

for i, key in ipairs(KEYS) do
    local value = tonumber(redis.call("GET", key) or "0")
    if value then
        table.insert(result, key)
    end
end

return result