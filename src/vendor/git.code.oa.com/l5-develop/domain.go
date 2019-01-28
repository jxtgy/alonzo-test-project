package l5

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	defaultDomainExpire = 30 * time.Second
	staticDomainReload  = 30 * time.Second
	overloadExpire      = 5 * time.Second
	staticDomainFiles   = []string{"/data/L5Backup/name2sid.backup", "/data/L5Backup/name2sid.cache.bin"}
	domainss            *domains
	anonymouss          *anonymous

	KGetSidFaild = errors.New("get sid failed")
)

type domains struct {
	store map[string]*Domain
	l     sync.Mutex
}

// 这个实现认为sid对应的mod/cmd是不会过期的
func (m *domains) Query(sid string) (*Domain, error) {
	now := time.Now()
	m.l.Lock()
	domain, exists := m.store[sid]
	if exists && domain.networking {
		m.l.Unlock()
		if domain.mod == 0 {
			domain.sidLock.Lock() // 只是为了获取锁, 下面已经把它锁住
			domain.sidLock.Unlock()

			if domain.mod != 0 {
				return domain, nil
			} else {
				return nil, KGetSidFaild
			}
		}
		return domain, nil
	}

	domain = &Domain{
		name:           sid,
		mod:            0,
		cmd:            0,
		expire:         now,
		overloadExpire: now.Add(overloadExpire),
		balancer:       NewBalancer(defaultBalancer),
		inited:         false,
		networking:     true,
	}
	m.store[sid] = domain

	domain.sidLock.Lock() // 在lock之前,不能让它在上面的某个地方lock住.所以先lock再unlock
	m.l.Unlock()
	defer domain.sidLock.Unlock()

	buf, err := dial(QOS_CMD_QUERY_SNAME, 0, domain.mod, domain.cmd, int32(os.Getpid()), int32(len(domain.name)), domain.name)
	if err != nil {
		domain.networking = false
		return nil, err
	}
	domain.mod = int32(defaultEndian.Uint32(buf[0:4]))
	domain.cmd = int32(defaultEndian.Uint32(buf[4:8]))
	domain.networking = true
	return domain, nil
}

func (m *domains) interval() {
	interval := time.NewTicker(staticDomainReload)
	var now time.Time
	for {
		select {
		case <-interval.C:
			var (
				err error
				fp  *os.File
			)
			now = time.Now()
			for _, v := range staticDomainFiles {
				if fp, err = os.Open(v); err != nil {
					// log.Printf("open file failed: %s", err.Error())
					continue
				}
				for {
					var (
						name string
						mod  int32
						cmd  int32
					)
					if n, fail := fmt.Fscanln(fp, &name, &mod, &cmd); n == 0 || fail != nil {
						break
					}
					m.l.Lock()
					_, exists := m.store[name]
					if !exists {
						m.store[name] = &Domain{
							name:           name,
							mod:            mod,
							cmd:            cmd,
							expire:         now,
							overloadExpire: now.Add(overloadExpire),
							balancer:       NewBalancer(defaultBalancer),
							inited:         false,
							networking:     true,
						}
					}
					m.l.Unlock()
				}
				fp.Close()
			}
		}
	}
}

type Domain struct {
	name string
	mod  int32
	cmd  int32

	l        sync.RWMutex
	expire   time.Time
	balancer Balancer

	// [MOD/CMD模式使用]下面的变量都是第一次初始化临时使用的
	inited         bool
	overloadExpire time.Time
	updateLock     sync.Mutex // 为了不和Lock冲突

	// [sid模式下使用]
	sidLock    sync.Mutex
	networking bool
}

func (d *Domain) Mod() int32 {
	return d.mod
}

func (d *Domain) Cmd() int32 {
	return d.cmd
}

func (d *Domain) Name() string {
	return d.name
}

type anonymous struct {
	store map[int32]map[int32]*Domain
	sync.Mutex
}

func (a *anonymous) Get(mod int32, cmd int32) *Domain {
	a.Lock()
	if _, exists := a.store[mod]; !exists {
		a.store[mod] = make(map[int32]*Domain)
	}

	domain, exists := a.store[mod][cmd]
	if exists {
		a.Unlock()
		return domain
	}
	now := time.Now()
	a.store[mod][cmd] = &Domain{
		name:           "",
		mod:            mod,
		cmd:            cmd,
		expire:         now,
		overloadExpire: now.Add(overloadExpire),
		balancer:       NewBalancer(defaultBalancer),
		inited:         false,
	}
	a.Unlock()
	return a.store[mod][cmd]
}

func (d *Domain) SetBalancer(b Balancer) {
	d.l.Lock()
	d.balancer = b
	d.l.Unlock()
}

func init() {
	domainss = &domains{store: make(map[string]*Domain)}
	anonymouss = &anonymous{store: make(map[int32]map[int32]*Domain)}
	go domainss.interval()
}
