import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '30s', target: 10 },
    { duration: '0s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
}

const BASE_URL = "https://jsonplaceholder.typicode.com/posts/1"

export default function () {
  const pingRes = http.get(`${BASE_URL}`)
  check(pingRes, {
    'ping is 200': r => r.status === 200
  })
}
