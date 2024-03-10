package main_test

import (
	"github.com/dop251/goja"
	"os"
	v8 "rogchap.com/v8go"
	"testing"
)

func BenchmarkGoja(b *testing.B) {
	vm := goja.New()
	filePath := "../script/add.js"
	jsData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// 执行js代码
	for i := 0; i < b.N; i++ {
		_, err = vm.RunString(string(jsData))
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkV8go(b *testing.B) {
	ctx := v8.NewContext()
	filePath := "../script/add.js"
	jsData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// 执行js代码
	for i := 0; i < b.N; i++ {
		_, err = ctx.RunScript(string(jsData), "math.js")
		if err != nil {
			panic(err)
		}
	}
}
