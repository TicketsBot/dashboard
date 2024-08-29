<div class="root">
    <ul>
        <slot/>
    </ul>
</div>

<style>
    .root {
        display: flex;
        flex-direction: row;
    }

    ul {
        width: 100%;
        list-style-type: none;
        margin: 2px 0 0 55px;
        padding: 0;
    }
</style>