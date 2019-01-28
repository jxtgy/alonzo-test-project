package l5

import "testing"

func TestApiGetSid_And_ApiGetRoute(t *testing.T) {
	name := "prot_proxy"
	domain, err := ApiGetSid(name)
	if err != nil {
		t.Error("ApiGetSid Errof:%s", err.Error())
		return
	} else {
		t.Logf("ApiGetSid(%s) Mod:%d Cmd:%d", name, domain.Mod, domain.Cmd)
	}
	srv, err := ApiGetRoute(domain.mod, domain.cmd)
	if err != nil {
		t.Error("ApiGetRoute Errof:%s", err.Error())
		return
	} else {
		t.Logf("ApiGetRoute(%s) Ip:%s Port:%d", name, Ip2String(srv.ip), srv.port)
	}
}

func TestApiGetRouteBySid_And_ApiRouteResultUpdate(t *testing.T) {
	name := "prot_proxy"
	srv, err := ApiGetRouteBySid(name)
	if err != nil {
		t.Error("ApiGetRouteBySid Errof:%s", err.Error())
		return
	}
	err = ApiRouteResultUpdate(srv, 0, 0)
	if err != nil {
		t.Error("ApiRouteResultUpdate Errof:%s", err.Error())
		return
	}
}
