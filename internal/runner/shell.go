package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type Shell struct {
	Base
	Cmds []string
}

var _ Executable = (*Shell)(nil)

func (s Shell) DryRun(ctx context.Context, w io.Writer) {
	sh, ok := os.LookupEnv("SHELL")
	if !ok {
		sh = "/bin/sh"
	}

	for _, cmd := range s.Cmds {
		_, _ = fmt.Fprintf(w, "%s in %s: %s\n", sh, s.Dir, cmd)
	}
}

func (s Shell) Run(ctx context.Context) error {
	sh, ok := os.LookupEnv("SHELL")
	if !ok {
		sh = "/bin/sh"
	}

	for _, cmd := range s.Cmds {
		if err := execSingle(ctx, sh, s.Dir, cmd, s.Stdin, s.Stdout, s.Stderr); err != nil {
			return err
		}
	}

	return nil
}

func execSingle(ctx context.Context, sh, dir, cmd string, stdin io.Reader, stdout, stderr io.Writer) error {
	c := exec.CommandContext(ctx, sh, []string{"-c", cmd}...)
	c.Dir = dir
	c.Stderr = stderr
	c.Stdout = stdout
	c.Stdin = stdin

	return errors.Wrapf(c.Run(), "failed to run command %q", cmd)
}