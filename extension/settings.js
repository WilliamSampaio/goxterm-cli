const inputBackend = document.getElementById("inputBackend");
const inputDefaultPort = document.getElementById("inputDefaultPort");
const inputDefaultUser = document.getElementById("inputDefaultUser");
const inputDefaultPassword = document.getElementById("inputDefaultPassword");
const btnSaveSettings = document.getElementById("btnSaveSettings");

document.addEventListener("DOMContentLoaded", async () => {
    const config = await getFromStorage("config");
    // console.info("Config from storage:", config);

    if (config && config.backend_url) {
        inputBackend.value = config.backend_url;
    }

    if (config && config.ssh_default_port) {
        inputDefaultPort.value = config.ssh_default_port;
    }

    if (config && config.ssh_default_user) {
        inputDefaultUser.value = config.ssh_default_user;
    }

    if (config && config.ssh_default_password) {
        inputDefaultPassword.value = config.ssh_default_password;
    }
});

btnSaveSettings.addEventListener("click", saveSettings);

function saveSettings() {
    const backend = inputBackend.value;
    const defaultPort = inputDefaultPort.value;
    const defaultUser = inputDefaultUser.value;
    const defaultPassword = inputDefaultPassword.value;

    if (!isValidPort(defaultPort)) {
        notify("Error! | GoXterm Extension", "Invalid default port. It should be a number between 1 and 65535.");
        return;
    }

    if (!isValidUser(defaultUser)) {
        notify("Error! | GoXterm Extension", "Invalid default user. It should be 1-32 characters long and can include letters, numbers, dots, underscores, and hyphens.");
        return;
    }

    saveOnStorage("config", {
        backend_url: backend,
        ssh_default_port: Number(defaultPort),
        ssh_default_user: defaultUser,
        ssh_default_password: defaultPassword || null
    }).then(() => {
        notify("GoXterm Extension", "Settings saved successfully.");
    }).catch((error) => {
        console.error("Error saving settings:", error);
        notify("Error! | GoXterm Extension", "Failed to save settings.");
    });
}
