<script lang="ts">
  import { Direction } from "../enums/Direction";
  import type { PlayerIntentPosition } from "../interfaces/PlayerIntentPosition";
  import FiestaTile from "./FiestaTile.svelte";

  // CONSTANTS
  const BOARD_TOTAL_HEIGHT_TILE_SIZE = 50;
  const BOARD_TOTAL_WIDTH_TILE_SIZE = 50;

  export let fiestaChunk = 1;
  export let occupiedTiles = {
    x: 1,
    y: 1,
  };
  export let handlePlayerMovement;

  // WHERE THE HEAD PLAYER TILE IS LOCATED IN Y
  let PLAYER_TILE_Y = 15;
  // WHERE THE HEAD PLAYER TILE IS LOCATED IN X
  let PLAYER_TILE_X = 15;

  let PLAYER_WISH_MOVE = Direction.Down;

  // 2 DIMENSIONAL ARRAY FOR X/Y AXIS OF THE FIESTA
  let FIESTA_TILES: FiestaTileUI[][] = [];

  // POPULATES THE ABOVE
  let FIESTA_ROW: FiestaTileUI[] = [];

  // INITIALIZES THE BOARD, THIS IS NOT NEEDED (mostly) DURING GAME
  for (let i = 0; i < BOARD_TOTAL_HEIGHT_TILE_SIZE; i++) {
    FIESTA_ROW = [];

    for (let j = 0; j < BOARD_TOTAL_WIDTH_TILE_SIZE; j++) {
      if (PLAYER_TILE_Y == i && PLAYER_TILE_X == j) {
        FIESTA_ROW.push({
          id: "x" + i.toString() + "y" + j.toString(),
          color: "black",
        });
      } else {
        FIESTA_ROW.push({
          id: "x" + i.toString() + "y" + j.toString(),
          color: "white",
        });
      }
    }

    FIESTA_TILES.push(FIESTA_ROW);
  }

  setInterval(() => {
    let WISH_MOVEMENT_X = PLAYER_TILE_X;
    let WISH_MOVEMENT_Y = PLAYER_TILE_Y;

    switch (PLAYER_WISH_MOVE) {
      case Direction.Right:
        WISH_MOVEMENT_X = PLAYER_TILE_X + 1;
        WISH_MOVEMENT_Y = PLAYER_TILE_Y;
        break;
      case Direction.Down:
        WISH_MOVEMENT_X = PLAYER_TILE_X;
        WISH_MOVEMENT_Y = PLAYER_TILE_Y + 1;
        break;
      case Direction.Left:
        WISH_MOVEMENT_X = PLAYER_TILE_X - 1;
        WISH_MOVEMENT_Y = PLAYER_TILE_Y;
        break;
      case Direction.Up:
        WISH_MOVEMENT_X = PLAYER_TILE_X;
        WISH_MOVEMENT_Y = PLAYER_TILE_Y - 1;
        break;
      default:
        WISH_MOVEMENT_X = PLAYER_TILE_X;
        WISH_MOVEMENT_Y = PLAYER_TILE_Y;
    }

    // MAKING BAKDEL WHITE
    FIESTA_TILES[PLAYER_TILE_Y][PLAYER_TILE_X] = {
      id: "x" + PLAYER_TILE_X.toString() + "y" + PLAYER_TILE_Y.toString(),
      color: "white",
    };

    // CHECKING X AXIS FIRST
    if (WISH_MOVEMENT_X >= BOARD_TOTAL_WIDTH_TILE_SIZE) {
      WISH_MOVEMENT_X = 0;
    } else if (WISH_MOVEMENT_X <= 0 - 1) {
      WISH_MOVEMENT_X = BOARD_TOTAL_WIDTH_TILE_SIZE - 1;
    }

    // CHECKING Y AXIS SECOND
    if (WISH_MOVEMENT_Y >= BOARD_TOTAL_HEIGHT_TILE_SIZE) {
      WISH_MOVEMENT_Y = 0;
    } else if (WISH_MOVEMENT_Y <= 0 - 1) {
      WISH_MOVEMENT_Y = BOARD_TOTAL_HEIGHT_TILE_SIZE - 1;
    }

    console.log(occupiedTiles);
    FIESTA_TILES[occupiedTiles.y][occupiedTiles.x] = {
      id: "x" + occupiedTiles.x.toString() + "y" + occupiedTiles.y.toString(),
      color: "pink",
    };

    // MAKING THE MOVE, DEPENDING ON THE PREVIOUS MANIPULATION OF THE WISHED STATE
    // FIESTA_TILES[WISH_MOVEMENT_Y][WISH_MOVEMENT_X] = {
    //   id: "x" + WISH_MOVEMENT_X.toString() + "y" + WISH_MOVEMENT_Y.toString(),
    //   color: "green",
    // };

    // SETTING CURRENT TILE POSITIONS
    PLAYER_TILE_X = WISH_MOVEMENT_X;
    PLAYER_TILE_Y = WISH_MOVEMENT_Y;

    // console.log("POS: ", PLAYER_TILE_X, PLAYER_TILE_Y);
    let move: PlayerIntentPosition = {
      x: PLAYER_TILE_X,
      y: PLAYER_TILE_Y,
      direction: PLAYER_WISH_MOVE
    };

    $: handlePlayerMovement(move);
  }, 100);

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
  {#each FIESTA_TILES as item, index (item)}
    {#each FIESTA_TILES[index] as cell (cell)}
      <FiestaTile id={cell.id} color={cell.color} />
    {/each}
  {/each}
</div>
<svelte:window on:keydown={onKeyPress} />

<style>
  div {
    background: blue;
    height: 500px;
    width: fit-content;
    display: grid;
    /* REMEMBER TO CHANGE THESE VALUES IF YOU WANT THE BOARD SIZE TO CHANGE */
    grid-template-columns: repeat(50, 1fr);
    grid-template-rows: repeat(50, 1fr);
    margin: auto;
  }
</style>
