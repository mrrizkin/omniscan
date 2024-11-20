package app

import (
	"fmt"
	"io"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"

	"github.com/mrrizkin/omniscan/app/providers/asset"
)

type (
	Vite struct {
		entry string
		token *tokens.Token
		vite  *asset.Asset
	}

	ViteReactRefresh struct {
		token *tokens.Token
		vite  *asset.Asset
	}
)

func viteParser(asset *asset.Asset) func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	return func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
		cs := &Vite{
			token: p.Current(),
			vite:  asset,
		}

		entry := args.Match(tokens.String)
		if entry == nil {
			return nil, args.Error("A vite name is required for vite controlStructure.", nil)
		}

		cs.entry = entry.Val
		if !args.Stream().End() {
			return nil, args.Error("Malformed vite controlStructure args.", nil)
		}

		return cs, nil
	}
}

func viteReactRefreshParser(asset *asset.Asset) func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	return func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
		cs := &ViteReactRefresh{
			token: p.Current(),
			vite:  asset,
		}

		if !args.Stream().End() {
			return nil, args.Error("Malformed vite react refresh controlStructure args.", nil)
		}

		return cs, nil
	}
}

func (cs *Vite) Position() *tokens.Token {
	return cs.token
}
func (cs *Vite) String() string {
	t := cs.Position()
	return fmt.Sprintf("ViteControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (cs *Vite) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	_, err := io.WriteString(r.Output, cs.vite.Entry(cs.entry))
	return err
}

func (cs *ViteReactRefresh) Position() *tokens.Token {
	return cs.token
}
func (cs *ViteReactRefresh) String() string {
	t := cs.Position()
	return fmt.Sprintf("ViteReactRefreshControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (cs *ViteReactRefresh) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	_, err := io.WriteString(r.Output, cs.vite.ReactRefresh())
	return err
}
