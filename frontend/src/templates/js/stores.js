import {get, writable} from "svelte/store";

const loadingCount = writable(0);
export const loadingScreen = writable(true);
export const dropdown = writable(false);

export function addLoadingScreenTicket() {
    loadingCount.update(n => n + 1);
    loadingScreen.set(true);
}

export function removeLoadingScreenTicket() {
    loadingCount.update(n => n - 1);
    if (get(loadingCount) === 0) {
        loadingScreen.set(false);
    }
}

export const notifyModal = writable(false);
export const notifyTitle = writable("");
export const notifyMessage = writable("");