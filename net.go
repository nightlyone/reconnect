package reconnect

import "net"

// Reconn Minimal wrapper around net.Conn that implements the Reconnectable interface
type Reconn struct {
	net.Conn
}

// Try to reconnect once and consider any network related error
// besides timeouts and temporary errors (as classified by net package) as final.
// To be used by the various retry handling functions in this package.
// NOTE: Doesn't handle close.
func (r *Reconn) Reconnect() error {
	addr := r.Conn.RemoteAddr()
	c, err := net.Dial(addr.Network(), addr.String())
	if err == nil {
		r = &Reconn{c}
		return nil
	}
	neterr, ok := err.(net.Error)
	if !ok {
		return FinalError{RealError: err}
	}

	if neterr.Timeout() || neterr.Temporary() {
		return err
	}

	return FinalError{RealError: err}
}
