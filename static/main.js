function q(sel) {
    return document.querySelector(sel);
}

let modalIsVisible = false;

function hideModal() {
    if (modalIsVisible) {
        q("#add-bookmark").style.display = "none";
        q("#toggle-add-btn").style.display = "block";
        modalIsVisible = false;
    }
}

function showModal() {
    if (! modalIsVisible) {
        q("#add-bookmark").style.display = "flex";
        q("#toggle-add-btn").style.display = "none";
        q("#txt-title").value = "";
        q("#txt-url").value = "";
        modalIsVisible = true;
    }
}

q("#toggle-add-btn").addEventListener('click', function () {
    if (modalIsVisible) {
        hideModal();
    } else {
        showModal();
    }
});

q("#confirm-btn").addEventListener('click', function () {
    hideModal();
});

q("#cancel-btn").addEventListener('click', function () {
    hideModal();
});
