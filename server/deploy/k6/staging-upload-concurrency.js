import http from 'k6/http';
import { check } from 'k6';

export const options = { vus: 80, duration: '4m' };
const base = __ENV.BASE_URL || 'https://localhost';
const token = __ENV.JWT_TOKEN || '';

export default function () {
  const r = http.get(`${base}/api/v1/scans/upload-url?content_type=image/jpeg`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  check(r, { 'upload-url status': (x) => x.status === 200 || x.status === 401 });
}
