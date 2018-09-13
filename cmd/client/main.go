package main

import (
    "log"
    "github.com/mediocregopher/radix.v3"
    "strconv"
)

//local result = ""
//for i, k in ipairs(nodeKeys) do
//  local values = redis.call("HGETALL", nodekeys[i])
//  result = result .. values["key"] .. "-" .. values["ip"] .. ":" .. values["port"] .. "\t"
//result = result .. "lol-"
//end

const GetRandomNodes =
// Use case: evalsha <SHA_key:string> 0, <max_number_of_nodes:int>
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

func main() {
    pool, err := radix.NewPool("tcp", ":6379", 4)
    if err != nil {
        log.Fatalln(err)
    }

    scriptID := compileScript(pool, GetRandomNodes)
    log.Printf("ScriptID is [%s]", scriptID)
    result := runScript(pool, scriptID, 16)
    log.Printf("Result is [%s]", result)
    //for _, v := range result {
    //    log.Printf("Result is [%s]", v)
    //}
}

func compileScript(pool *radix.Pool, script string) (outKey string) {
    if err := pool.Do(radix.Cmd(&outKey, "SCRIPT",  "LOAD", script)); err != nil {
        log.Fatal(err)
    }
    return outKey
}

func runScript(pool *radix.Pool, scriptID string, count int) (result []byte) {
    if err := pool.Do(radix.Cmd(&result, "EVALSHA", scriptID, "0", strconv.Itoa(count))); err != nil {
        log.Fatalln(err)
    }
    return result
}