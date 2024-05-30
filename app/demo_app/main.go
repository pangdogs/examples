package main

import "git.golaxy.org/framework"

func main() {
	// 创建app
	framework.NewApp().
		Setup("demo1", framework.ServiceGenericT[DemoService]{}).
		Setup("demo2", framework.ServiceGenericT[DemoService]{}).
		Run()
}
