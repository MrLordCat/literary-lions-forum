function addComment(event, postID) {
    event.preventDefault(); // Останавливаем стандартное поведение формы

    const form = document.getElementById(`add-comment-form-${postID}`);
    const formData = new FormData(form);
    const commentsList = document.getElementById(`comments-list-${postID}`);
    const noComments = document.getElementById(`no-comments-${postID}`);

    fetch('/add-comment', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        if (noComments) {
            noComments.style.display = 'none'; // Скрываем сообщение "No comments yet"
        }

        // Создаем новый элемент комментария
        const newComment = document.createElement('li');
        newComment.id = `comment-${data.ID}`;
        newComment.innerHTML = `
            <p>${data.Content}</p>
            <small>Posted at ${new Date(data.CreatedAt).toLocaleString()}</small>
            <div>
                Likes: 0
                <form method="POST" action="/like-comment" style="display:inline;">
                    <input type="hidden" name="comment_id" value="${data.ID}">
                    <input type="hidden" name="like_type" value="1">
                    <button id="postbutton" type="submit">👍 Like</button>
                </form>
            </div>
        `;

        commentsList.appendChild(newComment); // Добавляем новый комментарий в список
        form.reset(); // Сбрасываем форму
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
