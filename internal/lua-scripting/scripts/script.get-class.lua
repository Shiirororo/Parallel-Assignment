local result = {}

for i, key in ipairs(KEYS) do
    local value = redis.call("GET", key)
    if value then
        table.insert(result, key)
        table.insert(result, value)
    end
end

return result
