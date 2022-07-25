class EventBus {
  constructor () {
    this.events = {}
  }

  on (eventName, fn) {
    this.events[eventName] = this.events[eventName] || []
    this.events[eventName].push(fn)
  }

  off (eventName, fn) {
    if (this.events[eventName]) {
      for (let i = 0; i < this.events[eventName].length; i++) {
        if (this.events[eventName][i] === fn) {
          this.events[eventName].splice(i, 1)
          break
        }
      };
    }
  }

  emit (eventName, data) {
    if (this.events[eventName]) {
      this.events[eventName].forEach(function (fn) {
        fn(data)
      })
    }
  }

  notifySuccess (message:string, timeout:number = 3000) {
    this.emit('notify', { message, timeout, color: 'success' })
  }

  notifyError (message:string, timeout:number = 3000) {
    this.emit('notify', { message, timeout, color: 'error' })
  }

  unhandledAPIError (err) {
    this.emit('unhandled-api-error', err)
  }
}

export default new EventBus()
