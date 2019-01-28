package l5

import "testing"

func newTestDomain() (*Domain, error) {
	return domainss.Query("prot_proxy")
}

func TestDomains_Query(t *testing.T) {
	d, err := newTestDomain()
	if err != nil {
		t.Errorf("GetDomain Error:%s", err.Error())
		return
	} else {
		t.Logf("GetDomain(%s) Success Mod: %d Cmd:%d", d.Name, d.Mod, d.Cmd)
	}

	if srv, err := d.Get(); err != nil {
		t.Errorf("GetServer(%s) Error:%s", d.Name, err.Error())
		return
	} else {
		t.Logf("GetServer(%s) Success Ip:%s, Port: %d", d.Name, Ip2String(srv.ip), srv.port)
	}
}
