local value = redis.call('GET', KEYS[1])
if not value or value ~= ARGV[1] then
    return 0
end

redis.call('DEL', KEYS[1])
redis.call('SET', KEYS[2], ARGV[2], 'PX', ARGV[3])
return 1
