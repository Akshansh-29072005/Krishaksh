import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = { vus: 30, duration: '2m' };
const base = __ENV.BASE_URL || 'https://localhost/api/v1';
const token = __ENV.JWT_TOKEN || '';
const orderId = __ENV.ORDER_ID || '';

export default function () {
  const payload = JSON.stringify({ order_id: orderId });
  const params = { headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' } };
  const r = http.post(`${base}/payments/create-order`, payload, params);
  check(r, { 'payment order status': (x) => x.status === 200 || x.status === 201 || x.status === 400 });
  sleep(1);
}
