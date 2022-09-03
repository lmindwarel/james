import { usePlayerStore } from "@/plugins/store/player";
import { PlayerStatus, WebsocketMessage } from "@/types";

// @ts-ignore
const address = `${import.meta.env.VITE_JAMES_API_ADDRESS.replace('http', 'ws')}/ws`
console.log(address)
const websocket = new WebSocket(address)

websocket.onmessage = function (event) {
    const websocketMessage = JSON.parse(event.data) as WebsocketMessage
    handleMessage(websocketMessage)
}

websocket.onopen = function () {
    console.log("Successfully connected to the websocket server.")
}

websocket.onerror = function (err) {
    console.error("failed to connect to the websocket server", err)
}

function handleMessage(message: WebsocketMessage){
    switch(message.topic){
        case "player-status":
            const playerStore = usePlayerStore()
            playerStore.updateFromPlayerStatus(message.data as PlayerStatus)

        default:
            console.warn("Unknown message topic: " + message.topic)
    }
}