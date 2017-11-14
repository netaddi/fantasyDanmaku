const server = '127.0.0.1';
const port = '8081';

const wsUrl = 'ws://' + server + ':' + port + '/ws';
const regUrl = '//' +  server + ':' + port + '/reg';
const loginUrl = '//' +  server + ':' + port + '/login';
const sendUrl = '//' +  server + ':' + port + '/send';

const config = {server : server,
                port : port,
                wsUrl : wsUrl};

const wsToken = '__Danmaku_WS_Connected_Token___';