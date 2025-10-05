if (typeof browser === "undefined") {
    var browser = chrome;
}

if (!browser.notifications) {
    browser.notifications = chrome.notifications;
}

if (!browser.runtime) {
    browser.runtime = chrome.runtime;
}

if (!browser.storage) {
    browser.storage = chrome.storage;
}

if (!browser.tabs) {
    browser.tabs = chrome.tabs;
}