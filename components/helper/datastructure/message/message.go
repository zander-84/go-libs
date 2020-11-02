package message

import (
	"errors"
	"github.com/zander-84/go-libs/components/helper/datastructure/queue"
	"sync"
	"sync/atomic"
	"time"
)

type message struct {
	queue        queue.Queue // 消费队列
	failQueue    queue.Queue // 失败队列
	surplusQ     int64       // 剩余个数
	failSurplusQ int64       // 失败剩余个数
	cond         *sync.Cond  // 通知锁
	failcond     *sync.Cond  // 失败通知锁
	pool         sync.Pool   // 消息池
	trytimes     int         // 尝试次数
}

func NewMessage(typ int) Message {
	this := new(message)
	this.queue = queue.NewQueue(typ)
	this.failQueue = queue.NewQueue(typ)
	this.cond = sync.NewCond(new(sync.Mutex))
	this.failcond = sync.NewCond(new(sync.Mutex))
	this.pool.New = func() interface{} {
		return newMessages()
	}
	this.surplusQ = 0
	this.failSurplusQ = 0
	this.trytimes = 3
	return this
}

func (this *message) Produce(data interface{}) {
	this.queue.Push(data)
	atomic.AddInt64(&this.surplusQ, 1)
	this.cond.Broadcast()
}

func (this *message) ProducePriority(data interface{}, priority int) {
	this.queue.PushPriority(data, priority)
	atomic.AddInt64(&this.surplusQ, 1)
	this.cond.Broadcast()
}

// 消费
// 返回error的并且没有消费的，尝试3次，如果失败，放入失败队列
func (this *message) Consume(f func(messages *Messages) error) {
	for {
		this.cond.L.Lock()
		this.cond.Wait()
		messages := this.pool.Get().(*Messages)
		messages.isLock = true

		for this.queue.Size() > 0 {
			messages.clean()
			messages.Data, messages.Priority = this.queue.PopPriority()

			if messages.isLock {
				this.cond.L.Unlock()
				messages.isLock = false
			}

			if messages.Data == nil {
				continue
			}

			err := this.safeDo(messages, f)

			// 失败后重试3次
			if !messages.consumed && err != nil {
				doSucess := false
				for i := 0; i < this.trytimes; i++ {
					if err := this.safeDo(messages, f); err == nil {
						doSucess = true
						break
					}
				}
				// 没有消费和不成功情况下，尝试3次后还是失败 写到失败队列
				if !messages.consumed && !doSucess {
					this.failcond.Broadcast()
					this.failQueue.PushPriority(messages.Data, messages.Priority)
					atomic.AddInt64(&this.failSurplusQ, 1)
				}
			}
			atomic.AddInt64(&this.surplusQ, -1)
		}
		if messages.isLock {
			this.cond.L.Unlock()
		}

		// 防止漏处理
		if this.surplusQ != 0 {
			if messages.Data == nil {
				time.Sleep(time.Second) // 防止过渡消费cpu
			}
			this.cond.Broadcast()
		}

		this.pool.Put(messages)
	}
}

func (this *message) ConsumeFailQueue(f func(message *Messages) error) {
	for {
		this.failcond.L.Lock()
		this.failcond.Wait()
		messages := this.pool.Get().(*Messages)
		messages.isLock = true

		for this.failQueue.Size() > 0 {
			messages.clean()
			messages.Data, messages.Priority = this.failQueue.PopPriority()

			if messages.isLock {
				this.failcond.L.Unlock()
				messages.isLock = false
			}

			if messages.Data == nil {
				continue
			}

			_ = this.safeDo(messages, f)
			atomic.AddInt64(&this.failSurplusQ, -1)
		}

		if messages.isLock {
			this.failcond.L.Unlock()
		}

		// 防止漏处理
		if this.failSurplusQ != 0 {
			if messages.Data == nil {
				time.Sleep(time.Second) // 防止过渡消费cpu
			}
			this.failcond.Broadcast()
		}

		this.pool.Put(messages)
	}
}

func (this *message) Size() int {
	return this.queue.Size()
}

func (this *message) safeDo(val interface{}, f func(messages *Messages) error) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.New(rerr.(string))
		}
	}()
	err = f(val.(*Messages))
	return err
}
