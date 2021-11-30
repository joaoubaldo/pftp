package pftp

type result struct {
	code int
	msg  string
	err  error
	log  *logger
}

func (r *result) Response(handler *clientHandler) error {
	if r.log != nil && r.err != nil {
		handler.eventC.Send(Event{name: ErrorEventType, payload: ErrorEvent{RemoteAddr: handler.conn.RemoteAddr().String(), ErrorMessage: r.err.Error()}})
		r.log.err("command error response: %s", r.err)
	}

	if r.code != 0 {
		return handler.writeMessage(r.code, r.msg)
	}
	return nil
}
