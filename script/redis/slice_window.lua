-- 对当个 IP 进行限流
local key = KEYS[1] -- [prefix]:[ip]
local window = tonumber(ARGV[1]) -- 窗口大小
local rate = tonumber(ARGV[2]) --
local r = tonumber(ARGV[3]) -- 窗口右边界
local l = r - window -- 窗口左边界

redis.call('ZREMRANGEBYSCORE', key, '-inf', l) -- 只保留 l ~ r 的 key

local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')

if cnt >= rate then
    return "true" -- 限流
else
    redis.call("ZADD", key, r, r) -- ZADD KEY_NAME SCORE1 VALUE1
    redis.call("PEXPIRE", key, window) -- 毫秒数 用户不在访问，清理
    return "false" -- 不限流
end