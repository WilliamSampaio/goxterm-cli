(async () => {
    const IN_PAGE_ID = "___goxterm_in_page_script___";

    if (document.getElementById(IN_PAGE_ID)) return;

    const config = await getFromStorage("config");

    const connections = await getFromStorage("connections") || [];

    const inPage = document.createElement("div");
    inPage.id = IN_PAGE_ID;
    inPage.style.display = "none";

    const overlay = document.createElement("div");
    overlay.classList.add("goxterm_in_page_script_overlay");

    overlay.addEventListener("click", () => {
        inPage.style.display = "none";
    });

    inPage.appendChild(overlay);

    const dialog = document.createElement("div");
    dialog.classList.add("goxterm_in_page_script_basic", "goxterm_in_page_script_dialog");

    const content = document.createElement("div");
    content.style.display = "flex";
    content.style.flexDirection = "column";

    const label = document.createElement("b");
    label.textContent = "GoXterm | Search SSH:";

    const inputSearch = document.createElement("input");
    inputSearch.id = IN_PAGE_ID + "inputSearch";
    inputSearch.type = "text";
    inputSearch.classList.add("goxterm_in_page_script_basic", "goxterm_in_page_script_input");
    inputSearch.autocomplete = "off";

    content.appendChild(label);
    content.appendChild(inputSearch);
    dialog.appendChild(content);

    const listConnections = document.createElement("div");
    listConnections.id = IN_PAGE_ID + "listConnections";
    dialog.appendChild(listConnections);

    inPage.appendChild(dialog);
    document.body.appendChild(inPage);

    document.addEventListener("keydown", (e) => {
        const el = document.getElementById(IN_PAGE_ID);

        if (!el) return;

        if (e.code === "Enter" && connections.length > 0) {
            listConnections.children[1].click();
        }

        if (e.ctrlKey && e.code === "Backslash") {
            el.style.display = el.style.display === "none" ? "block" : "none";
            inputSearch.value = "";

            if (el.style.display === "block") inputSearch.focus();

            return;
        }

        if (el.style.display !== "none" && e.code === "Escape") {
            el.style.display = "none";
            inputSearch.value = "";
            return;
        }
    });

    function renderConnections(filter = "") {
        listConnections.innerHTML = "";

        if (filter === "") return;

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
            const url = new URL(`${config.backend_url}/web`);
            url.searchParams.append("ssh", true);
            url.searchParams.append("host", conn.host);
            url.searchParams.append("port", conn.port);
            url.searchParams.append("user", conn.user);

            if (conn.password) {
                url.searchParams.append("password", conn.password);
            }

            const item = document.createElement("div");
            item.classList.add("goxterm_in_page_script_list_item");

            const alias = document.createElement("strong");
            alias.textContent = conn.alias;

            item.appendChild(alias);
            item.appendChild(document.createElement("br"));

            const link = document.createElement("i");
            link.textContent = `${conn.user}@${conn.host}:${conn.port}`;

            item.addEventListener("click", () => {
                inPage.style.display = "none";
                window.open(url.toString(), "_blank");
            });

            item.appendChild(link);

            listConnections.appendChild(document.createElement("hr"));
            listConnections.appendChild(item);
        });
    }

    inputSearch.addEventListener("input", () => {
        renderConnections(inputSearch.value);
    });
})();