import http from 'k6/http';
import { check } from 'k6';

export const options = { vus: 100, duration: '3m' };
const base = __ENV.BASE_URL || 'https://localhost';

export default function () {
  const payload = JSON.stringify({
    id: `evt_${__VU}_${Date.now()}`,
    event: 'payment.captured',
    payload: { payment: { entity: { id: 'pay_test', order_id: 'order_test', amount: 100, status: 'captured' } } },
  });
  const r = http.post(`${base}/api/v1/payments/webhook`, payload, {
    headers: { 'Content-Type': 'application/json', 'X-Razorpay-Signature': 'invalid' },
  });
  check(r, { 'webhook spoof rejected': (x) => x.status === 401 });
}
