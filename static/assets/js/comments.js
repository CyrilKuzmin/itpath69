function newCommentHTML(id, text, modified_at, ) {
    return '<div id="' + id +
        '" class="card border-0 mb-3"><div class="card-header bg-transparent"><label class="text-muted fs-6">Последнее обновление: ' +
        '</label><label id="last-update-time_' +
        id + '" class="text-muted fs-6 ms-1">' +
        modified_at + '</label><button id="delete-comment_' +
        id + '" class="btn-delete-comment badge badge-warning bg-warning ms-2 float-sm-end border-0">Удалить</button><button id="update-comment_' +
        id + '" class="btn-update-comment badge badge-info bg-info badge-pill ms-2 float-sm-end border-0" onclick="makeCommentEditable(\'' +
        id + '\')">Изменить</button></div><div id="text-' +
        id + '" class="multiline-card card card-body bg-light">' +
        text + '</div><p></p></div>'
}

function updateCommentFormHTML(id, text) {
    return '<form action="comment" method="put" role="form" class="update-comment-form" id="update-comment-form-for-' +
        id + '"><div class="mb-3"><input type="hidden" name="id" value="' +
        id + '"><textarea class="form-control" aria-label="text" name="text" id="updatedComment' +
        id + '" maxlength="1000">' +
        text + '</textarea></div><button type="submit" id="btn-submit-update_' +
        id + '" class="btn btn-success btn-update-comment">Отправить</button>' +
        '<button type="cancel" class="btn btn-light ms-2" onclick="cancelCommentEdit(\'' +
        id + '\')">Отмена</button></form>'
}

const tx = document.getElementsByTagName("textarea");
for (let i = 0; i < tx.length; i++) {
    tx[i].setAttribute("style", "height:" + (tx[i].scrollHeight) + "px;overflow-y:hidden;");
    tx[i].addEventListener("input", OnInput, false);
}

function OnInput() {
    this.style.height = "auto";
    this.style.height = (this.scrollHeight) + "px";
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
            if (!resp.ok) {
                return resp.text();
            } else {
                return resp.json(); // or resp.text() or whatever the server sends
            }
        }).then((body) => {
            commentsBlockID = "comments-for-part-" + partId
            newDeleteBtnID = "delete-comment_" + body.id
            var newCommentBlock = newCommentHTML(body.id, body.text, body.modified_at)
            document.getElementById(commentsBlockID).innerHTML += newCommentBlock;
            newDeleteBtn = document.getElementById(newDeleteBtnID)
            newDeleteBtn.addEventListener('click', event => deleteCommentEvent(event, body.id))
        }).catch((error) => {
            displayError(form, error);
        });
    })
})


document.querySelectorAll('.btn-delete-comment').forEach(item => {
    var commentId = item.id.split("_").at(-1);
    item.addEventListener('click', event => deleteCommentEvent(event, commentId))
})

function deleteCommentEvent(event, commentId) {
    var result = confirm("Удалить комментарий?");
    if (result) {
        event.preventDefault();
        fetch('/comment?id=' + commentId, {
            method: 'DELETE',
        }).then((resp) => {
            return resp.text(); // or resp.text() or whatever the server sends
        }).then((body) => {
            var elem = document.getElementById(commentId);
            elem.remove();
        }).catch((error) => {
            console.log(error) // TODO handle it better
        })
    };
}

function updateCommentEvent(event, commentId) {
    event.preventDefault();
    var form = document.getElementById("update-comment-form-for-" + commentId);
    let formData = new FormData(form);
    fetch('/comment', {
        method: 'PUT',
        body: formData,
    }).then((resp) => {
        if (!resp.ok) {
            return resp.text();
        } else {
            return resp.json(); // or resp.text() or whatever the server sends
        }
    }).then((body) => {
        form.remove();
        saveUpdatesComment(body.id, body.text, body.modified_at);
    }).catch((error) => {
        displayError(form, error);
    });
}

function makeCommentEditable(commentId) {
    var commentDiv = document.getElementById(commentId);
    var textField = document.getElementById("text-" + commentId);
    var originalText = textField.innerHTML;
    textField.style.display = 'none'; // hide, 'flex' to show it again
    document.getElementById('update-comment_' + commentId).style.display = 'none';
    document.getElementById('delete-comment_' + commentId).style.display = 'none';
    commentDiv.innerHTML += updateCommentFormHTML(commentId, originalText);
    newUpdateBtnID = "btn-submit-update_" + commentId
    newUpdateBtn = document.getElementById(newUpdateBtnID)
    newUpdateBtn.addEventListener('click', event => updateCommentEvent(event, commentId))
}

function cancelCommentEdit(commentId) {
    var textField = document.getElementById("text-" + commentId);
    var form = document.getElementById("update-comment-form-for-" + commentId)
    form.remove()
    textField.style.display = 'flex';
    document.getElementById('update-comment_' + commentId).style.display = 'inline-block';
    document.getElementById('delete-comment_' + commentId).style.display = 'inline-block';
}

function saveUpdatesComment(commentId, text, modified_at) {
    textField = document.getElementById("text-" + commentId);
    textField.innerHTML = text;
    textField.style.display = 'flex';
    deleteBtn = document.getElementById('delete-comment_' + commentId);
    document.getElementById('update-comment_' + commentId).style.display = 'inline-block';
    document.getElementById('last-update-time_' + commentId).innerHTML = modified_at;
    deleteBtn.style.display = 'inline-block';
    deleteBtn.addEventListener('click', event => deleteCommentEvent(event, commentId))
}

function displayError(thisForm, error) {
    thisForm.querySelector('.loading').classList.remove('d-block');
    thisForm.querySelector('.error-message').innerHTML = error;
    thisForm.querySelector('.error-message').classList.add('d-block');
}