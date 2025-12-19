package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"
)

type workerMode string

const (
	modeAll     workerMode = "all"
	modeInit    workerMode = "init"
	modeWebhook workerMode = "webhook"
)

func main() {
	cmd := &cli.Command{
		Name:  "workers",
		Usage: "Start N worker processes",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "n",
				Aliases:  []string{"count"},
				Required: true,
				Usage:    "number of worker processes to start",
			},
			&cli.StringFlag{
				Name:    "bin",
				Aliases: []string{"b"},
				Value:   "./worker",
				Usage:   "path to worker binary",
			},
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Value:   string(modeAll),
				Usage:   "worker mode: all | init | webhook",
			},
			&cli.DurationFlag{
				Name:  "stop-timeout",
				Value: 10 * time.Second,
				Usage: "graceful stop timeout before SIGKILL",
			},
		},
		Action: run,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, c *cli.Command) error {
	n := c.Int("n")
	if n <= 0 {
		return errors.New("n must be > 0")
	}

	bin := c.String("bin")
	if _, err := os.Stat(bin); err != nil {
		return fmt.Errorf("worker binary not found (%s): %w", bin, err)
	}

	m := workerMode(strings.TrimSpace(strings.ToLower(c.String("mode"))))
	var kinds string
	switch m {
	case modeAll:
		kinds = ""
	case modeInit:
		kinds = "initial_sync"
	case modeWebhook:
		kinds = "webhook_upsert,webhook_delete"
	default:
		return fmt.Errorf("unknown mode %q (allowed: all|init|webhook)", m)
	}

	cmds, err := startN(bin, n, kinds)
	if err != nil {
		terminateAll(cmds, 5*time.Second)
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case <-ctx.Done():
		log.Printf("context canceled, stopping workers...")
	case sig := <-sigCh:
		log.Printf("received signal: %v, stopping workers...", sig)
	}

	terminateAll(cmds, c.Duration("stop-timeout"))
	return nil
}

func startN(bin string, n int, kinds string) ([]*exec.Cmd, error) {
	cmds := make([]*exec.Cmd, 0, n)

	for i := 0; i < n; i++ {
		cmd := exec.Command(bin)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if strings.TrimSpace(kinds) != "" {
			cmd.Env = append(cmd.Env, "WORKER_KINDS="+kinds)
		}

		if err := cmd.Start(); err != nil {
			return cmds, fmt.Errorf("start worker #%d: %w", i+1, err)
		}

		log.Printf("started worker #%d pid=%d mode=%s kinds=%q", i+1, cmd.Process.Pid, kindsMode(kinds), kinds)
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func kindsMode(kinds string) string {
	k := strings.TrimSpace(kinds)
	if k == "" {
		return string(modeAll)
	}
	if strings.Contains(k, "initial_sync") && !strings.Contains(k, "webhook_") {
		return string(modeInit)
	}
	if strings.Contains(k, "webhook_") && !strings.Contains(k, "initial_sync") {
		return string(modeWebhook)
	}
	return "custom"
}
func terminateAll(cmds []*exec.Cmd, timeout time.Duration) {
	pids := make(map[int]struct{}, len(cmds))

	for _, cmd := range cmds {
		if cmd == nil || cmd.Process == nil {
			continue
		}
		pid := cmd.Process.Pid
		pids[pid] = struct{}{}
		_ = syscall.Kill(pid, syscall.SIGTERM)
	}

	deadline := time.Now().Add(timeout)

	reapNonBlocking := func() {
		for {
			var ws syscall.WaitStatus
			pid, err := syscall.Wait4(-1, &ws, syscall.WNOHANG, nil)
			if err != nil {
				return
			}
			if pid == 0 {
				return
			}
			if _, ok := pids[pid]; ok {
				delete(pids, pid)
			}
		}
	}

	for len(pids) > 0 && time.Now().Before(deadline) {
		reapNonBlocking()
		if len(pids) == 0 {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	for pid := range pids {
		_ = syscall.Kill(pid, syscall.SIGKILL)
	}

	for len(pids) > 0 {
		var ws syscall.WaitStatus
		pid, err := syscall.Wait4(-1, &ws, 0, nil)
		if err != nil {
			return
		}
		if _, ok := pids[pid]; ok {
			delete(pids, pid)
		}
	}
}
