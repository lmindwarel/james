package controller

const (
	EventJamesStatusChange = "change:james-status"
)

func (ctrl *Controller) AddEventListener(event string, handler func(data interface{})) {
	log.Debugf("add event listener for event: %s", event)
	ctrl.listeners[event] = append(ctrl.listeners[event], handler)
}

func (ctrl *Controller) triggerListeners(event string, data interface{}) {
	eventListeners, hasEventListeners := ctrl.listeners[event]
	if !hasEventListeners {
		log.Warningf("No listener for event %s", event)
		return
	}

	for _, l := range eventListeners {
		l(data)
	}
}
