package main

import "git.golaxy.org/framework"

func main() {
	// 创建app
	framework.NewApp().
		Setup("demo1", &DemoService{}).
		Setup("demo2", &DemoService{}).
		Run()
}
