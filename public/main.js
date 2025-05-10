const api = 'https://jubilant-tribble-7rwgxwp5q6p2pjjp-8081.app.github.dev';

document.addEventListener('DOMContentLoaded', () => {
    fetch(`${api}/civilisations`)
        .then(res => res.json())
        .then(civs => {
            const select = document.getElementById('civ');
            civs.forEach(civ => {
                const option = document.createElement('option');
                option.value = civ.code;
                option.textContent = civ.name;
                select.appendChild(option);
            });
        })
        .catch(err => {
            console.error('Failed to load civilisations:', err);
        });
});