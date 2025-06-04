import http from 'k6/http';
import { check } from 'k6';

export const options = {
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 500, // Matches server capacity (~1,770 iters/s from prior tests)
            timeUnit: '1s',
            duration: '30s',
            preAllocatedVUs: 500,
            maxVUs: 1000,
        },
    },
};

export default function () {
    // Define the payload
    const payload = JSON.stringify({
        userID: '842a7be5-59bb-457d-9e59-e76abd0a7044', // Valid UUID
        status: 'confirmed',
        paymentStatus: 'paid',
    });

    // Set headers
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Send POST request with dynamic event ID
    const eventID = '9458e0b3-0463-4e3e-83f5-f6adf7bdff51'; // Valid UUID
    const res = http.post(`http://localhost:8080/events/${eventID}/register`, payload, params);

    // Check response
    const checkResult = check(res, { 'status is 201': (r) => r.status === 201 }); // Expect 201
    if (!checkResult) {
        console.log(`Failed request: status=${res.status}, error=${res.error}`);
    }
}
