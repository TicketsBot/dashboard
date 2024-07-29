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
    return `#${i.toString(16).padStart(6, '0')}`
}

export function nullIfBlank(s) {
    return s === '' ? null : s;
}

export function setBlankStringsToNull(obj) {
    // Set all blank strings in the object, including nested objects, to null
    for (const key in obj) {
        if (obj[key] === "" || obj[key] === "null") {
            obj[key] = null;
        } else if (typeof obj[key] === "object") {
            setBlankStringsToNull(obj[key]);
        }
    }
}

export function removeBlankEmbedFields(obj) {
    for (const key in obj) {
        if (obj[key] === null || obj[key] === undefined) {
            delete obj[key];
        }

        if (typeof obj[key] === "string" && obj[key] === "") {
            delete obj[key];
        }

        if (typeof obj[key] === "object") {
            removeBlankEmbedFields(obj[key]);
        }

        if (Array.isArray(obj[key]) && obj[key].length === 0) {
            delete obj[key];
        }
    }

    for (const key in obj) {
        if (typeof obj[key] === "object" && Object.keys(obj[key]).length === 0) {
            delete obj[key];
        }
    }
}

export function checkForParamAndRewrite(param) {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get(param) === "true") {
        const newUrl = new URL(window.location.href);
        newUrl.searchParams.delete(param);

        window.history.pushState(null, '', newUrl.toString());
        return true;
    }

    return false;
}
