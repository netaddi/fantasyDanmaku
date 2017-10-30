let server = '127.0.0.1';
let port = '8081';

let wsUrl = 'ws://' + server + ':' + port + '/ws';
let regUrl = '//' +  server + ':' + port + '/reg';
let loginUrl = '//' +  server + ':' + port + '/login';
let sendUrl = '//' +  server + ':' + port + '/send';

let config = {server : server,
                port : port,
                wsUrl : wsUrl};