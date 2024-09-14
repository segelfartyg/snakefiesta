<div style={styleObject}></div>

<style>
 div{
    background: white;
    height:10px;
    width:10px;
    border:solid rgba(0, 0, 0, 0.027) 1px;
 }

</style>

<script>
    export let color;
    export let id;
    let styleObject = "background: " + color;
</script>