local failures = tonumber(redis.call('HGET', KEYS[1], 'failures') or '0')
local level = tonumber(redis.call('HGET', KEYS[1], 'level') or '0')
local lastFailure = tonumber(redis.call('HGET', KEYS[1], 'last_failure') or '0')
local lockedUntil = tonumber(redis.call('HGET', KEYS[1], 'locked_until') or '0')
local now = tonumber(ARGV[1])
local maxFailures = tonumber(ARGV[2])
local resetAfter = tonumber(ARGV[3])

if lockedUntil > now then
    return {0, level, lockedUntil}
end

if lastFailure == 0 or now - lastFailure >= resetAfter then
    failures = 0
    level = 0
end

failures = failures + 1
lockedUntil = 0

if failures >= maxFailures then
    level = math.min(level + 1, 5)
    lockedUntil = now + tonumber(ARGV[3 + level])
    failures = 0
end

redis.call('HSET', KEYS[1], 'failures', failures)
redis.call('HSET', KEYS[1], 'level', level)
redis.call('HSET', KEYS[1], 'last_failure', now)
redis.call('HSET', KEYS[1], 'locked_until', lockedUntil)
redis.call('EXPIRE', KEYS[1], resetAfter + tonumber(ARGV[8]))
return {failures, level, lockedUntil}
