import { usePlayerStore } from "@/plugins/store/player";
import { PlayerStatus, QueuedTrack, WebsocketMessage } from "@/types";
import eventbus from '@/services/eventbus'

// @ts-ignore
const address = `${import.meta.env.VITE_JAMES_API_ADDRESS.replace('http', 'ws')}/ws`

function connect() {
    const websocket = new WebSocket(address)

    websocket.onmessage = function (event) {
        const websocketMessage = JSON.parse(event.data) as WebsocketMessage
        handleMessage(websocketMessage)
    }

    websocket.onopen = function () {
        console.log("Successfully connected to the websocket server.")
    }

    websocket.onclose = function (e) {
        console.log('Websocket disconnected. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(connect, 1000);
    };

    websocket.onerror = function (err) {
        console.error("failed to connect to the websocket server", err)
    }
}

connect()

function handleMessage(message: WebsocketMessage) {
    const playerStore = usePlayerStore()
    // console.log("received message", message)
    switch (message.topic) {
        case "player-status":
            playerStore.updateFromPlayerStatus(message.data as PlayerStatus)
            break

        case "player-queue":
            playerStore.queue = message.data || [] as QueuedTrack[]
            break

        default:
            console.warn("Unknown message topic: " + message.topic)
    }

    eventbus.emit(`ws-${message.topic}`, message.data)
}