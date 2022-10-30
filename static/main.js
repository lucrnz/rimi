function q(sel) {
    return document.querySelector(sel);
}

function renderBookmark(bookmark) {
    const { title, url } = bookmark;
    if (typeof title === 'string' &&
        typeof url === 'string' &&
        title.length > 0 &&
        url.length > 0) {
        const a = document.createElement("a");
        a.innerText = title;
        a.href = url;
        q("#bookmarks").appendChild(a);
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

q("#add-btn").addEventListener('click', async function () {
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
