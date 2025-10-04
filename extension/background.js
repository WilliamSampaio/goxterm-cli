importScripts('polyfill.js');

const defaultConfig = {
    backend_url: 'http://localhost:8080',
    ssh_default_port: 22,
    ssh_default_user: 'root',
    ssh_default_password: null
};

chrome.runtime.onInstalled.addListener(async (details) => {
    if (details.reason === "install") {
        console.log("Install!");
        await chrome.storage.local.set({ config: defaultConfig });
    } else if (details.reason === "update") {
        console.log("Update!");
        const stored = await chrome.storage.local.get("config");
        const newConfig = { ...defaultConfig, ...stored.config };
        await chrome.storage.local.set({ config: newConfig });
    }
});
