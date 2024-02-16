// Проверяем, есть ли в URL параметр token
const urlParams = new URLSearchParams(window.location.search);
const token = urlParams.get('token');
if (token) {
    // Сохраняем почту из параметра token в локальном хранилище
    localStorage.setItem('email', token);
}

// При отправке формы нового пароля
document.getElementById('new-password-form').addEventListener('submit', async function (event) {
    event.preventDefault(); // Предотвращаем отправку формы

    // Получаем значения из формы
    const email = localStorage.getItem('email');
    const newPassword = document.getElementById('new-password').value;

    const data = {
        email: email,
        Password: newPassword
    };

    try {
        // Отправляем данные на сервер
        const response = await fetch('http://localhost:8083/v1/newpass', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            console.log('Новый пароль успешно установлен');
            // Перенаправляем пользователя на страницу успешной смены пароля
            window.location.href = '/';
        } else {
            const errorMessage = await response.text();
            showError(errorMessage);
        }
    } catch (error) {
        console.error('Ошибка при отправке данных:', error.message);
    }
});
