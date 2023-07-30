local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

redis.call('ZREMRANGEBYSCORE', key, '-inf', now - window)

local prevCount = redis.call('ZCARD', key)

local first = redis.call('ZRANGE', key, 0, 0)
local reset
if next(first) == nil then
    reset = window
else
    reset = window - (now - tonumber(first[1]))
end

if prevCount >= limit then
    return { -1, reset }
end

redis.call('ZADD', key, now, now)
redis.call('EXPIRE', key, window / 1000)
return { limit - prevCount - 1, reset }
