function newCommentHTML(id, text, modified_at, ) {
    return '<div id="' + id +
        '" class="card border-0 mb-3"><div class="card-header bg-transparent"><label class="text-muted fs-6">Последнее обновление: ' +
        modified_at + '</label><button id="delete-comment_' +
        id + '" class="btn-delete-comment badge badge-danger bg-danger float-sm-end">Удалить</button></div><div class="card card-body bg-light">' +
        text + '</div><p></p></div>'
}

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
            newDeleteBtnID = "delete-comment_" + body.id
            var newCommentBlock = newCommentHTML(body.id, body.text, body.modified_at)
            document.getElementById(commentsBlockID).innerHTML += newCommentBlock
            newDeleteBtn = document.getElementById(newDeleteBtnID)
            newDeleteBtn.addEventListener('click', event => deleteCommentEvent(event, newDeleteBtn))
        }).catch((error) => {
            console.log(error) // TODO handle it better
        });
    })
})


document.querySelectorAll('.btn-delete-comment').forEach(item => {
    item.addEventListener('click', event => deleteCommentEvent(event, item))
})

function deleteCommentEvent(event, item) {
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
}