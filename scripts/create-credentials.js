#!/usr/bin/env node

/**
 * Script to create authentication credentials for a user via API
 * Usage: node scripts/create-credentials.js [email] [password] [user_id] [api_base_url]
 */

const https = require('https');
const http = require('http');

const email = process.argv[2] || 'customer1@b2b.local';
const password = process.argv[3] || 'password123';
const userId = process.argv[4] || '550e8400-e29b-41d4-a716-446655440004';
const apiBaseUrl = process.argv[5] || 'http://localhost:8090';

const url = new URL(`${apiBaseUrl}/auth/basic/credentials`);

const requestData = JSON.stringify({
  username: email,
  password: password,
  user_id: userId,
});

const options = {
  hostname: url.hostname,
  port: url.port || (url.protocol === 'https:' ? 443 : 80),
  path: url.pathname,
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Content-Length': Buffer.byteLength(requestData),
  },
};

console.log('Creating credentials for user:', email);
console.log('User ID:', userId);
console.log('API Base URL:', apiBaseUrl);
console.log('');

const client = url.protocol === 'https:' ? https : http;

const req = client.request(options, (res) => {
  let data = '';

  res.on('data', (chunk) => {
    data += chunk;
  });

  res.on('end', () => {
    if (res.statusCode === 201) {
      console.log('✓ Credentials created successfully!');
      console.log('Response:', data);
      console.log('');
      console.log('You can now login with:');
      console.log(`  Email: ${email}`);
      console.log(`  Password: ${password}`);
    } else if (res.statusCode === 400 || res.statusCode === 403) {
      console.log(`✗ Failed to create credentials (HTTP ${res.statusCode})`);
      console.log('Response:', data);
      console.log('');
      console.log('Note: This endpoint may require authentication or a registration token.');
      console.log('If you have an admin token, you can use it in the Authorization header.');
    } else {
      console.log(`✗ Unexpected response (HTTP ${res.statusCode})`);
      console.log('Response:', data);
    }
  });
});

req.on('error', (error) => {
  console.error('✗ Request failed:', error.message);
  process.exit(1);
});

req.write(requestData);
req.end();


