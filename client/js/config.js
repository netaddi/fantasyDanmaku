const server = 'db.zjufantasy.club';
const port = '8082';


// const wsUrl = 'ws://' + server + ':' + port + '/ws';
const wsUrl = `ws://${server}:${port}/ws`;
const httpUrlBase = `http://${server}:${port}/`;
const userListUrl = httpUrlBase + 'getUserList';
const questionResultUrl = httpUrlBase + 'getQuestionResult';
const commentRankingUrl = httpUrlBase + 'getCommentRanking';


const config = {
    server : server,
    port : port,
    wsUrl : wsUrl
};


// use this token as a signal of successfully connected websocket.
// have to be the same as the one on server side.
const wsToken = '__Danmaku_WS_Connected_Token___';
