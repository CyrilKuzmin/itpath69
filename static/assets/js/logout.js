btn = document.getElementById("logout-btn")
btn.addEventListener('click', event => logoutEvent(event, commentId))


function logoutEvent(event, commentId) {
    event.preventDefault();
    fetch('/logout', {
        method: 'POST',
        headers: [("Content-Type", "application/x-www-form-urlencoded")]
    }).then((resp) => {
        if (resp.redirected) {
            window.location.href = "/";
        }
    }).catch((error) => {
        console.log(error) // TODO handle it better
    })
}