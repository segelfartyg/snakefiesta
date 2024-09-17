<script lang="ts">
  import FiestaBoard from "$lib/FiestaBoard.svelte";
  import { SERVER_URI } from "../consts/consts";
  import { onMount } from "svelte";
  import type { PlayerIntentDirection } from "../interfaces/PlayerIntentDirection";
  import type { PlayerMovementIntent } from "../interfaces/PlayerMovementIntent";
  import type { ServerPlayerTiles } from "../interfaces/ServerPlayerTiles";
  let websiteName = "Snake Fiesta";
  let fiestaChunk = 1;
  let occupiedTiles: ServerPlayerTiles[] = [];
  let previousOccupiedTiles: ServerPlayerTiles[] = [];
  let username= "";
  let webSocket: WebSocket;

  onMount(() => {
    webSocket = new WebSocket("ws://" + SERVER_URI + "/ws");

    webSocket.onclose = function (evt) {
      console.log("SOCKET CONNECTION CLOSED");
    };

    webSocket.onmessage = (msg) => {
      console.log("RECEIVED MESSAGE");
      console.log(msg.data)
      let jsonRes = JSON.parse(msg.data);
      console.log(Object.keys(jsonRes))
      console.log(jsonRes)
      previousOccupiedTiles = occupiedTiles
      occupiedTiles = []
      Object.keys(jsonRes).forEach(playerId => {
        console.log(jsonRes[playerId].Chunk);
       occupiedTiles.push({
          playerId: jsonRes[playerId].playerId,
          x: jsonRes[playerId].X,
          y: jsonRes[playerId].Y
        })
      });
    };
  });

  function handlePlayerMovement(positionIntent: PlayerIntentDirection) {
    let intent: PlayerMovementIntent = {
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
  <FiestaBoard {occupiedTiles} {previousOccupiedTiles} {handlePlayerMovement} />
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
