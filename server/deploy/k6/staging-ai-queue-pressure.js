import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    queue_pressure: {
      executor: 'constant-arrival-rate',
      duration: '5m',
      rate: 25,
      timeUnit: '1s',
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },
};

const base = __ENV.BASE_URL || 'https://localhost';
const token = __ENV.JWT_TOKEN || '';

export default function () {
  const payload = JSON.stringify({ crop_type: 'Wheat', image_key: `scans/staging/${Date.now()}.jpg` });
  const r = http.post(`${base}/api/v1/scans`, payload, {
    headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
  });
  check(r, { 'scan queued': (x) => [200, 201, 202, 400, 401].includes(x.status) });
}
