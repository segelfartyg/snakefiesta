<script lang="ts">
  import { onMount } from "svelte";
  import { Direction } from "../enums/Direction";
  import type { PlayerIntentDirection } from "../interfaces/PlayerIntentDirection";
  import type { ServerPlayerTiles } from "../interfaces/ServerPlayerTiles";

  // CONSTANTS
  const BOARD_TOTAL_HEIGHT_TILE_SIZE = 500;
  const BOARD_TOTAL_WIDTH_TILE_SIZE = 1500;
  const FIELD_COLOR = "#ffffff";
  const GRID_COLOR = "#ffffff";
  const SNAKE_COLOR = "#206020";
  const CELL_SIZE = 10;
  const width = BOARD_TOTAL_WIDTH_TILE_SIZE;
  const height = BOARD_TOTAL_HEIGHT_TILE_SIZE;
  let image: any;
  export let occupiedTiles: ServerPlayerTiles[] = [];
  export let previousOccupiedTiles: ServerPlayerTiles[] = [];
  export let handlePlayerMovement;

  let canvas: any;
  let context: any;

  let PLAYER_WISH_MOVE = Direction.Up;

  // LOOP FOR SENDING NEW INTENT EVENTS TO SERVER
  setInterval(() => {
    let move: PlayerIntentDirection = {
      direction: PLAYER_WISH_MOVE,
    };
    handlePlayerMovement(move);
  }, 100);

  onMount(() => {
    context = canvas.getContext("2d");
    image =  document.getElementById("image");
    draw_field();
  });

  // LOOP FOR SETTING THE BOARD FROM TILE ARRAYS POPULATED FROM WS EVENTS
  setInterval(() => {
    previousOccupiedTiles.forEach((tile) => {
      context.fillStyle = FIELD_COLOR;
      context.strokeStyle = GRID_COLOR;
      context.fillRect(tile.x * 10, tile.y * 10, CELL_SIZE, CELL_SIZE);
      context.strokeRect(tile.x * 10, tile.y * 10, CELL_SIZE, CELL_SIZE);
    });
    occupiedTiles.forEach((tile) => {
      context.fillStyle = SNAKE_COLOR;
      context.strokeStyle = GRID_COLOR;
      context.textBaseLine ="middle"
      context.font = "10px Arial";
      context.drawImage(image, tile.x * 10, tile.y * 10, CELL_SIZE, CELL_SIZE);
      // context.fillRect(tile.x * 10, tile.y * 10, CELL_SIZE, CELL_SIZE);
      // context.strokeRect(tile.x * 10, tile.y * 10, CELL_SIZE, CELL_SIZE);
    });
  }, 10);

  const draw_field = function () {
    context.fillStyle = FIELD_COLOR;
    context.strokeStyle = FIELD_COLOR;
    context.fillRect(
      0,
      0,
      BOARD_TOTAL_WIDTH_TILE_SIZE,
      BOARD_TOTAL_HEIGHT_TILE_SIZE
    );
    
    for (let i = CELL_SIZE; i < BOARD_TOTAL_HEIGHT_TILE_SIZE; i += CELL_SIZE) {
      context.moveTo(0, i);
      context.lineTo(BOARD_TOTAL_WIDTH_TILE_SIZE, i);
      context.stroke();
    }
    for (let i = CELL_SIZE; i < BOARD_TOTAL_WIDTH_TILE_SIZE; i += CELL_SIZE) {
      context.moveTo(i, 0);
      context.lineTo(i, BOARD_TOTAL_HEIGHT_TILE_SIZE);
      context.stroke();
    }
  };

  function onKeyPress(event: { key: any }) {
    switch (event.key) {
      case "ArrowRight":
        PLAYER_WISH_MOVE = Direction.Right;
        break;
      case "ArrowDown":
        PLAYER_WISH_MOVE = Direction.Down;
        break;
      case "ArrowLeft":
        PLAYER_WISH_MOVE = Direction.Left;
        break;
      case "ArrowUp":
        PLAYER_WISH_MOVE = Direction.Up;
        break;
    }
  }
</script>

<div>
  <canvas {width} {height} bind:this={canvas} />
</div>

<img id="image" src="src\lib\favicon.png">
<svelte:window on:keydown={onKeyPress} />

<style>
  div {
    background: rgb(255, 255, 255);
    width: fit-content;
    display: grid;
    /* REMEMBER TO CHANGE THESE VALUES IF YOU WANT THE BOARD SIZE TO CHANGE */

    border: solid black 2px;
    margin: auto;
  }

  canvas {
    width: 1500px;
    height: 500px;
    border: solid green 1px;
  }

  #image {
    display:none;
  }
</style>
