package pool

import (
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"sort"
)

type IdleConns []*Conn

func (c *IdleConns) Pop(sig esig.ESIG, maxOpenEndpoints int) *Endpoint {

	type finder struct {
		connIdx     int
		endpointIdx int
		order       int
		cap         int
		usedAt      int64
	}

	var finders []finder
	var findConnIdx int
	var findEndpoint *Endpoint

	conns := *c

	for idx, conn := range conns {

		if len(conn.endpoints) == 0 {
			finders = append(finders, finder{idx, -1, 0, 0, conn.UsedAt().Unix()})
			continue
		}

		capEnd := len(conn.endpoints)

		for i, endpoint := range conn.endpoints {

			if esig.Equal(endpoint.sig, sig) {
				findEndpoint = endpoint
				findConnIdx = idx
				break
			}

			orderByte := 2

			if esig.HighBoundEqual(endpoint.sig, sig.High()) && uuid.Equal(endpoint.sig.Low(), uuid.Nil) {
				orderByte = 1
			}

			finders = append(finders, finder{idx, i, orderByte, capEnd, endpoint.UsedAt().Unix()})

		}

		if findEndpoint != nil {
			break
		}

	}

	if findEndpoint != nil {
		c.remove(findConnIdx)
		return findEndpoint
	}

	if len(finders) == 0 {
		return nil
	}

	sort.Slice(finders, func(i, j int) bool {
		if finders[i].order < finders[j].order {
			return true
		}
		if finders[i].order > finders[j].order {
			return false
		}

		if finders[i].cap < finders[j].cap {
			return true
		}
		if finders[i].cap > finders[j].cap {
			return false
		}

		return finders[i].usedAt < finders[j].usedAt
	})

	f := finders[0]

	conn := conns[f.connIdx]

	switch f.order {

	case 0:
		findEndpoint = &Endpoint{
			conn: conn,
		}
	case 1:
		findEndpoint = conn.endpoints[f.endpointIdx]
	case 2:

		if len(conn.endpoints) < maxOpenEndpoints {
			findEndpoint = &Endpoint{
				conn: conn,
			}
		} else {
			findEndpoint = conn.endpoints[f.endpointIdx]
		}

	}

	c.remove(f.connIdx)
	return findEndpoint
}

func (c *IdleConns) remove(i int) {

	conns := *c
	conns[i] = conns[len(conns)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	*c = conns[:len(conns)-1]

}
