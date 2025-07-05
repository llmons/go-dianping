-- 比较标识
if ARGV[1] == redis.call("GET", KEYS[1]) then
    -- 释放锁
    return redis.call("DEL", KEYS[1])
end
return 0
