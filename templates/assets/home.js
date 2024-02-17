const createTeamButton = document.getElementById('create-team');
const teamNameInput = document.getElementById('team-name');
const teamContainer = document.getElementById('team-container');

createTeamButton.addEventListener('click', function() {
    const teamName = teamNameInput.value;
    if (teamName.trim() !== '') {
        const teamBlock = document.createElement('div');
        teamBlock.classList.add('team-block');
        teamBlock.innerHTML = `
            <h3>${teamName}</h3>
            <button class="join-team">Вступить</button>
        `;
        teamContainer.appendChild(teamBlock);
        teamNameInput.value = ''; // Очищаем поле ввода после создания команды
    } else {
        alert('Введите название команды');
    }
});

teamContainer.addEventListener('click', function(event) {
    if (event.target.classList.contains('join-team')) {
        const teamName = event.target.previousElementSibling.textContent;
        alert(`Вы вступили в команду "${teamName}"`);
    }
});
