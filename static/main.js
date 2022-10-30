// @license magnet:?xt=urn:btih:5ac446d35272cc2e4e85e4325b146d0b7ca8f50c&dn=unlicense.txt Unlicense
// rimi bookmark manager: https://git.lucdev.net/luc/rimi

function q(sel) {
    return document.querySelector(sel);
}

let isDeleteMode = false;

async function deleteBookmarks(urlArray) {
    const promises = [];
    for (const url of urlArray ) {
        promises.push(fetch("/api/bookmarks", {
            method: "DELETE",
            cache: 'no-cache',
            headers: {
                'Content-Type': "text/json"
            },
            body: JSON.stringify({
                url
            })
        }));
    }
    return Promise.all(promises);
}

function renderBookmark(bookmark) {
    const { title, url } = bookmark;
    if (typeof title === 'string' &&
        typeof url === 'string' &&
        title.length > 0 &&
        url.length > 0) {
        const div = document.createElement("div");
        const checkmark = document.createElement("input");
        checkmark.type = "checkbox";
        checkmark.style.display = "none";
        
        const link = document.createElement("a");
        link.classList.add("bookmark");
        link.innerText = title;
        link.href = url;

        div.appendChild(checkmark);
        div.appendChild(link);
        q("#bookmarks").appendChild(div);
    }
}

async function loadBookmarks() {
    const response = await fetch("/api/bookmarks");
    const data = await response.json();
    
    if (Array.from(q("#bookmarks").children).length > 0) {
        q("#bookmarks").innerHTML = ""; 
    }

    data.forEach(bookmark => renderBookmark(bookmark));
}

q("#del-btn").addEventListener('click', async () => {
    if ( q("#del-btn").dataset.mode === "check" ) {
        const elements = Array.from(q("#bookmarks").children);
        if (elements.length > 0) {
            for (const el of elements) {
                const check = el.querySelector('input[type=checkbox]');
                if (check) {
                    check.style.display = "block";
                }
            }
        }
        q("#del-btn").value = "Confirm delete";
        q("#del-btn").dataset.mode = "confirm";
    } else if ( q("#del-btn").dataset.mode === "confirm") {
        const elements = Array.from(q("#bookmarks").children);

        const urlArray = [];

        if (elements.length > 0) {
            for (const el of elements) {
                const check = el.querySelector('input[type=checkbox]');
                if (check) {
                    check.style.display = "none";
                    if (check.checked) {
                        const link = el.querySelector('a.bookmark');
                        if (link && link.href.length > 0) {
                            urlArray.push(link.href);
                        }
                    }
                }
            }
        }

        if (urlArray.length > 0 ) {
            await deleteBookmarks(urlArray);
            await loadBookmarks();
        }
        
        q("#del-btn").value = "Delete";
        q("#del-btn").dataset.mode = "check"; 
    }
});

q("#add-btn").addEventListener('click', async function () {
    if (isDeleteMode) {
        return;
    }

    const url = q("#txt-url").value;
    const title = q("#txt-title").value;

    if (url.length > 0 && title.length > 0) {
        await fetch("/api/bookmarks", {
            method: "POST",
            cache: "no-cache",
            headers: {
                'Content-Type': "text/json"
            },
            body: JSON.stringify({
                title,
                url
            })
        })
        await loadBookmarks();
        q("#txt-url").value = "";
        q("#txt-title").value = "";
    }
});

(async () => {
    await loadBookmarks();
})();

// @license-end