import { readable } from "svelte/store"
import { SERVER_URI } from "../consts/consts";


let socket; 


export const socketStore = readable(null, (set) => {

    let webSocket = new WebSocket("ws://" + SERVER_URI + "/ws");

    webSocket.onclose = function (evt) {
      console.log("SOCKET CONNECTION CLOSED");
    };

    webSocket.addEventListener("open", (event) => {
        console.log("CONNECTION OPENED");
    });

})