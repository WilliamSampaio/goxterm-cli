if (typeof browser === "undefined") {
    var browser = chrome;
}

if (!browser.notifications) {
    browser.notifications = chrome.notifications;
}

if (!browser.storage) {
    browser.storage = chrome.storage;
}
