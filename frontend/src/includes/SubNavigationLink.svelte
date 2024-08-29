<li>
    <div class="wrapper">
        <a href="{href}" class:active>
            {#if icon}
                <i class="fas {icon}"/>
            {/if}

            <span>
            <slot/>
        </span>
        </a>
    </div>
</li>

<style>
    li {
        /*padding: 2px 0;*/
    }

    .wrapper {
        width: 100%;
    }

    a.active {
        font-weight: 600 !important;
    }

    a:not(.active) {
        opacity: 0.75;
    }

    a, a:link, a:hover, a:visited, a:active {
        color: inherit;
        text-decoration: none;
    }

    i {
        width: 20px;
        text-align: center;
    }
</style>

<script>
    export let currentRoute;
    export let icon;
    export let href = "#";
    export let routePrefix;

    let active = href !== "/" && ((routePrefix || href).toLowerCase() === currentRoute.name.toLowerCase() ||
        currentRoute.name.toLowerCase().startsWith((routePrefix || href).toLowerCase()));

    $: active;
</script>
