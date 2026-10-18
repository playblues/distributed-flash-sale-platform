const GATEWAY_URL = 'http://localhost:3000';
const INVENTORY_URL = 'http://localhost:3001';

export async function login() {
  const res = await fetch(`${GATEWAY_URL}/api/v1/login`, { method: 'POST' });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  const data = await res.json();
  return data.token;
}

export async function buyTicket({ userId, eventId, price, token }) {
  const idempotencyKey = `ui-${Date.now()}-${Math.random().toString(36).substring(2, 7)}`;
  
  const res = await fetch(`${GATEWAY_URL}/api/v1/buy`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      'X-Idempotency-Key': idempotencyKey,
    },
    body: JSON.stringify({ userId, eventId, price }),
  });

  const data = await res.json();
  return { status: res.status, data };
}

export async function initStock(amount = 10) {
  const res = await fetch(`${INVENTORY_URL}/api/v1/inventory/init?amount=${amount}`, {
    method: 'POST',
  });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return await res.json();
}