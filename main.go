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
	gojaOut := gojaFunc(req.Params, script)
	ret := map[string]interface{}{
		"goja_out": gojaOut,
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
	//user js code
	ctx.RunScript(script, "")
	//sys js code
	ctx.RunScript(fmt.Sprintf(`const result = run(...%v)`, jsonStr(input)), "")
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
