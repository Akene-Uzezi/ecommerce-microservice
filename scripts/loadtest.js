import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '5s', target: 10 },
    { duration: '5s', target: 10 },
    { duration: '0s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
}

const BASE_URL = "http://localhost:3001/api/v1"

export function setup() {
  const createUserRes = http.post(`${BASE_URL}/create_user`,
    JSON.stringify({ email: 'test@test.com', password: 'testpassword' }),
    { headers: { 'Content-Type': 'application/json' } }
  )
  check(createUserRes, { 'creat user is 200': r => r.status === 201 })
  const loginRes = http.post(
    `${BASE_URL}/login`,
    JSON.stringify({ email: 'test@test.com', password: 'testpassword' }),
    { headers: { 'Content-Type': 'application/json' } }
  )
  check(loginRes, { 'login is 200': r => r.status === 200 })

  const token = JSON.parse(loginRes.body).token;
  return { token: token }
}

export default function (data) {
  const pingRes = http.get(`${BASE_URL}/ping`)
  check(pingRes, {
    'ping is 200': r => r.status === 200
  })

  const authHeaders = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${data.token}`
    }
  }

  const productRes = http.get(`${BASE_URL}/products`, authHeaders);
  check(productRes, { 'products is 200': r => r.status === 200 })

  const orderRes = http.post(
    `${BASE_URL}/orders`,
    JSON.stringify({
      customer_id: '1',
      items: [{ product_id: '3', quantity: 2, price: 10.62 }],
    }),
    authHeaders
  )
  check(orderRes, { 'orders is 200': r => r.status === 200 })

  sleep(1);
}
