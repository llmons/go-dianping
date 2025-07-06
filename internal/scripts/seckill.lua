-- 1. 参数列表
-- 1.1. 优惠券 id
local voucherId = ARGV[1]
-- 1.2. 用户 id
local userId = ARGV[2]

-- 2. 数据 key
-- 2.1. 库存 key
local stockKey = 'seckill:stock:' .. voucherId
-- 2.2. 订单 key
local orderKey = 'seckill:order:' .. voucherId

-- 3. 脚本业务
-- 3.1. 判断库存是否充足
if ((tonumber(redis.call('get', stockKey))) <= 0) then
    -- 3.2. 库存不足，返回1
    return 1
end
-- 3.3. 判断用户是否下单
if redis.call('sismember',orderKey,userId)==1 then
    -- 3.4. 存在，重复下单，返回2
    return 2
end
-- 3.5. 扣减库存
redis.call('decr', stockKey)
-- 3.6. 下单（保存用户到集合）
redis.call('sadd', orderKey, userId)
return 0