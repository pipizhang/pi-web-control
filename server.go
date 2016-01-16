package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Commander struct {
	Name        string
	Description string
	Command     string
	Args        []string
	RawOutput   []byte
	Error       error
}

func NewCommander() *Commander {
	return &Commander{}
}

func (c *Commander) Run() *Commander {
	c.RawOutput, c.Error = exec.Command(c.Command, c.Args...).Output()
	return c
}

func (c *Commander) Parse(s string) *Commander {
	arr := strings.Split(s, " ")
	c.Command = arr[0]
	c.Args = arr[1:]
	return c
}

func (c *Commander) isError() bool {
	return c.Error != nil
}

func (c *Commander) Output() string {
	return string(c.RawOutput)
}

func main() {

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	// Homepage
	m.Get("/", func(r render.Render) {
		data := make(map[string]interface{})
		data["Title"] = "Home Page"
		r.HTML(200, "index", &data)
	})
}
