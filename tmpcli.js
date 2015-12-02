#!/usr/bin/env node

// This is solely a test client for the sake of testing out the certs on both sides.
// We setup a TLS connection with a specific CA and client key pair against the go server
// locally, then pipe data from stdin.

var tls = require('tls');

options = {
  host: "127.0.0.1",
  port: 27001,
	key: process.env.KEY,
  cert: process.env.CERT,
  ca: [process.env.CA]
}

var socket = tls.connect(options);

socket.on('secureConnect', function() {
  console.log('client connected',
              socket.authorized ? 'authorized' : 'unauthorized');
  process.stdin.pipe(socket);
  process.stdin.resume();
});

socket.setEncoding('ascii');
socket.on('data', function(data) {
	console.log(data);
});

socket.on('end', function() {
	console.log("THE END");
});
