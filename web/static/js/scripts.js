document.addEventListener("DOMContentLoaded", function () {
    // Toggle form visibility
    document.querySelectorAll('.toggle-form').forEach(button => {
        button.addEventListener('click', function () {
            const formId = this.getAttribute('data-target');
            const form = document.getElementById(formId);
            if (form) {
                if (form.classList.contains('hidden')) {
                    form.classList.remove('hidden');
                } else {
                    form.classList.add('hidden');
                }
            }
        });
    });

    // Toggle comments visibility
    document.querySelectorAll('.toggle-comments').forEach(button => {
        button.addEventListener('click', function () {
            const commentsId = this.getAttribute('data-target');
            const comments = document.getElementById(commentsId);
            if (comments) {
                if (comments.classList.contains('hidden')) {
                    comments.classList.remove('hidden');
                } else {
                    comments.classList.add('hidden');
                }
            }
        });
    });
});
function toggleDropdown(dropdownId) {
    const dropdown = document.getElementById(dropdownId);
    if (dropdown.classList.contains('hidden')) {
        dropdown.classList.remove('hidden');
    } else {
        dropdown.classList.add('hidden');
    }
}
