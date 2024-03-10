package main

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/robertkrimen/otto" // 导入Otto包
	"net/http"
	"os"
	v8 "rogchap.com/v8go"
)

func jsonStr(in interface{}) string {
	b, _ := json.Marshal(in)
	return string(b)
}

type Req struct {
	Params []interface{} `json:"params"`
	Script string        `json:"script"`
}

func runJsHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			println(err)
		}
	}()

	// req
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get script
	content, err := os.ReadFile(req.Script)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	script := string(content)

	// run
	v8goOut := v8goFunc(req.Params, script)
	//gojaOut := gojaFunc(req.Params, script)
	ret := map[string]interface{}{
		//"goja_out": gojaOut,
		"v8go_out": v8goOut,
	}
	fmt.Fprintf(w, "%s", jsonStr(ret))
}

func main() {
	http.HandleFunc("/runjs", runJsHandler)
	if err := http.ListenAndServe(":19800", nil); err != nil {
		panic(err)
	}
}

func gojaFunc(input []interface{}, script string) string {
	vm := goja.New()
	// 执行 JavaScript 字符串
	_, _ = vm.RunString(script)

	// 定义 req 对象
	req := vm.ToValue(map[string]interface{}{
		"item": input,
		"all": func(call goja.FunctionCall) goja.Value {
			// 返回所有对象
			data := call.Argument(0).ToObject(vm).Get("item").Export().([]interface{})
			return vm.ToValue(data)
		},
		"first": func(call goja.FunctionCall) goja.Value {
			// 返回第一个对象
			data := call.Argument(0).ToObject(vm).Get("item").Export().([]interface{})
			if len(data) > 0 {
				return vm.ToValue(data[0])
			}
			return nil
		},
		"last": func(call goja.FunctionCall) goja.Value {
			// 返回最后一个对象
			data := call.Argument(0).ToObject(vm).Get("item").Export().([]interface{})
			if len(data) > 0 {
				return vm.ToValue(data[len(data)-1])
			}
			return nil
		},
	})

	// 将 req 对象设置为全局对象
	vm.Set("req", req)

	// 调用 JavaScript 函数 run
	_, err := vm.RunString("const result = run(req); result;")
	if err != nil {
		return err.Error()
	}

	// 返回结果
	result := vm.Get("result")
	return result.String()
}

func gojaFunc1(input []interface{}, script string) string {
	vm := goja.New()
	// 执行js字符串
	_, _ = vm.RunString(script)
	// 校验函数并返回
	runfunc, ok := goja.AssertFunction(vm.Get("run"))
	if !ok {
		panic("not function")
	}
	// 将 input 转换为 JavaScript 中的值
	jsInput := make([]goja.Value, len(input))
	for i, v := range input {
		jsInput[i] = vm.ToValue(v)
	}
	// 调用 JavaScript 函数 run
	result, err := runfunc(nil, jsInput...)
	if err != nil {
		return err.Error()
	}
	return result.String()
}

func v8goFunc(input []interface{}, script string) string {
	ctx := v8.NewContext()

	//禁止访问 require 和 process 对象。 todo 待完善。
	_, _ = ctx.RunScript("delete require; delete process;", "")

	//user js code
	ctx.RunScript(script, "")
	//sys js code
	ctx.RunScript(fmt.Sprintf(`
				args = %v 
				let req = {
					// item
					item: args,
			
					// 返回所有对象
					all: function() {
						return args;
					},
					// 返回第一个对象
					first: function() {
						return args.length > 0 ? args[0] : undefined;
					},
					// 返回最后一个对象
					last: function() {
						return args.length > 0 ? args[args.length - 1] : undefined;
					},
					// 环境参数，和本demo无关。
					params: {
						'limits': 10,
					}
					// json  todo
					// input.context.noItemsLeft 和本demo无关。
				};
				const result = run(req)`, jsonStr(input)), "")
	val, err := ctx.RunScript("result", "value.js")
	if err != nil {
		println("run err:", err.Error())
		return ""
	} else {
		return val.String()
	}
}

func ottoFunc() {
	vm := otto.New() // 创建一个Otto实例

	// 定义一个Go函数，用于向js注入上下文
	hello := func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s!\n", call.Argument(0).String())
		return otto.Value{}
	}

	// 将Go函数注册到Otto实例中，命名为helloWorld
	vm.Set("helloWorld", hello)

	// 执行一段js代码，调用helloWorld函数，并传递参数Alice
	vm.Run(`
        helloWorld("Alice");
    `)
}
