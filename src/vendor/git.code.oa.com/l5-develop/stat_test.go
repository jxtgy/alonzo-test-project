package l5

import "testing"

func TestServer_StatUpdate(t *testing.T) {
	d, err := newTestDomain()
	if err != nil {
		t.Errorf("Domains.Query Error:%s", err.Error())
		return
	}
	srv, err := d.Get()
	if err != nil {
		t.Errorf("GetServer(%s) Error:%s", d.Name, err.Error())
		return
	}
	succCount := srv.stat.SuccCount
	srv.StatUpdate(0, 1000)
	if succCount+1 != srv.stat.SuccCount {
		t.Errorf("Server.StatUpdate SuccCount: %d != %d", succCount+1, srv.stat.SuccCount)
	}
}

func TestServer_Allocate(t *testing.T) {
	d, err := newTestDomain()
	if err != nil {
		t.Errorf("Domains.Query Error:%s", err.Error())
		return
	}
	srv, err := d.Get()
	if err != nil {
		t.Errorf("GetServer(%s) Error:%s", d.Name, err.Error())
		return
	}
	alloc := srv.stat.AllocCount
	if srv = srv.allocate(); srv.stat.AllocCount != alloc+1 {
		t.Errorf("AllocCount: %d != %d", alloc+1, srv.stat.AllocCount)
	}
}

func TestServer_Report(t *testing.T) {
	d, err := newTestDomain()
	if err != nil {
		t.Errorf("Domains.Query Error:%s", err.Error())
		return
	}
	srv, err := d.Get()
	if err != nil {
		t.Errorf("GetServer(%s) Error:%s", d.Name, err.Error())
		return
	}
	if err := srv.report(true); err != nil {
		t.Errorf("Server.Report fail: %s", err.Error())
		return
	}
	if srv.stat.AllocCount != 0 {
		t.Errorf("Server.Report fail: counter：%d not empty", srv.stat.AllocCount)
	}
}
