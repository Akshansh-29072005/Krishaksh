import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = { vus: 20, duration: '1m' };
const base = __ENV.BASE_URL || 'https://localhost';

export default function () {
  const res = http.get(`${base}/health`);
  check(res, { 'health 200': (r) => r.status === 200 });
  sleep(1);
}
