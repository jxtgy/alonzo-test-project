package l5

import (
	"os"
	"time"
)

var (
	statErrorReportInterval = time.Second
	statReportInterval      = 5 * time.Second
	statMaxErrorCount       = 16
	statMaxErrorRate        = 0.2
)

type Stat struct {
	AllocCount uint32
	SuccCount  uint32
	SuccDelay  uint64
	ErrCount   uint32
	ErrDelay   uint64
	LastReport time.Time
}

func (s *Server) needReport() bool {
	now := time.Now()
	interval := statReportInterval
	if s.stat.ErrCount > 0 {
		interval = statErrorReportInterval
	}
	if s.stat.LastReport.Add(interval).After(now) {
		return true
	}
	if s.stat.ErrCount >= uint32(statMaxErrorCount) {
		return true
	}
	if float64(s.stat.AllocCount)*statMaxErrorRate <= float64(s.stat.ErrCount) {
		return true
	}
	return false
}

func (s *Server) report(args ...bool) error {
	force := false
	if len(args) > 0 {
		force = args[0]
	}
	s.l.RLock()
	if !force && !s.needReport() {
		s.l.RUnlock()
		return nil
	} else {
		s.l.RUnlock()
	}
	s.l.Lock()
	defer s.l.Unlock()
	if s.domain != nil {
		var err error
		_, err = dial(QOS_CMD_GET_STAT, 0, s.domain.mod, s.domain.cmd, s.ip,
			uint32(s.port), s.stat.AllocCount, int32(os.Getpid()))
		if err != nil {
			return err
		}
		_, err = dial(QOS_CMD_CALLER_UPDATE_BIT64, 0, int32(0), int32(0),
			s.domain.mod, s.domain.cmd, s.ip, uint32(s.port), int32(-1),
			s.stat.ErrCount, s.stat.ErrDelay, int32(os.Getpid()))
		if err != nil {
			return err
		}
		_, err = dial(QOS_CMD_CALLER_UPDATE_BIT64, 0, int32(0), int32(0),
			s.domain.mod, s.domain.cmd, s.ip, uint32(s.port), int32(0),
			s.stat.SuccCount, s.stat.SuccDelay, int32(os.Getpid()))
		if err != nil {
			return err
		}
	}
	s.stat.AllocCount = 0
	s.stat.SuccCount = 0
	s.stat.SuccDelay = 0
	s.stat.ErrCount = 0
	s.stat.ErrDelay = 0
	s.stat.LastReport = time.Now()
	return nil
}

func (s *Server) StatUpdate(result int32, usetime uint64) error {
	s.l.Lock()
	defer s.l.Unlock()
	if result >= 0 {
		s.stat.SuccCount++
		s.stat.SuccDelay += usetime
	} else {
		s.stat.ErrCount++
		s.stat.ErrDelay += usetime
	}
	return nil
}

func (s *Server) allocate() *Server {
	s.l.Lock()
	defer s.l.Unlock()
	s.stat.AllocCount++
	return s
}
