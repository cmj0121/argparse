package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/cmj0121/argparse"
)

type IFace struct {
	argparse.Help

	IFace net.Interface

	*net.Interface `name:"iface"`
}

func main() {
	c := IFace{}
	parser := argparse.MustNew(&c)
	if err := parser.Run(); err == nil {
		data, _ := json.MarshalIndent(c, "", "    ")
		fmt.Println(string(data))
	}
}
