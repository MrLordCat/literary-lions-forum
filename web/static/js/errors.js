document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.getElementById("login-form");
   

    if (loginForm) {
        loginForm.addEventListener("submit", function (event) {
            event.preventDefault();

            const loginData = new FormData(loginForm);
            const loginError = document.getElementById("login-error");
            const passwordError = document.getElementById("password-error");

            loginError.textContent = "";
            passwordError.textContent = "";

            fetch("/login", {
                method: "POST",
                body: loginData
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    if (data.error.includes("User not found")) {
                        loginError.textContent = data.error;
                    } else if (data.error.includes("Invalid password")) {
                        passwordError.textContent = data.error;
                    }
                } else {
                    window.location.href = "/";
                }
            })
            .catch(error => {
                console.error("Error:", error);
            });
        });
    }
});
document.addEventListener("DOMContentLoaded", function () {
    const registerForm = document.getElementById("register-form");
    if (registerForm) {
        registerForm.addEventListener("submit", function (event) {
            event.preventDefault();

            const registerData = new FormData(registerForm);
            const usernameError = document.getElementById("username-error");
            const emailError = document.getElementById("email-error");

            usernameError.textContent = "";
            emailError.textContent = "";

            fetch("/register", {
                method: "POST",
                body: registerData
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    if (data.error.includes("Invalid email format")) {
                        emailError.textContent = data.error;
                    } else if (data.error.includes("Username or email already exists")) {
                        usernameError.textContent = data.error;
                    }
                } else {
                    window.location.href = "/";
                }
            })
            .catch(error => {
                console.error("Error:", error);
            });
        });
    }

});