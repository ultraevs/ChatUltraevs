// Обработчик для формы восстановления пароля
document.getElementById('forgot-form').addEventListener('submit', async function (event) {
    event.preventDefault(); // Prevent form submission

    // Retrieve form data for forgot password
    const email = document.getElementById('email').value;

    const data = {
        email: email
    };

    try {
        // Отправляем запрос на сервер
        const response = await fetch('http://localhost:8083/v1/forgot', {
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
            window.location.href = '/';
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