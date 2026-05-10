import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    upload: { executor: 'constant-arrival-rate', rate: 10, timeUnit: '1s', duration: '2m', preAllocatedVUs: 20, maxVUs: 100 },
  },
};

const base = __ENV.BASE_URL || 'https://localhost/api/v1';
const token = __ENV.JWT_TOKEN || '';

export default function () {
  const params = { headers: { Authorization: `Bearer ${token}` } };
  const r = http.get(`${base}/scans/upload-url?content_type=image/jpeg`, params);
  check(r, { 'upload-url 200': (x) => x.status === 200 });
  sleep(0.2);
}
