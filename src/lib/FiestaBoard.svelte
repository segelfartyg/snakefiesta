<script  lang="ts">
  import FiestaTile from "./FiestaTile.svelte";

  
  let BOARD_TOTAL_HEIGHT_TILE_SIZE  = 50;
  let BOARD_TOTAL_WIDTH_TILE_SIZE = 100;



interface FiestaTileUI {id: String, color: String};

  let BOARD_TILES: FiestaTileUI[][] = [];
  let PLAYER_TILE_H = 6;
  let PLAYER_TILE_W = 6;



  let BOARD_ROW: FiestaTileUI[] = [];
  

for(let i = 0; i< BOARD_TOTAL_HEIGHT_TILE_SIZE; i++){
    BOARD_ROW = []

for(let j = 0; j < BOARD_TOTAL_WIDTH_TILE_SIZE; j++){

if(PLAYER_TILE_H == i && PLAYER_TILE_W == j){
    BOARD_ROW.push({id:"w"+ i.toString() + "h"+ j.toString(), color: "black"})
}
else{
    BOARD_ROW.push({id:"w"+ i.toString() + "h"+ j.toString(), color: "white"})
}


}

BOARD_TILES.push(BOARD_ROW);

}


  setInterval(() => {
PLAYER_TILE_W++;
    BOARD_TILES[6][PLAYER_TILE_W] = {id:"w"+ PLAYER_TILE_W.toString() + "h"+ PLAYER_TILE_H.toString(), color: "black"}

BOARD_TILES = BOARD_TILES;

}, 100);

</script>
<div>
    {#each BOARD_TILES as item, index (item)}
    {#each BOARD_TILES[index] as cell (cell)}
        <FiestaTile id={cell.id} color={cell.color}/>
    {/each}
  {/each}
</div>



<style>
    div {
        background: blue;
        height: 500px;
        width:1000px;
        display: grid;
        grid-template-columns: repeat(100, 1fr);
        grid-template-rows: repeat(50, 1fr);
        margin:auto;
    }
</style>