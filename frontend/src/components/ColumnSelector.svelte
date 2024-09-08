<section>
    <Button on:click={toggleDropdown} bind:clientWidth={buttonWidth}>{label}</Button>
    <div class="dropdown" bind:this={dropdown} style="min-width: {buttonWidth}px">
        {#each options as option}
            <div class="option">
                <input type="checkbox" checked={selected.includes(option)} on:change={() => handleChange(option)} />
                <span>{option}</span>
            </div>
        {/each}
    </div>
</section>

<style>
    .dropdown {
        display: none;
        position: absolute;
        margin-top: 6px;
        padding: 5px 10px;
        background-color: var(--background);
        border-radius: 4px;
        border: 1px solid var(--primary);
        box-shadow: 0 14px 14px rgba(0, 0, 0, 0.25);
    }

    .option {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 4px;
        font-size: 18px;
    }

    input {
        height: 16px;
        width: 16px;
    }
</style>

<script>
    import Button from "./Button.svelte";
    import { createEventDispatcher } from 'svelte';

    export let label = "Select Columns";
    export let options = [];
    export let selected = [];

    const dispatch = createEventDispatcher();

    function handleChange(option) {
        if (!selected) {
            selected = [];
        }

        if (selected.includes(option)) {
            selected = selected.filter((item) => item !== option);
        } else {
            selected = [...selected, option];
        }

        dispatch('change', selected)
    }

    let dropdown;
    let buttonWidth;

    function toggleDropdown() {
        if (dropdown.style.display === 'block') {
            dropdown.style.display = 'none';
        } else {
            dropdown.style.display = 'block';
        }
    }

    document.addEventListener('click', (e) => {
        let current = e.target;
        let dropdownFound = false;

        while (current) {
            if (current.attributes) {
                if (current.hasAttribute('istrigger')) {
                    dropdownFound = true;
                    break;
                }
            }

            if (current === dropdown) {
                dropdownFound = true;
                break;
            } else {
                current = current.parentNode;
            }
        }

        if (!dropdownFound) {
            dropdown.style.display = 'none';
        }
    });
</script>