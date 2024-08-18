<div>
    {websiteName}
    <FiestaBoard />
</div>


<script>
    import FiestaBoard from "$lib/FiestaBoard.svelte";

    let websiteName = "Snake Fiesta"
</script>


<style>
    div {
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }
</style>