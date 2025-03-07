let intervalId;
let taskIntervalId;

function showMessage(message) {
    const messageDiv = document.getElementById('message');
    messageDiv.innerText = message;
}

function updateInterval() {
    const interval = document.getElementById('interval').value * 1000;
    if (intervalId) {
        clearInterval(intervalId);
    }
    intervalId = setInterval(fetchExpressions, interval);
    if (taskIntervalId) {
        clearInterval(taskIntervalId);
    }
    taskIntervalId = setInterval(fetchTasks, interval);
}

function sendCalculationRequest() {
    const host = document.getElementById('host').value;
    const expression = document.getElementById('expression').value;
    fetch(`${host}/api/v1/calculate`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ expression })
    })
        .then(response => response.json())
        .then(data => showMessage(`Запрос отправлен. Id: ${data.id}`))
        .catch(error => console.error('Ошибка:', error));
}

function fetchExpressions() {
    const tableBody = document.getElementById('expressionsTable').getElementsByTagName('tbody')[0];
    tableBody.innerHTML = '';
    const host = document.getElementById('host').value;
    fetch(`${host}/api/v1/expressions`)
        .then(response => {
            if (response.status === 404) {
                const row = tableBody.insertRow();
                const cell = row.insertCell(0);
                cell.colSpan = 3;
                cell.innerText = 'Нет выражений';
                return;
            }
            return response.json();
        })
        .then(data => {
            if (data) {
                data.expressions.forEach(expression => {
                    const row = tableBody.insertRow();
                    row.insertCell(0).innerText = expression.id;
                    row.insertCell(1).innerText = expression.status;
                    row.insertCell(2).innerText = expression.result;
                });
            }
        })
        .catch(error => console.error('Ошибка:', error));
}

function fetchTasks() {
    const tableBody = document.getElementById('tasksTable').getElementsByTagName('tbody')[0];
    tableBody.innerHTML = '';
    const host = document.getElementById('host').value;
    fetch(`${host}/internal/tasks`)
        .then(response => response.json())
        .then(data => {
            data.tasks.forEach(task => {
                const row = tableBody.insertRow();
                row.insertCell(0).innerText = task.id;
                row.insertCell(1).innerText = task.arg1;
                row.insertCell(2).innerText = task.arg2;
                row.insertCell(3).innerText = task.operation;
                row.insertCell(4).innerText = task.operation_time;
                row.insertCell(5).innerText = task.is_busy ? 'Да' : 'Нет';
            });
        })
        .catch(error => console.error('Ошибка:', error));
}

function fetchExpressionById() {
    const host = document.getElementById('host').value;
    const id = document.getElementById('expressionId').value;
    fetch(`${host}/api/v1/expressions/${id}`)
        .then(response => response.json())
        .then(data => {
            const details = document.getElementById('expressionDetails');
            details.innerHTML = `
                    <p>Id: ${data.expression.id}</p>
                    <p>Status: ${data.expression.status}</p>
                    <p>Result: ${data.expression.result}</p>
                `;
        })
        .catch(error => console.error('Ошибка:', error));
}

function fetchTask() {
    const host = document.getElementById('host').value;
    fetch(`${host}/internal/task`)
        .then(response => response.json())
        .then(data => {
            const details = document.getElementById('taskDetails');
            details.innerHTML = `
                    <p>Id: ${data.task.id}</p>
                    <p>Arg1: ${data.task.arg1}</p>
                    <p>Arg2: ${data.task.arg2}</p>
                    <p>Operation: ${data.task.operation}</p>
                    <p>OperationTime: ${data.task.operation_time}</p>
                `;
        })
        .catch(error => console.error('Ошибка:', error));
}

function sendTaskResult() {
    const host = document.getElementById('host').value;
    const id = document.getElementById('taskId').value;
    const result = document.getElementById('taskResult').value;
    fetch(`${host}/internal/task`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id, result: parseFloat(result) })
    })
        .then(response => {
            if (response.ok) {
                showMessage('Результат отправлен');
            } else {
                showMessage('Ошибка отправки результата');
            }
        })
        .catch(error => console.error('Ошибка:', error));
}

document.getElementById('interval').addEventListener('change', updateInterval);
updateInterval(); // Initialize interval on page load
fetchExpressions(); // Fetch expressions on page load
fetchTasks(); // Fetch tasks on page load