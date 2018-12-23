const server = 'localhost';
const port = '8081';


// const wsUrl = 'ws://' + server + ':' + port + '/ws';
const wsUrl = `ws://${server}:${port}/ws`;
const httpUrlBase = `http://${server}:${port}/`;
const userListUrl = httpUrlBase + 'getUserList';
const questionResultUrl = httpUrlBase + 'getQuestionResult';


const config = {
    server : server,
    port : port,
    wsUrl : wsUrl
};


// use this token as a signal of successfully connected websocket.
// have to be the same as the one on server side.
const wsToken = '__Danmaku_WS_Connected_Token___';
const speed = 1;
const defaultSize = 48;
const danmakuLineContiueThreshold = 80;