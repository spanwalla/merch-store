import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

export const options = {
    scenarios: {
        constant_rps: {
            executor: 'constant-arrival-rate',
            rate: 1000,
            timeUnit: '1s',
            duration: '5m',
            preAllocatedVUs: 200,
            maxVUs: 1000,
        },
    },
}

// const baseUrl = 'http://localhost:8080/api'
const baseUrl = 'http://host.docker.internal:8080/api'
// const baseUrl = 'http://90.156.159.123:8080/api'

const users = new SharedArray('users', function () {
    let array = [];
    for (let i = 0; i < 1000; i++) {
        array.push({username: `stress_test_user_${i}`, password: "1Test!Password49"});
    }
    return array;
})

const tokens = {};

function getToken(user) {
    if (tokens[user.username]) {
        return tokens[user.username];
    }

    let loginRes = http.post(`${baseUrl}/auth`, JSON.stringify(user), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(loginRes, { 'logged in successfully': (res) => res.status === 200 });

    let token = loginRes.json('token');
    if (token) {
        tokens[user.username] = token;
    }
    return token;
}

export default function () {
    let user = users[Math.floor(Math.random() * users.length)];
    let transferBody = {
        toUser: users[Math.floor(Math.random() * users.length)].username,
        amount: 1,
    };
    let token = getToken(user);
    if (!token) return;

    let params = {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': `application/json`,
        },
    };

    let res = http.get(`${baseUrl}/buy/socks`, params);
    check(res, { 'status is not 500': (r) => r.status !== 500 });
    sleep(1);

    res = http.post(`${baseUrl}/sendCoin`, JSON.stringify(transferBody), params);
    check(res, { 'status is not 500': (r) => r.status !== 500 });
    sleep(1);

    res = http.get(`${baseUrl}/info`, params);
    check(res, { 'status is not 500': (r) => r.status !== 500 });
    sleep(1);
}