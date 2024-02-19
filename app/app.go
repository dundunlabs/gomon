package app

import (
	"os"
	"os/exec"
	"syscall"
)

func New(args ...string) *App {
	return &App{
		args: args,
	}
}

type App struct {
	args []string

	cmd *exec.Cmd
}

func (app *App) Start() error {
	app.cmd = exec.Command("go", app.args...)

	app.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	app.cmd.Stderr = os.Stderr
	app.cmd.Stdout = os.Stdout

	return app.cmd.Start()
}

func (app *App) Stop() error {
	return syscall.Kill(-app.cmd.Process.Pid, syscall.SIGKILL)
}

func (app *App) Restart() error {
	if err := app.Stop(); err != nil {
		return err
	}
	return app.Start()
}
