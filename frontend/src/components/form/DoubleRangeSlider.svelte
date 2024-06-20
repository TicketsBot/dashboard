<section>
    <span class="form-label">{label}</span>
    <div class="line" bind:clientWidth={width}>
        <div class="slider" bind:this={leftParent} on:mousedown={onMouseDown} on:touchstart={onMouseDown} style="
            left: min({rightParent?.style?.left || `${width}px`}, max(0px, min({width-(sliderDiameter/2)}px, {leftOffset}px)));
            ">
            <div bind:this={leftSlider} style="height: {sliderDiameter}px; width: {sliderDiameter}px">
                <span class="label" bind:this={leftLabel}>{start}</span>
            </div>
        </div>
        <div class="slider" bind:this={rightParent} on:mousedown={onMouseDown} on:touchstart={onMouseDown} style="
            left: max({leftParent?.style?.left || `${width}px`}, max(0px, min({width-(sliderDiameter/2)}px, {rightOffset}px)));
            ">
            <div bind:this={rightSlider} style="height: {sliderDiameter}px; width: {sliderDiameter}px">
                <span class="label" bind:this={rightLabel}>{end}</span>
            </div>
        </div>
    </div>
</section>

<svelte:window
    on:mouseup={onMouseUp}
    on:touchend={onMouseUp}
    on:mousemove={onMouseMove}
    on:touchmove={onTouchMove}
/>

<style>
    section {
        display: flex;
        flex-direction: column;
        gap: 20px;
        width: 100%;
    }

    .line {
        height: 5px;
        border-radius: 5px;
        background-color: #3472f7;
        width: 100%;
        margin-bottom: 12px;
    }

    .slider {
        position: absolute;
    }

    .label {
        position: relative;
        top: -25px;
        left: 0;
        color: #9a9a9a;
        font-size: 12px;
        user-select: none;
    }

    .slider > div {
        position: relative;
        top: -8px;
        left: 0;
        border-radius: 50%;
        background-color: white;
        user-select: none;
        text-align: center;
    }
</style>

<script>
    import {onMount} from "svelte";

    export let label = "Slider";

    export let min = 0;
    export let max = 100;

    export let start;
    export let end;

    const sliderDiameter = 20;

    let leftSlider, rightSlider;
    let leftParent, rightParent;
    let leftLabel, rightLabel;
    let moving;

    let prevWidth = -1;
    let width;
    let leftOffset = 0;
    let rightOffset = 0;

    $: {
        if (prevWidth !== width) {
            leftOffset = (width - (sliderDiameter / 2)) * ((start - min) / (max - min));
            rightOffset = (width - (sliderDiameter / 2)) * ((end - min) / (max - min));
        }

        prevWidth = width;
    }

    function onMouseDown(e) {
        if (e.target === rightSlider || e.target === rightParent) {
            moving = rightSlider;
        } else if (e.target === leftSlider || e.target === leftParent) {
            moving = leftSlider;
        }
    }

    let previousTouch;

    function onTouchMove(e) {
        const touch = e.touches[0];

        if (previousTouch && moving) {
            e.movementX = touch.pageX - previousTouch.pageX;
            onMouseMove(e);
        }

        previousTouch = touch;
    }

    function onMouseMove(e) {
        if (moving === rightSlider) {
            rightOffset += e.movementX;

            const ratio = parseOffset(rightParent.style.left) / (width - (sliderDiameter / 2));
            end = Math.ceil(ratio * (max - min) + min);
        } else if (moving === leftSlider) {
            leftOffset += e.movementX;

            const ratio = parseOffset(leftParent.style.left) / (width - (sliderDiameter / 2));
            start = Math.ceil(ratio * (max - min) + min);
        }
    }

    function onMouseUp() {
        moving = null;
        previousTouch = null;
    }

    // calc(123px) -> 123
    function parseOffset(offset) {
        const regex = /calc\((\d+).*\)/;
        const match = offset.match(regex);
        if (match) {
            return parseInt(match[1]);
        }
    }

    onMount(() => {
        if (!start) start = min;
        if (!end) end = max;

        leftOffset = (width - (sliderDiameter / 2)) * ((start - min) / (max - min));
        rightOffset = (width - (sliderDiameter / 2)) * ((end - min) / (max - min));
    });
</script>