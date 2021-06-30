import * as Stores from './stores'
import {navigateTo} from "svelte-router-spa";

export async function withLoadingScreen(func) {
    Stores.addLoadingScreenTicket();
    await func();
    Stores.removeLoadingScreenTicket();
}

export function errorPage(message) {
    navigateTo(`/error?message=${encodeURIComponent(message)}`)
}

export function notify(title, message) {
    Stores.notifyTitle.set(title);
    Stores.notifyMessage.set(message);
    Stores.notifyModal.set(true);
}

export function notifyError(message) {
    notify('Error', message);
}

export function notifySuccess(message) {
    notify('Success', message);
}

export function notifyRatelimit() {
    notifyError("You're doing that too fast: please wait a few seconds and try again");
}

export function closeNotificationModal() {
    Stores.notifyModal.set(false);
}

export function colourToInt(colour) {
    return parseInt(`0x${colour.slice(1)}`);
}

export function intToColour(i) {
    return `#${i.toString(16)}`
}
