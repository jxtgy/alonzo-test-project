package l5

import "testing"

func TestWeightedRoundRobin(t *testing.T) {
	lb := NewBalancer(CL5_LB_TYPE_WRR)
	srv, err := lb.Get()
	if err == nil {
		t.Errorf("empty lb should return error!")
		return
	} else {
		t.Logf("info:%s", err.Error())
	}

	weights := []int32{0, 200, 300}
	for i := 2; i >= 0; i-- {
		err = lb.Set(&Server{ip: int32(i), port: int16(i), weight: weights[i], total: 0})
		if err != nil {
			t.Errorf("lb Set Server error:%s", err.Error())
			return
		}
	}

	err = lb.Remove(&Server{ip: 2, port: 2})
	if err != nil {
		t.Errorf("lb Remove Server error:%s", err.Error())
		return
	}

	for i := 0; i < 2; i++ {
		srv, err = lb.Get()
		if err != nil {
			t.Errorf("lb Get error: %s!", err.Error())
			return
		} else {
			t.Logf("lb Get srv %s:%d", Ip2String(srv.ip), srv.port)
		}

		if srv.ip != int32(i) || srv.port != int16(i) {
			t.Errorf("lb has error!")
			return
		}
	}
}
