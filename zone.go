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
	Hosts map[string]net.IP
	ttl   uint32
}

func (h *zoneImpl) Records(q dns.Question) []dns.RR {
	if ip, ok := h.Hosts[q.Name]; ok {
		rr := make([]dns.RR, 0)
		if len(ip) == net.IPv6len {
			rr = append(rr, &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    h.ttl,
				},
				AAAA: ip.To16(),
			})
		}
		if len(ip) == net.IPv4len {
			rr = append(rr, &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    h.ttl,
				},
				A: ip.To4(),
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
	hosts := make(map[string]string)
	err = json.Unmarshal(confRaw, &hosts)
	if err != nil {
		panic(err)
	}
	for k, v := range hosts {
		if !strings.HasSuffix(k, ".") {
			k += "."
		}
		ip := net.ParseIP(v)
		if ip == nil {
			continue
		}
		if strings.Contains(v, ".") {
			ip = ip.To4()
		} else if strings.Contains(v, ":") {
			ip = ip.To16()
		} else {
			continue
		}
		impl.Hosts[k] = ip
	}
	impl.ttl = 120
	return impl
}
