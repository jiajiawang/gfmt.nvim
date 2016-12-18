package main

import (
	"github.com/jiajiawang/gfmt.nvim/src/sort"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "Sort", NArgs: "?", Range: ".", Eval: "*"}, sort.Sort)
		return nil
	})
}
