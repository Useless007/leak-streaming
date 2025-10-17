import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s'
};

const BASE_URL = __ENV.API_BASE_URL || 'http://localhost:8080';

export default function () {
  const response = http.post(`${BASE_URL}/movies/sample-movie/playback-token`);
  check(response, {
    'status is 200': (res) => res.status === 200
  });
  sleep(1);
}
