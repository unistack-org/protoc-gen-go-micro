package main

import (
	"fmt"
	"go/token"
	"go/types"
	"os"
	"sort"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/packages"
	"google.golang.org/protobuf/compiler/protogen"
)

func (g *Generator) fieldAlign(plugin *protogen.Plugin) error {
	if !g.fieldaligment {
		return nil
	}

	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes,
		Fset: fset,
		Dir:  g.tagPath,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return fmt.Errorf("fieldalignment load: %w", err)
	}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			continue
		}

		makePass := func(a *analysis.Analyzer, resultOf map[*analysis.Analyzer]interface{}) *analysis.Pass {
			return &analysis.Pass{
				Analyzer:          a,
				Fset:              fset,
				Files:             pkg.Syntax,
				Pkg:               pkg.Types,
				TypesInfo:         pkg.TypesInfo,
				TypesSizes:        pkg.TypesSizes,
				ResultOf:          resultOf,
				Report:            func(d analysis.Diagnostic) {},
				ImportObjectFact:  func(obj types.Object, fact analysis.Fact) bool { return false },
				ExportObjectFact:  func(obj types.Object, fact analysis.Fact) {},
				ImportPackageFact: func(pkg *types.Package, fact analysis.Fact) bool { return false },
				ExportPackageFact: func(fact analysis.Fact) {},
				AllObjectFacts:    func() []analysis.ObjectFact { return nil },
				AllPackageFacts:   func() []analysis.PackageFact { return nil },
			}
		}

		inspectResult, err := inspect.Analyzer.Run(makePass(inspect.Analyzer, nil))
		if err != nil {
			return fmt.Errorf("inspect run: %w", err)
		}

		var fixes []analysis.SuggestedFix
		alignPass := makePass(fieldalignment.Analyzer, map[*analysis.Analyzer]interface{}{
			inspect.Analyzer: inspectResult,
		})
		alignPass.Report = func(d analysis.Diagnostic) {
			fixes = append(fixes, d.SuggestedFixes...)
		}

		if _, err := fieldalignment.Analyzer.Run(alignPass); err != nil {
			return fmt.Errorf("fieldalignment run: %w", err)
		}

		if err := applyTextEdits(fset, fixes); err != nil {
			return fmt.Errorf("fieldalignment apply: %w", err)
		}
	}

	return nil
}

func applyTextEdits(fset *token.FileSet, fixes []analysis.SuggestedFix) error {
	byFile := make(map[string][]analysis.TextEdit)
	for _, fix := range fixes {
		for _, edit := range fix.TextEdits {
			fname := fset.Position(edit.Pos).Filename
			byFile[fname] = append(byFile[fname], edit)
		}
	}

	for fname, edits := range byFile {
		content, err := os.ReadFile(fname)
		if err != nil {
			return err
		}

		sort.Slice(edits, func(i, j int) bool {
			return edits[i].Pos > edits[j].Pos
		})

		for _, edit := range edits {
			start := fset.Position(edit.Pos).Offset
			end := fset.Position(edit.End).Offset
			content = append(content[:start], append(edit.NewText, content[end:]...)...)
		}

		if err := os.WriteFile(fname, content, 0o644); err != nil {
			return err
		}
	}

	return nil
}
