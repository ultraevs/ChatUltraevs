// Получаем кнопки "Логин" и "Регистрация"
const loginButton = document.getElementById('login-button');
const registerButton = document.getElementById('register-button');

// Получаем формы логина и регистрации
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');

// Показываем форму логина при нажатии на кнопку "Логин"
loginButton.addEventListener('click', function() {
    loginForm.style.display = 'block';
    registerForm.style.display = 'none';
});

// Показываем форму регистрации при нажатии на кнопку "Регистрация"
registerButton.addEventListener('click', function() {
    loginForm.style.display = 'none';
    registerForm.style.display = 'block';
});

document.getElementById('forgot-password').addEventListener('click', function() {
    window.location.href = '/forgot';
});
// Обработчик для формы логина
document.getElementById('login-form').addEventListener('submit', async function (event) {
    event.preventDefault(); // Prevent form submission

    // Retrieve form data for login
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const data = {
        email: email,
        password: password
    };

    try {
        // Отправляем запрос на сервер
        const response = await fetch('http://localhost:8083/v1/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data) // Преобразуем данные в формат JSON
        });

        // Проверяем успешность запроса
        if (response.ok) {
            console.log('Данные успешно отправлены');
            // Перенаправляем пользователя на /home
            window.location.href = '/home';
        } else {
            // Отображаем ошибку
            const errorMessage = await response.text(); // Получаем текст ошибки
            showError(errorMessage);
        }
    } catch (error) {
        console.error('Ошибка при отправке данных:', error.message);
        // Здесь можно обработать ошибку
    }
});

// Обработчик для формы авторизации
document.getElementById('register-form').addEventListener('submit', async function (event) {
    event.preventDefault(); // Prevent form submission

    // Retrieve form data for registration
    const username = document.getElementById('username').value;
    const email = document.getElementById('emailreg').value;
    const password = document.getElementById('passwordreg').value;

    const data = {
        username: username,
        email: email,
        password: password
    };

    try {
        // Отправляем запрос на сервер
        const response = await fetch('http://localhost:8083/v1/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data) // Преобразуем данные в формат JSON
        });

        // Проверяем успешность запроса
        if (response.ok) {
            console.log('Данные успешно отправлены');
            // Перенаправляем пользователя на /home
            // window.location.href = '/home';
        } else {
            // Отображаем ошибку
            const errorMessage = await response.text(); // Получаем текст ошибки
            showError(errorMessage);
        }
    } catch (error) {
        console.error('Ошибка при отправке данных:', error.message);
        // Здесь можно обработать ошибку
    }
});

function showError(message) {
    const errorDiv = document.getElementById('error-message');
    errorDiv.textContent = message;
    errorDiv.style.color = 'red';
}