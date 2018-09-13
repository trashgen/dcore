package persistence

const (
    getRandomNodesName = "GetRandomNodes"
    getRandomNodes =
// Use case: evalsha XXX 0, <max_number_of_nodes:int>
`
    local hgetall = function (key)
      local bulk = redis.call('HGETALL', key)
        local result = {}
        local nextkey
        for i, v in ipairs(bulk) do
            if i % 2 == 1 then
                nextkey = v
            else
                result[nextkey] = v
            end
        end
        return result
    end
    local nodeKeys = redis.call("SRANDMEMBER", "nodekeys", ARGV[1])
    local result = ""
    for i, k in ipairs(nodeKeys) do
        local values = hgetall(nodeKeys[i])
        result = result .. values["key"] .. "-" .. values["ip"] .. ":" .. values["port"] .. "\t"
    end
    return result
`
)