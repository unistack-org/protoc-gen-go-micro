package main

import (
	"context"
	"fmt"
	"go/types"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func clone(srcRepo string, dstDir string) error {
	if !strings.HasPrefix(srcRepo, "https://") {
		srcRepo = "https://" + srcRepo
	}

	u, err := url.Parse(srcRepo)
	if err != nil {
		return err
	}

	var rev string
	if idx := strings.Index(u.Path, "@"); idx > 0 {
		rev = u.Path[idx+1:]
	}

	cloneOpts := &git.CloneOptions{
		URL: srcRepo,
		//	Progress: os.Stdout,
	}

	if len(rev) == 0 {
		cloneOpts.SingleBranch = true
		cloneOpts.Depth = 1
	}

	if err := cloneOpts.Validate(); err != nil {
		return err
	}

	repo, err := git.CloneContext(context.Background(), memory.NewStorage(), nil, cloneOpts)
	if err != nil {
		return err
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dstDir, os.FileMode(0755)); err != nil {
		return err
	}

	if err := cleanDir(dstDir); err != nil {
		return err
	}

	err = tree.Files().ForEach(func(file *object.File) error {
		if file == nil {
			return types.Error{Msg: "file pointer is empty"}
		}

		fmode, err := file.Mode.ToOSFileMode()
		if err != nil {
			return err
		}

		switch file.Mode {
		case filemode.Executable:
			return writeFile(file, dstDir, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fmode)
		case filemode.Dir:
			return os.MkdirAll(filepath.Join(dstDir, file.Name), fmode)
		case filemode.Regular:
			return writeFile(file, dstDir, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fmode)
		default:
			return fmt.Errorf("unsupported filetype %v for %s", file.Mode, file.Name)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
