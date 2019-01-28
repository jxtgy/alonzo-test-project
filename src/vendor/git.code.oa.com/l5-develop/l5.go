// Copyright (c) 2017, Tencent. All rights reserved.
// Authors esliznwang cheaterlin

package l5

func ApiGetSid(sid string) (*Domain, error) {
	return domainss.Query(sid)
}

func ApiGetRouteBySid(sid string) (*Server, error) {
	domain, err := domainss.Query(sid)
	if err != nil {
		return nil, err
	}
	return domain.Get()
}

func ApiGetRoute(mod int32, cmd int32) (*Server, error) {
	return anonymouss.Get(mod, cmd).Get()
}

func ApiRouteResultUpdate(s *Server, result int32, usetime uint64) error {
	if s == nil {
		return nil
	}

	return s.StatUpdate(result, usetime)
}
