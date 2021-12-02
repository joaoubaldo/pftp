package pftp

import (
	"reflect"
	"strings"
	"testing"
)

func Test_EventChan_SendNoReceiver(t *testing.T) {
	eventC := NewEventChan(0)
	err := eventC.Send(Event{name: ClientCommandEventType, payload: ClientCommandEvent{}})
	if err == nil {
		t.Errorf("expecting Send to fail")
	} else {
		if !strings.Contains(err.Error(), "no receiver") {
			t.Errorf("should error with no receiver")
		}
	}
}

func Test_EventChan_Send(t *testing.T) {
	eventC := NewEventChan(1)
	err := eventC.Send(Event{name: ClientCommandEventType, payload: ClientCommandEvent{RemoteAddr: "1.1.1.1"}})
	if err != nil {
		t.Errorf("expecting Send not to fail")
	} else {
		select {
		case event := <-eventC:
			event.Name()
			switch event.Name() {
			case ClientCommandEventType:
				if !reflect.DeepEqual(event.Payload().(ClientCommandEvent), ClientCommandEvent{RemoteAddr: "1.1.1.1"}) {
					t.Errorf("expected event payload does not match")
				}
			default:
				t.Errorf("expected event does not match")
			}

		default:
			t.Errorf("event expected")
		}
	}
}
