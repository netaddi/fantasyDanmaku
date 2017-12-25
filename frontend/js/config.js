const server = '127.0.0.1';
const port = '8081';

const wsUrl = 'ws://' + server + ':' + port + '/ws';
const regUrl = '//' +  server + ':' + port + '/reg';
const loginUrl = '//' +  server + ':' + port + '/login';
const sendUrl = '//' +  server + ':' + port + '/send';

const config = {server : server,
                port : port,
                wsUrl : wsUrl};


// use this token as a signal of successfully connected websocket.
// have to be the same as the one on server side.
const wsToken = '__Danmaku_WS_Connected_Token___';
const speed = 1;
const defaultSize = 36;
const danmakuLineContiueThreshold = 100;