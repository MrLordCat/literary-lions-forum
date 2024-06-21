document.addEventListener("DOMContentLoaded", () => {
    const carousels = document.querySelectorAll('.post-images');

    carousels.forEach((carousel) => {
        let currentSlide = 0;
        const items = carousel.querySelectorAll('.carousel-item');
        const prevButton = carousel.querySelector('.prev');
        const nextButton = carousel.querySelector('.next');

        if (items.length <= 1) {
            prevButton.style.display = 'none';
            nextButton.style.display = 'none';
        }

        function showSlide(index) {
            if (index >= items.length) {
                currentSlide = 0;
            } else if (index < 0) {
                currentSlide = items.length - 1;
            } else {
                currentSlide = index;
            }

            items.forEach((item, i) => {
                item.classList.toggle('active', i === currentSlide);
            });
        }

        function nextImage() {
            showSlide(currentSlide + 1);
        }

        function prevImage() {
            showSlide(currentSlide - 1);
        }

        if (items.length > 1) {
            nextButton.addEventListener('click', nextImage);
            prevButton.addEventListener('click', prevImage);
        }

        showSlide(currentSlide);
    });
});
