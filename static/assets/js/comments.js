document.querySelectorAll('.new-comment-form-submit').forEach(item => {
    item.addEventListener('click', event => {
        event.preventDefault();
        var partId = item.id.split("-").at(-1)
        var form = document.getElementById("new-comment-form-for-" + partId);
        let formData = new FormData(form);
        fetch(form.action, {
            method: form.method,
            body: formData,
        }).then((resp) => {
            return resp.json(); // or resp.text() or whatever the server sends
        }).then((body) => {
            commentsBlockID = "comments-for-part-" + partId
            var newCommentBlock = '<div id="' + body.id + '">' +
                '<button id="delete-comment-' +
                body.id +
                '" class="btn-delete-comment badge badge-danger bg-danger ms-1 float-end">X</button>' +
                '<div class="card card-body bg-light">' +
                body.text +
                '</div><p></p></div>'
            document.getElementById(commentsBlockID).innerHTML += newCommentBlock
        }).catch((error) => {
            console.log(error) // TODO handle it better
        });
    })
})

document.querySelectorAll('.btn-delete-comment').forEach(item => {
    item.addEventListener('click', event => {
        event.preventDefault();
        var commentId = item.id.split("_").at(-1)
        fetch('/comment?id=' + commentId, {
            method: 'DELETE',
        }).then((resp) => {
            return resp.text(); // or resp.text() or whatever the server sends
        }).then((body) => {
            var elem = document.getElementById(commentId);
            elem.remove();
        }).catch((error) => {
            console.log(error) // TODO handle it better
        });
    })
})