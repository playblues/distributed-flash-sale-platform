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

export async function initStock(amount = 10, eventId = 1) {
  const res = await fetch(`${INVENTORY_URL}/api/v1/inventory/init?event_id=${eventId}&amount=${amount}`, {
    method: 'POST',
  });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return await res.json();
}

export async function getStock(eventId = 1) {
  try {
    const res = await fetch(`${INVENTORY_URL}/api/v1/inventory/stock?event_id=${eventId}`);
    if (!res.ok) return null;
    return await res.json();
  } catch {
    return null;
  }
}
