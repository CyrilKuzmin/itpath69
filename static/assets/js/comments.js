function newCommentHTML(id, text, modified_at, ) {
    return '<div id="' + id +
        '" class="card border-0 mb-3"><div class="card-header bg-transparent"><label class="text-muted fs-6">Последнее обновление: ' +
        modified_at + '</label><button id="update-comment_' +
        id + '" class="btn-update-comment badge badge-primary bg-primary float-sm-end" onclick="makeCommentEditable(\'' +
        id + '\'">Изменить</button><button id="delete-comment_' +
        id + '" class="btn-delete-comment badge badge-danger bg-danger float-sm-end">Удалить</button></div><div id="text-' +
        id + '" class="card card-body bg-light">' +
        text + '</div><p></p></div>'
}

function updateCommentFormHTML(id, text) {
    '<form action="comment" method="put" role="form" class="update-comment-form" id="update-comment-form-for-' +
    id + '"><div class="mb-3"><input type="hidden" name="id" value="' +
        id + '"><label for="updatedComment' +
        id + '" class="form-label">Комментарий</label><textarea class="form-control" aria-label="text" name="text" id="updatedComment' +
        id + '" maxlength="1000">' +
        text + '</textarea></div><button type="submit" id="btn-update-update-' +
        id + '" class="btn btn-success btn-update-comment">Отправить</button>' +
        '<button type="cancel" onclick="cancelCommentEdit(' +
        id + '">Cancel</button></form>'
}

document.querySelectorAll('.btn-submit-comment').forEach(item => {
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
    var commentId = item.id.split("_").at(-1);
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

function updateCommentEvent(event, item) {
    event.preventDefault();
    var commentId = item.id.split("_").at(-1)
    var form = document.getElementById("update-comment-form-for-" + commentId);
    let formData = new FormData(form);
    fetch('/comment' + commentId, {
        method: 'PUT',
        body: formData,
    }).then((resp) => {
        return resp.text(); // or resp.text() or whatever the server sends
    }).then((body) => {
        var elem = document.getElementById(commentId);
        elem.remove();
    }).catch((error) => {
        console.log(error) // TODO handle it better
    });
}

function makeCommentEditable(commentId) {
    var textField = document.getElementById("text-" + commentId);
    var originalText = textField.innerHTML;
    textField.style.display = 'none'; // hide, 'flex' to show it again
    textField.innerHTML = updateCommentFormHTML(commentId, originalText);
    newUpdateBtnID = "btn-submit-update-" + body.id
    newUpdateBtn = document.getElementById(newUpdateBtnID)
    newUpdateBtn.addEventListener('click', event => updateCommentEvent(event, newUpdateBtn))
}

function cancelCommentEdit(commentId) {
    var textField = document.getElementById("text-" + commentId);
    var form = document.getElementById("update-comment-form-for-" + commentId)
    form.remove()
    textField.style.display = 'flex';
}

function saveUpdatesComment(commentId, text) {

}