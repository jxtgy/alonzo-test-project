package l5

import "testing"

func TestDomain_Get(t *testing.T) {
	d, err := newTestDomain()
	if err != nil {
		t.Errorf("Domains.Query Error:%s", err.Error())
		return
	} else {
		t.Logf("Domains.Query(%s) Mod:%d, Cmd:%d", d.Name, d.Mod, d.Cmd)
	}

	var srv *Server
	srv, err = d.Get()
	if err != nil {
		t.Errorf("GetServer(%s) Error:%s", d.Name, err.Error())
		return
	} else {
		t.Logf("GetServer(%s) Success Ip:%s Port:%d", d.Name, Ip2String(srv.ip), srv.port)
	}
}
