package main

import (
	"encoding/json"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/hashicorp/mdns"
	"github.com/miekg/dns"
)

type zoneImpl struct {
	Hosts map[string]net.IP `json:"hosts"`
	ttl   uint32
}

func (h *zoneImpl) Records(q dns.Question) []dns.RR {
	if ip, ok := h.Hosts[q.Name]; ok {
		rr := make([]dns.RR, 0)
		if ip6 := ip.To16(); ip6 != nil {
			rr = append(rr, &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   "zhimiao mdns server",
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    h.ttl,
				},
				AAAA: ip6,
			})
		}
		if ip4 := ip.To4(); ip4 != nil {
			rr = append(rr, &dns.A{
				Hdr: dns.RR_Header{
					Name:   "zhimiao mdns server",
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    h.ttl,
				},
				A: ip4,
			})
		}
		slog.Info("answer", "rr", rr)
		return rr
	}
	return nil
}

func NewZone() mdns.Zone {
	confRaw, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	impl := &zoneImpl{
		Hosts: map[string]net.IP{},
	}
	err = json.Unmarshal(confRaw, &impl.Hosts)
	if err != nil {
		panic(err)
	}
	for k, v := range impl.Hosts {
		if !strings.HasSuffix(k, ".") {
			impl.Hosts[k+"."] = v
		}
	}
	impl.ttl = 120
	return impl
}
