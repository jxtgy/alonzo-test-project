package l5

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	staticServerFiles   = []string{"/data/L5Backup/current_route.backup", "/data/L5Backup/current_route_v2.backup"}
	defaultServerExpire = 30 * time.Second
	KL5Overload         = errors.New("l5 overload")
)

func (d *Domain) getFromAgentAndLocal() error {
	// 从agent中获取
	var buf []byte
	var list []*Server
	var err error
	var tmpIp int32
	var tmpPort int16
	buf, err = dial(QOS_CMD_BATCH_GET_ROUTE_WEIGHT, 0, d.mod, d.cmd, int32(os.Getpid()), int32(gVersion))
	if err == nil && len(buf) > 16 {
		size := len(buf) - 16
		list = make([]*Server, size/14)
		for k, _ := range list {
			tmpIp = int32(defaultEndian.Uint32(buf[16+k*14 : 20+k*14]))
			tmpPort = int16(defaultEndian.Uint16(buf[20+k*14 : 22+k*14]))
			list[k] = &Server{
				domain:  d,
				ip:      tmpIp,
				port:    tmpPort,
				weight:  int32(defaultEndian.Uint32(buf[22+k*14 : 26+k*14])),
				total:   int32(defaultEndian.Uint32(buf[26+k*14 : 30+k*14])),
				strIp:   Ip2StringLittle(tmpIp),
				intPort: int(uint16(tmpPort)),
			}
		}
		if len(list) < 1 {
			return ErrNotFound
		}
		if err = d.Set(list); err != nil {
			return err
		}

		return nil
	}

	// 从静态文件中获取
	var fp *os.File
	list = nil

	var (
		mod  int32
		cmd  int32
		ip   string
		port int16
		n    int
		fail error
	)
	for _, v := range staticServerFiles {
		if fp, err = os.Open(v); err != nil {
			continue
		}
		for {
			if n, fail = fmt.Fscanln(fp, &mod, &cmd, &ip, &port); n == 0 || fail != nil {
				break
			}
			if d.mod != mod || d.cmd != cmd {
				continue
			}
			list = append(list, &Server{
				domain:  d,
				ip:      String2IpLittle(ip),
				port:    HostInt16ToLittle(port),
				weight:  100, //default weight: 100
				total:   0,
				strIp:   ip,
				intPort: int(uint16(port)),
			})
		}
		fp.Close()
	}

	if len(list) < 1 {
		return ErrNotFound
	}
	if err = d.Set(list); err != nil {
		return err
	}
	return nil
}

func (d *Domain) Get() (*Server, error) {
	d.l.RLock() // protect balancer
	if d.balancer == nil {
		d.l.RUnlock()
		return nil, ErrNotBalancer
	}
	// 从balancer中获取
	var srv *Server
	var err error
	srv, err = d.balancer.Get()
	d.l.RUnlock()

	now := time.Now()
	if err == nil {
		var isExpired bool
		d.updateLock.Lock()
		isExpired = now.After(d.expire)
		if isExpired {
			d.expire = now.Add(defaultServerExpire)
		}
		d.updateLock.Unlock()

		if isExpired {
			go d.getFromAgentAndLocal()
		}

		return srv.allocate(), nil
	} else if err == ErrNotFound {
		// 判断是会否已经初始化
		// 如果已初始化, 那么就是overload了, 每OverloadExpire给予更新一次l5的机会
		if d.inited && !now.After(d.overloadExpire) {
			return nil, KL5Overload
		}

		d.inited = false
		// 第一次调用, 未初始化, 请求都阻塞在这里, 只发起一次更新操作
		err = nil
		d.updateLock.Lock()
		if !d.inited {
			d.expire = now.Add(defaultServerExpire)
			err = d.getFromAgentAndLocal()
			d.overloadExpire = now.Add(overloadExpire)
			d.inited = true
		}
		d.updateLock.Unlock()

		if err != nil {
			return nil, err
		}

		d.l.RLock() // protect balancer
		if d.balancer == nil {
			d.l.RUnlock()
			return nil, ErrNotBalancer
		}
		srv, err = d.balancer.Get()
		d.l.RUnlock()

		if err != nil {
			return nil, err
		} else {
			return srv.allocate(), nil
		}
	} else {
		return nil, err
	}
}

func (d *Domain) Set(list []*Server) error {
	if err := d.Destory(); err != nil {
		return err
	}
	d.l.Lock()
	defer d.l.Unlock()
	if d.balancer == nil {
		return ErrNotBalancer
	}
	for _, v := range list {
		if err := d.balancer.Set(v); err != nil {
			return err
		}
	}
	return nil
}

func (d *Domain) Destory() error {
	d.l.Lock()
	defer d.l.Unlock()
	if d.balancer == nil {
		return ErrNotBalancer
	}
	return d.balancer.Destory()
}

type Server struct {
	domain *Domain
	ip     int32 // Little
	port   int16 // Little
	weight int32
	total  int32
	stat   Stat

	strIp   string
	intPort int

	l sync.RWMutex
}

func (s *Server) Ip() string {
	return s.strIp
}

func (s *Server) Port() int {
	return s.intPort
}

func (s *Server) LittleIp() int32 {
	return s.ip
}

func (s *Server) LittlePort() int16 {
	return s.port
}
