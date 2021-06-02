import {writable} from "svelte/store";

export const loadingScreen = writable(true);

export const notifyModal = writable(false);
export const notifyTitle = writable("");
export const notifyMessage = writable("");