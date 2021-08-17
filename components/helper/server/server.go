package server

import (
	"context"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Server is a server interface.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type App struct {
	servers []Server
	cancel  func()

	opts options
}

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	startTimeout time.Duration
	stopTimeout  time.Duration
	signals      []os.Signal

	// eventsTimeOut 用来兜底执行事件事件，每个事件应该自己维护好各自事件
	eventsTimeOut     time.Duration
	beforeStartEvents []func() error
	afterStartEvents  []func() error
	beforeStopEvents  []func() error
	afterStopEvents   []func() error
}

// Signal with os signals.
func Signal(fn func(*App, os.Signal), sigs ...os.Signal) Option {
	return func(o *options) {
		if len(sigs) > 0 {
			o.signals = sigs
		}
	}
}
func EventsTimeOut(eventsTimeOut time.Duration) Option {
	return func(o *options) {
		o.eventsTimeOut = eventsTimeOut
	}
}

func StartTimeout(startTimeout time.Duration) Option {
	return func(o *options) {
		o.startTimeout = startTimeout
	}
}

func StopTimeout(stopTimeout time.Duration) Option {
	return func(o *options) {
		o.stopTimeout = stopTimeout
	}
}

func AppendAfterStopEvents(event func() error) Option {
	return func(o *options) {
		o.afterStopEvents = append(o.afterStopEvents, event)
	}
}

func AppendBeforeStopEvents(event func() error) Option {
	return func(o *options) {
		o.beforeStopEvents = append(o.beforeStopEvents, event)
	}
}

func AppendAfterStartEvents(event func() error) Option {
	return func(o *options) {
		o.afterStartEvents = append(o.afterStartEvents, event)
	}
}

func AppendBeforeStartEvents(event func() error) Option {
	return func(o *options) {
		o.beforeStartEvents = append(o.beforeStartEvents, event)
	}
}

func NewApp(opts ...Option) *App {
	options := options{
		startTimeout:  time.Second * 60,
		stopTimeout:   time.Second * 30,
		eventsTimeOut: time.Second * 60,
		signals: []os.Signal{
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGINT,
		},
	}
	for _, o := range opts {
		o(&options)
	}
	s := new(App)
	s.opts = options
	return s
}

func (s *App) Append(srv ...Server) {
	s.servers = append(s.servers, srv...)
}

func (s *App) Run() error {
	var ctx context.Context
	ctx, s.cancel = context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	s.beforeStart()

	wg := sync.WaitGroup{}
	for _, srv := range s.servers {
		srv := srv
		g.Go(func() error {
			<-ctx.Done()
			stopCtx, cancel := context.WithTimeout(context.Background(), s.opts.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})

		wg.Add(1)
		g.Go(func() error {
			wg.Done()
			startCtx, cancel := context.WithTimeout(context.Background(), s.opts.startTimeout)
			defer cancel()
			return srv.Start(startCtx)
		})

	}
	wg.Wait()
	s.afterStart()

	c := make(chan os.Signal, 1)
	signal.Notify(c, s.opts.signals...)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _ = <-c:
				s.beforeStop()
				s.cancel()
			}
		}
	})

	err := g.Wait()
	s.afterStop()
	return err
}

// beforeStart 服务停止前事件
func (s *App) beforeStart() {
	s.doEvents(s.opts.beforeStartEvents)
}

// afterStart 服务启动前事件
func (s *App) afterStart() {
	s.doEvents(s.opts.afterStartEvents)
}

// afterStop 服务停止后事件
func (s *App) afterStop() {
	s.doEvents(s.opts.afterStopEvents)
}

// beforeStop 服务停止后事件
func (s *App) beforeStop() {
	s.doEvents(s.opts.beforeStopEvents)
}

func (s *App) doEvents(events []func() error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.opts.eventsTimeOut)
	defer cancel()
	fin := make(chan struct{}, 1)
	go func() {
		for _, event := range events {
			_ = event()
		}
		fin <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return
	case <-fin:
		return
	}
}
