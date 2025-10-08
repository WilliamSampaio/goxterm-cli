const backendStatus = document.getElementById('backendStatus');
const inputAlias = document.getElementById('inputAlias');
const inputHost = document.getElementById('inputHost');
const inputPort = document.getElementById('inputPort');
const inputUser = document.getElementById('inputUser');
const inputPassword = document.getElementById('inputPassword');
// const inputPrivateKey = document.getElementById('inputPrivateKey');
// const inputPassphrase = document.getElementById('inputPassphrase');

const btnConnectSSH = document.getElementById('btnConnectSSH');
const btnSaveConnection = document.getElementById('btnSaveConnection');
const btnDeleteConnection = document.getElementById('btnDeleteConnection');

const terminalsMessage = document.getElementById('terminalsMessage');
const formSelectShell = document.getElementById('formSelectShell');

document.addEventListener("DOMContentLoaded", async () => {
    const config = await getFromStorage("config");
    // console.info("Config from storage:", config);

    const currentHost = String(await getCurrentTabHost()).split(':')[0];
    console.info("Current tab host:", currentHost);

    const existingConnection = await getConnectionByHost(currentHost);
    console.info("Existing connection for current host:", existingConnection);

    if (existingConnection) {
        inputAlias.value = existingConnection.alias;
        inputAlias.readOnly = true;
        inputHost.value = existingConnection.host;
        inputPort.value = existingConnection.port;
        inputUser.value = existingConnection.user;
        inputPassword.value = existingConnection.password;
    } else {
        inputHost.value = currentHost;

        if (config && config.ssh_default_port) {
            inputPort.value = config.ssh_default_port;
        }

        if (config && config.ssh_default_user) {
            inputUser.value = config.ssh_default_user;
        }

        if (config && config.ssh_default_password) {
            inputPassword.value = config.ssh_default_password;
        }
    }

    const backendGetInfo = async () => {
        if (formSelectShell.children.length > 1) {
            formSelectShell.children.slice(1).forEach(c => c.remove());
        };
        return fetch(`${config.backend_url}/api/info`)
            .then(async response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}.`);
                }
                const { shells } = await response.json();

                shells.forEach(s => {
                    if (!s.default) return;

                    const button = document.createElement('button');
                    button.style.marginBottom = '2px';
                    button.name = 'path';
                    button.value = s.path;
                    button.textContent = s.bin + ' (default)';

                    formSelectShell.append(button);
                });

                shells.forEach(s => {
                    if (s.default) return;

                    const button = document.createElement('button');
                    button.style.marginBottom = '2px';
                    button.name = 'path';
                    button.value = s.path;
                    button.textContent = s.bin;

                    formSelectShell.append(button);
                });
            })
            .catch(error => {
                console.error('Error fetching API data:', error);
            });
    }

    const setBackendStatus = (online = true) => {
        const inputs = document.querySelectorAll('input');
        const buttons = document.querySelectorAll('button');
        if (online) {
            backendStatus.textContent = "Online";
            backendStatus.style.color = "green";
            inputs.forEach(input => input.disabled = false);
            buttons.forEach(input => input.disabled = false);
            terminalsMessage.textContent = '';
        } else {
            backendStatus.textContent = "Offline";
            backendStatus.style.color = "red";
            inputs.forEach(input => input.disabled = true);
            buttons.forEach(input => input.disabled = true);
            terminalsMessage.textContent = 'No terminals available. Backend is offline.';
            terminalsMessage.style.color = "red";
        }
    }

    const backendIsOnline = async () => {
        return fetch(`${config.backend_url}/api/ping`)
            .then(async response => {
                if (!response.ok) {
                    setBackendStatus(false);
                    throw new Error(`HTTP error! status: ${response.status}.`);
                }

                const { alive } = await response.json();

                if (alive === true) {
                    setBackendStatus(true);
                    if (formSelectShell.children.length === 1) {
                        backendGetInfo();
                    }
                } else {
                    setBackendStatus(false);
                }
            })
            .catch(error => {
                setBackendStatus(false);
                console.error('Error fetching API data:', error);
            });
    }

    backendIsOnline();
    setInterval(backendIsOnline, 2000);
    setInterval(backendGetInfo, 30 * 1000);

    formSelectShell.action = (config && config.backend_url) ? `${config.backend_url}/web` : '';

    btnConnectSSH.addEventListener("click", () => {
        const host = inputHost.value.trim();
        const port = inputPort.value.trim();
        const user = inputUser.value.trim();
        const password = inputPassword.value;

        if (!isValidHostname(host)) {
            notify("Error! | GoXterm Extension", "Invalid hostname.");
            return;
        }

        if (!isValidPort(port)) {
            notify("Error! | GoXterm Extension", "Invalid port. It should be a number between 1 and 65535.");
            return;
        }

        if (!isValidUser(user)) {
            notify("Error! | GoXterm Extension", "Invalid user. It should be 1-32 characters long and can include letters, numbers, dots, underscores, and hyphens.");
            return;
        }

        if (!config || !config.backend_url) {
            notify("Error! | GoXterm Extension", "Backend URL is not configured. Please set it in the extension settings.");
            return;
        }

        const url = new URL(`${config.backend_url}/web`);
        url.searchParams.append("ssh", true);
        url.searchParams.append("host", host);
        url.searchParams.append("port", port);
        url.searchParams.append("user", user);

        if (password) {
            url.searchParams.append("password", password);
        }

        window.open(url.toString(), '_blank');
    });

    btnSaveConnection.addEventListener("click", async () => {
        const alias = inputAlias.value.trim();
        const host = inputHost.value.trim();
        const port = inputPort.value.trim();
        const user = inputUser.value.trim();
        const password = inputPassword.value;

        if (!alias) {
            notify("Error! | GoXterm Extension", "Alias is required.");
            return;
        }

        if (!isValidHostname(host)) {
            notify("Error! | GoXterm Extension", "Invalid hostname.");
            return;
        }

        if (!isValidPort(port)) {
            notify("Error! | GoXterm Extension", "Invalid port. It should be a number between 1 and 65535.");
            return;
        }

        if (!isValidUser(user)) {
            notify("Error! | GoXterm Extension", "Invalid user. It should be 1-32 characters long and can include letters, numbers, dots, underscores, and hyphens.");
            return;
        }

        const connection = {
            alias: alias,
            host: host,
            port: Number(port),
            user: user,
            password: password || null,
        };

        const connections = await getConnectionsList();

        const existingIndex = connections.findIndex(conn => conn.alias === alias);
        if (existingIndex !== -1) {
            connections[existingIndex] = connection;
        } else {
            connections.push(connection);
        }

        await saveOnStorage("connections", connections);

        inputAlias.readOnly = true;

        notify("GoXterm Extension", "Connection saved successfully.");
    });

    btnDeleteConnection.addEventListener("click", async () => {
        const alias = inputAlias.value.trim();

        if (!alias) {
            notify("Error! | GoXterm Extension", "Alias is required to delete a connection.");
            return;
        }

        const connections = await getConnectionsList();

        const updatedConnections = connections.filter(conn => conn.alias !== alias);
        if (updatedConnections.length === connections.length) {
            notify("Error! | GoXterm Extension", "No connection found with the given alias.");
            return;
        }

        await saveOnStorage("connections", updatedConnections);

        inputAlias.value = '';
        inputAlias.readOnly = false;
        inputHost.value = currentHost;

        if (config && config.ssh_default_port) {
            inputPort.value = config.ssh_default_port;
        } else {
            inputPort.value = '';
        }

        if (config && config.ssh_default_user) {
            inputUser.value = config.ssh_default_user;
        } else {
            inputUser.value = '';
        }

        if (config && config.ssh_default_password) {
            inputPassword.value = config.ssh_default_password;
        } else {
            inputPassword.value = '';
        }

        notify("GoXterm Extension", "Connection deleted successfully.");
    });
});

async function getConnectionsList() {
    return getFromStorage("connections").then(connections => connections || []);
}

async function getConnectionByHost(host) {
    return getConnectionsList().then(connections => connections.find(conn => conn.host === host));
}