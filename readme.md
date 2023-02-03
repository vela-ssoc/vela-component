# component

内置基础函数库 都是lua接口

## vela.args
- v = vela.args(0)
- 获取命令行参数

## vela.copy
- vela.copy.run(dst , src)
- vela.copy.spawn(dst , src) 异步

## vela.channel
- chl = vela.channel(10)
### 函数
- [chl.pop(function)]()
- [chl.push(function)]()
- [chl.close()]()

## vela.catch
- vela.catch(a , b , c , d) 
- 捕获错误是否为空

## sync.*
  
锁服务，提供rwmutex 和 mutx
```lua
    local mu = sync.mutex

    mu.lock()
    mu.unlock()

    
    local rw = sync.rwmutex
    rw.lock()
    rw.unlock()

    rw.read()
    rw.unread()
```

wait group接口
```lua
    local wg = sync.wait_group(100) 
    wg.add(100)
    wg.done(1)
    wg.wait()
```

## atomic

原子操作
```lua
    local a = atomic(10) 

    a.add(1)
    a.sub(10)
    print(a.value)

```




```lua
    -- time.now 处理接口
    print(time.now("mil") )   -- "2006-01-02 15:04:05.00"
    print(time.now("sec") )   -- "2006-01-02 15:04:05"
    print(time.now("min") )   -- "2006-01-02 15:04"
    print(time.now("hour"))   -- "2006-01-02.15"
    print(time.now("day") )   -- "2006-01-02"
    print(time.now("mon") )   -- "2006-01"
    print(time.now("year"))   -- "2006"

    print(time.now("2006-01-02.15.05.05.00")) -- "2022-02-22.18.13.11.33"

    local now = time.now()

    print(now.sec)    --59
    print(now.min)    --59
    print(now.hour)    --23
    print(now.day)    --30
    print(now.week)    --sunday
    print(now.month)    -- 12
    print(now.year)    -- 2022

    print(now.tt_sec)   --秒级别时间戳
    print(now.tt_milli) --毫秒级别时间戳
    print(now.tt_nano)  --纳秒级别时间戳
	print(now.today)    -- 2022-02-22

    print(now.format("2006-01-02"))    -- 2022-02-22
```

## std
- 控制台输出 支持lua.writer

```lua
    local out = std.out
    out.println(1111)
    out.print("helo")

    local err = std.err
    err.println("sss")
    err.print("helo error")
```

## cond
- 条件控制
- cond = vela.regex(v1 , v2 , v2 , ...)
- cond = vela.equal(v1 , v2 , v2 , ...)
- cond = vela.suffix(v1 , v2 , v2 , ...)
- cond = vela.prefix(v1 , v2 , v2 , ...)
- cond = vela.grep(v1 , v2 , v2 , ...)

#### 接口
- [cond.OK]()
- [cond.N(function , args...)]()
- [cond.Y(function , args...)]()
```lua

    local cond = vela.regex("aaa" , "a" , "ab")     -- true
    local cond = vela.equal("aaa" , "aaa" , "abc")  -- true
    local cond = vela.suffix("abcc-aa" ,"abc" , "abaa") -- true
    local cond = vela.grep("abcccc" , "ab*" , "acc*")   -- true
    --todo

    cond.Y(function(a) print(a) end, "y") --y
    cond.Y(function() print("x") end)     --x

    cond.Y().N()

```

## F

格式化字符串
```lua
    print(F("%s %d %v" , "helo" , 10 , userdata))
    -- 系统底层会默认调用String() 格式化对象
```