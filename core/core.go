package core

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/powerpuffpenguin/cronbackup/core/backend"
	"github.com/powerpuffpenguin/cronbackup/core/status"
	"github.com/robfig/cron/v3"
)

var ErrAlreadyClosed = errors.New(`core already closed`)

// Core core
type Core struct {
	*cron.Cron

	opts   coreOptions
	wait   sync.WaitGroup
	signal chan uint8
	close  chan struct{}
	once   sync.Once

	status *status.Status
}

// New new Core
func New(opt ...CoreOption) (core *Core, e error) {
	opts := defaultCoreOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	e = os.MkdirAll(opts.output, 0775)
	if e != nil {
		return
	}
	status, e := status.New(filepath.Join(opts.output, `status.db`))
	if e != nil {
		return
	}

	core = &Core{
		Cron:   cron.New(),
		opts:   opts,
		signal: make(chan uint8),
		close:  make(chan struct{}),

		status: status,
	}
	return
}

// Add add a job for crontab
func (c *Core) Add(spec string) (cron.EntryID, error) {
	return c.AddFunc(spec, c.job)
}

func (c *Core) job() {
	select {
	case c.signal <- 1:
	case <-c.close:
	default:
		log.Println(`job busy,ignore cron.Job`)
	}
}

// Serve Perform scheduled tasks
func (c *Core) Serve() (e error) {
	c.wait.Add(2)
	go c.run()
	if c.opts.immediate {
		select {
		case c.signal <- 1:
		case <-c.close:
			e = ErrAlreadyClosed
			return
		}
	}
	go c.start()
	c.wait.Wait()
	return
}
func (c *Core) start() {
	c.Start()
	c.wait.Done()
}
func (c *Core) run() {
	running := true
	for running {
		select {
		case <-c.close:
			running = false
		case <-c.signal:
			c.runJob()
		}
	}
	c.wait.Done()
}

// Stop Serve
func (c *Core) Stop() (e error) {
	closed := true
	c.once.Do(func() {
		closed = false
		close(c.close)
		e = c.Stop()
	})
	if closed {
		e = ErrAlreadyClosed
	}
	return nil
}

// UnsafeJob .
func (c *Core) UnsafeJob() {
	c.runJob()
}
func (c *Core) runJob() {
	log.Println(`runjob`)
	var (
		e      error
		ele    *status.Element
		newjob uint64
	)
	for newjob == 0 {
		ele, e = c.status.Next()
		if e != nil {
			log.Println(e)
			return
		}
		md := backend.Metadata{
			ID:       ele.ID,
			Host:     c.opts.serverHost,
			Port:     c.opts.serverPort,
			Username: c.opts.serverUsername,
			Password: c.opts.serverPassword,
			Output:   c.opts.output,
		}
		if ele.Step == status.Init {
			newjob = ele.ID
			log.Printf("new job %v\n", ele.ID)
			// do bakcup
			e = c.opts.backend.Backup(md)
			if e != nil {
				log.Println(e)
				return
			}
			// set Backuped
			ele.Step = status.Backuped
			ele.Backuped = time.Now()
			e = c.status.Update(ele)
			if e != nil {
				log.Println(e)
				return
			}
		} else {
			log.Printf("restore job %v on %v\n", ele.ID, ele.Step)
		}

		if ele.Step == status.Backuped {
			// do package
			e = c.opts.backend.Package(md)
			if e != nil {
				log.Println(e)
				return
			}
			// set Packaged
			ele.Step = status.Packaged
			e = c.status.Update(ele)
			if e != nil {
				log.Println(e)
				return
			}
		}
		if ele.Step == status.Packaged {
			// do remove expired
			e = c.opts.backend.RemoveExpired(md)
			if e != nil {
				log.Println(e)
				return
			}
			// set RemoveExpired
			ele.Step = status.RemoveExpired
			e = c.status.Update(ele)
			if e != nil {
				log.Println(e)
				return
			}
		}
		if ele.Step == status.RemoveExpired {
			// do finish
			e = c.opts.backend.Finish(md)
			if e != nil {
				log.Println(e)
				return
			}
			// set Finished
			ele.Step = status.Finished
			e = c.status.Update(ele)
			if e != nil {
				log.Println(e)
				return
			}
		}
	}
	log.Printf("runjob %v finished\n", newjob)
	if c.opts.description {
		e = c.GenerateDescription()
		if e != nil {
			log.Println(e)
		}
	}
}
func (c *Core) GenerateDescription() error {
	return c.status.GenerateDescription(filepath.Join(c.opts.output, `description.json`))
}
