const server = 'newyear.zjufantasy.com';
const port = '8081';

const wsUrl = 'ws://' + server + ':' + port + '/ws';
const regUrl = '//' +  server + ':' + port + '/reg';
const answerUrl = '//' +  server + ':' + port + '/answer';
const loginUrl = '//' +  server + ':' + port + '/login';
const sendUrl = '//' +  server + ':' + port + '/send';
const userListUrl = 'http://' +  server + ':' + port + '/getUserList';
const resultUrl = 'http://' +  server + ':' + port + '/getQuestionResult';


const config = {server : server,
                port : port,
                wsUrl : wsUrl};


// use this token as a signal of successfully connected websocket.
// have to be the same as the one on server side.
const wsToken = '__Danmaku_WS_Connected_Token___';
const speed = 1;
const defaultSize = 48;
const danmakuLineContiueThreshold = 80;