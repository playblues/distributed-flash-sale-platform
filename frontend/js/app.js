import { login, buyTicket, initStock } from './api.js';
import { NotificationSocket } from './websocket.js';

let authToken = '';
let wsClient = null;

const logsContainer = document.getElementById('logs');
const wsStatusBadge = document.getElementById('wsStatus');
const userIdInput = document.getElementById('userIdInput');
const priceInput = document.getElementById('priceInput');
const buyBtn = document.getElementById('buyBtn');
const initStockBtn = document.getElementById('initStockBtn');
const clearLogsBtn = document.getElementById('clearLogsBtn');

export function log(msg, type = 'info') {
  const time = new Date().toLocaleTimeString();
  let color = 'text-slate-300';
  if (type === 'success') color = 'text-emerald-400 font-medium';
  if (type === 'error') color = 'text-rose-400 font-medium';
  if (type === 'warn') color = 'text-amber-400';

  const entry = document.createElement('div');
  entry.className = color;
  entry.innerHTML = `<span class="text-slate-600 font-normal">[${time}]</span> ${msg}`;
  logsContainer.appendChild(entry);
  logsContainer.scrollTop = logsContainer.scrollHeight;
}

window.addEventListener('DOMContentLoaded', async () => {
  try {
    authToken = await login();
    log('JWT-токен успешно получен', 'success');
  } catch (e) {
    log('Ошибка логина через Gateway: ' + e.message, 'error');
  }

  initWebSocket();
  setupEventListeners();
});

function initWebSocket() {
  const userId = parseInt(userIdInput.value, 10) || 42;

  wsClient = new NotificationSocket(
    userId,
    (message) => {
      log(`Push от Kafka: ${message}`, 'success');
    },
    (isConnected, id) => {
      if (isConnected) {
        wsStatusBadge.className = 'text-xs px-3 py-1 bg-emerald-950/60 border border-emerald-800 text-emerald-300 rounded-full';
        wsStatusBadge.innerText = `WebSocket: Подключен (User ${id})`;
        log(`WebSocket соединен для пользователя ${id}`, 'success');
      } else {
        wsStatusBadge.className = 'text-xs px-3 py-1 bg-rose-950/60 border border-rose-800 text-rose-300 rounded-full';
        wsStatusBadge.innerText = 'WebSocket: Отключен';
        log('WebSocket соединение закрыто', 'warn');
      }
    }
  );

  wsClient.connect();
}

function setupEventListeners() {
  userIdInput.addEventListener('change', () => {
    initWebSocket();
  });

  clearLogsBtn.addEventListener('click', () => {
    logsContainer.innerHTML = '';
  });

  initStockBtn.addEventListener('click', async () => {
    try {
      await initStock(10);
      log('На склад завезено 10 билетов (Event 1)', 'success');
    } catch (e) {
      log('Ошибка инициализации склада: ' + e.message, 'error');
    }
  });

  buyBtn.addEventListener('click', async () => {
    const userId = parseInt(userIdInput.value, 10);
    const price = parseFloat(priceInput.value);

    buyBtn.disabled = true;
    buyBtn.innerText = 'Обработка...';

    try {
      const { status, data } = await buyTicket({
        userId,
        eventId: 1,
        price,
        token: authToken,
      });

      if (status === 201) {
        log(`Заказ №${data.id} создан (Статус: ${data.status}). Ожидание оплаты через Kafka...`, 'info');
      } else if (status === 409) {
        log(`Ошибка: ${data.error || 'Билеты закончились'}`, 'error');
      } else if (status === 429) {
        log(`Сработал Rate Limiter: Слишком много запросов!`, 'warn');
      } else {
        log(`Ошибка сервера (${status}): ${JSON.stringify(data)}`, 'error');
      }
    } catch (e) {
      log(`Ошибка сети: ${e.message}`, 'error');
    } finally {
      buyBtn.disabled = false;
      buyBtn.innerText = 'Забронировать и купить';
    }
  });
}