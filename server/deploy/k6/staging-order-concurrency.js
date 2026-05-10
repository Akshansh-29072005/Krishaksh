import http from 'k6/http';
import { check } from 'k6';

export const options = { vus: 40, duration: '4m' };
const base = __ENV.BASE_URL || 'https://localhost';
const token = __ENV.JWT_TOKEN || '';

export default function () {
  const payload = JSON.stringify({ shipping_metadata: { city: 'Delhi' }, notes: 'staging load test' });
  const r = http.post(`${base}/api/v1/orders`, payload, {
    headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
  });
  check(r, { 'order endpoint alive': (x) => [201, 400, 401].includes(x.status) });
}
