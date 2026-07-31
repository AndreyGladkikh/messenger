package login

import "net/netip"

type Command struct {
	Login string
	Password string
	UserAgent string
	IP netip.Addr
}