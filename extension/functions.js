async function saveOnStorage(key, value) {
    await browser.storage.local.set({ [key]: value });
}

async function getFromStorage(key) {
    const result = await browser.storage.local.get(key);
    return result[key];
}

async function removeFromStorage(key) {
    await browser.storage.local.remove(key);
}

async function clearStorage() {
    await browser.storage.local.clear();
}

function notify(title, message) {
    if (!browser.notifications) {
        console.warn("Notifications API not available.");
        return;
    }

    const options = {
        type: "basic",
        iconUrl: "icons/icon48.png",
        title: title,
        message: message
    };

    browser.notifications.create(options);
}

async function getCurrentTabHost() {
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

    if (!tab || !tab.url) return "0.0.0.0";

    try {
        const url = new URL(tab.url);
        return url.host;
    } catch {
        return "invalid";
    }
}

function isValidHostname(hostname) {
    // Regex to validate hostname (simplified version)
    const hostnameRegex = /^(?=.{1,253}$)(?!-)[A-Za-z0-9-]{1,63}(?<!-)(\.(?!-)[A-Za-z0-9-]{1,63}(?<!-))*$/;
    return hostnameRegex.test(hostname);
}

function isValidPort(port) {
    const portNumber = Number(port);
    return Number.isInteger(portNumber) && portNumber > 0 && portNumber <= 65535;
}

function isValidUser(user) {
    const userRegex = /^[a-zA-Z0-9._-]{1,32}$/;
    return userRegex.test(user);
}
