package main

import (
	"changeme/backend/infra"
	"changeme/backend/vo"
	"context"
	"os"
	"os/signal"
	"syscall"
)

const Lockfile = ".wows-fast-stats.lockfile"

type Lock struct {
	env        vo.Env
	logger     infra.Logger
	cancelFunc context.CancelFunc
}

func NewLock(env vo.Env, logger infra.Logger) *Lock {
	return &Lock{
		env:    env,
		logger: logger,
	}
}

func (l *Lock) Locked() bool {
	if !l.env.IsProduction() {
		return false
	}

	_, err := os.Stat(Lockfile)
	return err == nil
}

func (l *Lock) Lock() error {
	if !l.env.IsProduction() {
		return nil
	}

	if l.Locked() {
		return nil
	}

	if err := os.WriteFile(Lockfile, []byte(""), 0600); err != nil {
		l.logger.Error("Failed to create lockfile", err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	l.cancelFunc = cancel
	go l.waitShutdown(ctx)

	return nil
}

func (l *Lock) Unlock() error {
	if !l.env.IsProduction() {
		return nil
	}

	if !l.Locked() {
		return nil
	}

	if err := os.Remove(Lockfile); err != nil {
		l.logger.Error("Failed to remove lockfile", err)
		return err
	}

	defer l.cancelFunc()

	return nil
}

func (l *Lock) waitShutdown(ctx context.Context) {
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case s := <-trap:
		l.logger.Info("signal detect -> " + s.String())
		_ = l.Unlock()
		s.Signal()
	case <-ctx.Done():
		return
	}
}
