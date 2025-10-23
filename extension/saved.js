const inputSearch = document.getElementById("inputSearch");
const listConnections = document.getElementById("listConnections");

document.addEventListener("DOMContentLoaded", async () => {
    const config = await getFromStorage("config");
    // console.info("Config from storage:", config);

    const connections = await getFromStorage("connections") || [];
    console.info("Connections from storage:", connections);

    function renderConnections(filter = "") {
        listConnections.innerHTML = "";

        const filteredConnections = connections.filter(c =>
            c.alias.toLowerCase().includes(filter.toLowerCase()) ||
            c.host.toLowerCase().includes(filter.toLowerCase()) ||
            c.user.toLowerCase().includes(filter.toLowerCase())
        );

        if (filteredConnections.length === 0) {
            const noConnItem = document.createElement("div");
            noConnItem.textContent = "No connections found.";
            listConnections.appendChild(noConnItem);
            return;
        }

        filteredConnections.forEach(conn => {
            const url = new URL(`${config.backend_url}`);
            url.searchParams.append("ssh", true);
            url.searchParams.append("host", conn.host);
            url.searchParams.append("port", conn.port);
            url.searchParams.append("user", conn.user);

            if (conn.password) {
                url.searchParams.append("password", conn.password);
            }

            const item = document.createElement("div");
            item.style.display = "flex";

            const content = document.createElement("div");
            content.style.marginRight = "10px";
            content.style.marginBottom = "0";

            const alias = document.createElement("strong");
            alias.textContent = conn.alias;
            content.appendChild(alias);

            content.appendChild(document.createElement("br"));

            const link = document.createElement("a");
            link.href = url.toString();
            link.textContent = `${conn.user}@${conn.host}:${conn.port}`;
            link.target = "_blank";

            content.appendChild(link);

            const btnDelete = document.createElement("button");
            btnDelete.textContent = "Delete";
            btnDelete.classList.add("small");
            btnDelete.style.marginLeft = "auto";

            btnDelete.addEventListener("click", async () => {
                const index = connections.findIndex(c => c.alias === conn.alias);
                if (index !== -1) {
                    connections.splice(index, 1);
                    await saveOnStorage("connections", connections);
                    renderConnections(inputSearch.value);
                }
            });

            item.appendChild(content);
            item.appendChild(btnDelete);

            listConnections.appendChild(document.createElement("hr"));
            listConnections.appendChild(item);
        });
    }

    inputSearch.addEventListener("input", () => {
        renderConnections(inputSearch.value);
    });

    renderConnections();
});