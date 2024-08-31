export function isAnimated(icon) {
    if (icon === undefined || icon === "") {
        return false;
    } else {
        return icon.startsWith('a_')
    }
}

export function getIconUrl(id, icon) {
    if (!icon || icon === "") {
        return getDefaultIcon(id);
    }

    if (isAnimated(icon)) {
        return `https:\/\/cdn.discordapp.com/icons/${id}/${icon}.gif?size=256`
    } else {
        return `https:\/\/cdn.discordapp.com/icons/${id}/${icon}.webp?size=256`
    }
}

export function getDefaultIcon(id) {
    return `https://cdn.discordapp.com/embed/avatars/${Number((BigInt(id) >> BigInt(22)) % BigInt(6))}.png`
}