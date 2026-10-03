import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: "1m",
  thresholds: {
    'http_req_duration{name:ping}': ['p(95)<500'],
    'http_req_duration{name:createuser}': ['p(95)<500'],
    'http_req_duration{name:login}': ['p(95)<500'],
    'http_req_duration{name:product}': ['p(95)<500'],
    'http_req_duration{name:getproduct}': ['p(95)<500'],
    'http_req_duration{name:orders}': ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
}

const BASE_URL = "http://localhost:3000/api/v1"



export default function () {
  const uniqueEmail = `user-${__VU}-${__ITER}@test.com`;
  const uniqueName = `Widget-${__VU}-${__ITER}`;
  const pingRes = http.get(`${BASE_URL}/ping`, { tags: { name: 'ping' } })
  check(pingRes, {
    'ping is 200': r => r.status === 200
  })

  const createUserRes = http.post(`${BASE_URL}/create_user`,
    JSON.stringify({ email: uniqueEmail, password: 'testpassword' }),
    { headers: { 'Content-Type': 'application/json' }, tags: { name: 'createuser' } }
  )
  check(createUserRes, { 'create user is 200': r => r.status === 201 })
  const loginRes = http.post(
    `${BASE_URL}/login`,
    JSON.stringify({ email: uniqueEmail, password: 'testpassword' }),
    { headers: { 'Content-Type': 'application/json' }, tags: { name: 'login' } }
  )
  check(loginRes, { 'login is 200': r => r.status === 200 })
  const token = JSON.parse(loginRes.body).token;
  const prodRes = http.post(`${BASE_URL}/add_product`,
    JSON.stringify({ product: { name: uniqueName, price: 9.99, quantity: 100 } }),
    { headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` }, tags: { name: 'product' } }
  )
  check(prodRes, { 'add product is 200': r => r.status === 201 })


  const authHeaders = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  }

  const productRes = http.get(`${BASE_URL}/products`, authHeaders, { tags: { name: 'getproduct' } });
  check(productRes, { 'products is 200': r => r.status === 200 })

  const orderRes = http.post(
    `${BASE_URL}/orders`,
    JSON.stringify({
      customer_id: '1',
      items: [{ product_id: '3', quantity: 2, price: 10.62 }],
    }),
    authHeaders, { tags: { name: 'orders' } }
  )
  check(orderRes, { 'orders is 200': r => r.status === 201 })

  sleep(1);
}
