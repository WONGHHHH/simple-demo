package main

import "simple-demo/basics/demo"

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"msgs"`
}

func main() {
	demo.JsonDemo()
}