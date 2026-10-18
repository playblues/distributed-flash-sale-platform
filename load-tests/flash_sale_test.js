import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '3s', target: 50 },
        { duration: '10s', target: 100 },
        { duration: '3s', target: 0 },
    ],
};

const BASE_URL = 'http://localhost:3000';

export function setup() {
    const loginRes = http.post(`${BASE_URL}/api/v1/login`);
    const token = JSON.parse(loginRes.body).token;
    return { token: token };
}

export default function (data) {
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${data.token}`,
            'X-Idempotency-Key': `k6-${__VU}-${__ITER}-${Date.now()}`,
        },
    };

    const payload = JSON.stringify({
        userId: 42,
        eventId: 1,
        price: 5000.0,
    });

    const res = http.post(`${BASE_URL}/api/v1/buy`, payload, params);

    check(res, {
        'статус 201, 409 или 429': (r) => [201, 409, 429].includes(r.status),
    });

    sleep(0.05);
}