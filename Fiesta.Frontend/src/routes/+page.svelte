<script lang="ts">
  import FiestaBoard from "$lib/FiestaBoard.svelte";
  import { SERVER_URI } from "../consts/consts";
  import { onMount } from "svelte";
  import type { PlayerIntentPosition } from "../interfaces/PlayerIntentPosition";
  import type { PlayerMovementIntent } from "../interfaces/PlayerMovementIntent";
  let websiteName = "Snake Fiesta";
  let fiestaChunk = 1;
  let occupiedTiles = { x: 1, y: 1 };
  let username = "";
  let webSocket: WebSocket;

  onMount(() => {
    webSocket = new WebSocket("ws://" + SERVER_URI + "/ws");

    webSocket.onclose = function (evt) {
      console.log("SOCKET CONNECTION CLOSED");
    };

    webSocket.onmessage = (msg) => {
      console.log("RECEIVED MESSAGE");
      let jsonRes = JSON.parse(msg.data);

      if (jsonRes[username] != null || jsonRes[username] != undefined) {
        occupiedTiles = {
          x: jsonRes[username].X,
          y: jsonRes[username].Y,
        };
      }
    };
  });

  function handlePlayerMovement(positionIntent: PlayerIntentPosition) {
    let intent: PlayerMovementIntent = {
      x: positionIntent.x,
      y: positionIntent.y,
      playerId: username,
      timestamp: "utcnow",
      direction: positionIntent.direction 
    };

    if (webSocket != null) {
      webSocket.send(JSON.stringify(intent));
    }
  }
</script>

<div>
  {websiteName}
  <FiestaBoard {fiestaChunk} {occupiedTiles} {handlePlayerMovement} />
  <input bind:value={username} placeholder="enter your name" />
</div>

<style>
  div {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  input {
    position: absolute;
    top: 0;
  }
</style>
