var express = require('express');
var app = express();
var path = require('path');

// viewed at http://localhost:8080
const proxy = require('http-proxy-middleware')
var apiProxy = proxy('/api', {target: 'http://localhost:31112/function/storage'});
app.use(apiProxy)

app.get('/', function(req, res) {
    res.sendFile(path.join(__dirname + '/index.html'));
});

app.listen(8000);