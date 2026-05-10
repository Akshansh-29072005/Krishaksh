import http from 'k6/http';
import { check } from 'k6';

export const options = { vus: 100, duration: '3m' };
const base = __ENV.BASE_URL || 'https://localhost/api/v1';
const token = __ENV.JWT_TOKEN || '';

export default function () {
  const payload = JSON.stringify({ crop_type: 'Wheat', image_key: `scans/test/${Date.now()}.jpg` });
  const params = { headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' } };
  const r = http.post(`${base}/scans`, payload, params);
  check(r, { 'scan queued': (x) => x.status === 200 || x.status === 201 || x.status === 202 });
}
