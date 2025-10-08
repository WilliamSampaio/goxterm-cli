const inputBackend = document.getElementById("inputBackend");
const inputDefaultPort = document.getElementById("inputDefaultPort");
const inputDefaultUser = document.getElementById("inputDefaultUser");
const inputDefaultPassword = document.getElementById("inputDefaultPassword");
const btnSaveSettings = document.getElementById("btnSaveSettings");

const btnSaveBackup = document.getElementById("btnSaveBackup");
const formRestoreBackup = document.getElementById("formRestoreBackup");
const inputRestoreFile = document.getElementById("inputRestoreFile");

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

btnSaveBackup.addEventListener("click", saveBackup);

async function saveBackup() {

    const config = await getFromStorage("config");
    const connections = await getFromStorage("connections") || [];

    const backup = {
        config,
        connections
    }

    const jsonString = JSON.stringify(backup, null, null);
    const blob = new Blob([jsonString], { type: "application/json" });
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.style.display = "none";
    document.body.appendChild(a);

    a.href = url;
    a.download = `goxterm-backup-${new Date().toISOString().split("T")[0]}.json`;

    a.click();

    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    notify("GoXterm Extension", "Backup saved successfully.");
}

formRestoreBackup.addEventListener("submit", async (e) => {
    e.preventDefault();

    const file = inputRestoreFile.files[0];

    if (!file) {
        notify("Error | GoXterm Extension", "No file selected!");
        return;
    }

    try {
        const text = await file.text();
        const backup = JSON.parse(text);

        if (backup.config) await browser.storage.local.set({ config: backup.config });

        if (backup.connections) await browser.storage.local.set({ connections: backup.connections });

        notify("GoXterm Extension", "Backup restored successfully!");
    } catch (error) {
        console.error("Error restoring backup:", error);
        notify("Error | GoXterm Extension", "Failed to restore backup.");
    }
});
