# run javascript in go
  封装一个go http 服务，执行指定js文件中的run函数。  

## todo list
  * 按目录支持js

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
        }
    ],
    "script": "./script/es6.js"
}'
```
