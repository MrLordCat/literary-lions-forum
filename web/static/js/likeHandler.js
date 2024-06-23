document.addEventListener("DOMContentLoaded", function() {
    window.handleLike = function(event, postID, likeType) {
        event.preventDefault();

        fetch('/like', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `post_id=${postID}&like_type=${likeType}`
        })
        .then(response => {
            if (!response.ok) {
                if (response.status === 400) {
                    const messageDiv = document.getElementById(`login-message-${postID}`);
                    if (messageDiv) {
                        messageDiv.style.display = 'block';
                        setTimeout(() => {
                            messageDiv.style.display = 'none';
                        }, 5000);
                    }
                    throw new Error('You must be logged in to like/dislike.');
                }
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                document.getElementById(`post-likes-${postID}`).innerText = `${data.likes}`;
                document.getElementById(`post-dislikes-${postID}`).innerText = `${data.dislikes}`;
            } else {
                alert('Failed to update like/dislike');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    };

    window.likeComment = function(commentID, likeType) {
        fetch('/like-comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({
                'comment_id': commentID,
                'like_type': likeType
            })
        })
        .then(response => {
            if (!response.ok) {
                if (response.status === 400) {
                    const messageDiv = document.getElementById(`login-message-comment-${commentID}`);
                    if (messageDiv) {
                        messageDiv.style.display = 'block';
                        setTimeout(() => {
                            messageDiv.style.display = 'none';
                        }, 5000);
                    }
                    throw new Error('You must be logged in to like/dislike.');
                }
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                document.getElementById(`comment-likes-${commentID}`).textContent = data.likes;
            } else {
                alert('Failed to like comment.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    };
});
