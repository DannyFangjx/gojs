# run javascript in go
  封装一个go http 服务，执行指定js文件中的run函数。  

## todo list
  * 交叉编译 以及 x86_64 centos服务器上编译。
  * 安全方面待优化
  * 超时等异常情况处理。
  * js存在很多异步返回的场景，我们go工程同步返回结果不太友好。需要考虑异步返回。

## quick start
git clone
go build -o ./js_in_go ./main.go

./js_in_go


## test case
params 脚本入参
script 脚本路径
```
curl --location 'http://localhost:19800/runjs' \
--header 'Content-Type: application/json' \
--data '{
    "params": [
        "1",
        "fdasf",
        2,
        {
            "p1":"123",
            "p4": 4
        },
        {
            "te1":"1"
        }
    ],
    "script": "./script/test.js"
}'
```

### change log
* 20240310 utc+8
 1. 删除 goja的返回
 2. js脚本中 run(input) 只接受一个固定的input对象；input对象支持 all(),first(), last() 等方法。