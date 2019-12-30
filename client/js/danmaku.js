
const speed = 3.2;
const danmakuMoveInterval = 25;  // ms
const defaultSize = 64;
const maxDanmakuRailCount = 12;
const continuousDanmakuWaitTime = 512;  // ms

const danmakuDivId = 'danmaku';
const logDivId = 'log';
const logDiv = document.getElementById(logDivId);
const danmakuDiv = document.getElementById(danmakuDivId);

const windowWidth = window.innerWidth;
const windowHeight = window.innerHeight;
let connected;

class KeywordFilter{
    constructor() {
        this.bannedKeywordList = [];
    }

    addBannedKeyword(keyword) {
        this.bannedKeywordList.push(keyword)
    }

    checkCommentBanned(comment) {
        return this.bannedKeywordList.reduce(function (previousValue, thisKeyword) {
            return comment.indexOf(thisKeyword) > -1 || previousValue
        }, false)
    }
}

const keywordFilter = new KeywordFilter();

function moveDanmaku(danmakuItem, rail){
    const freeLength = window.innerWidth / 10 ;
    const totalLength = window.innerWidth + danmakuItem.offsetWidth;
    let freeRailMoveLength = danmakuItem.offsetWidth + freeLength;
    let movedLength = 0;
    let railFreed = false;
    let startPosition = danmakuItem.getBoundingClientRect().left;
    let timer = setInterval(function () {
        movedLength += speed;
        // if (movedLength > 100 ) return;
        danmakuItem.style.left = (startPosition - movedLength).toString() + 'px';
        // danmakuItem.style.transform = 'translateX(-' + movedLength + 'px)';
        if(!railFreed && movedLength > freeRailMoveLength){
            railFreed = true;
            danmakuController.releaseRail(rail)
        }
        if(movedLength > totalLength){
            clearInterval(timer);
            danmakuItem.remove();
        }
    }, danmakuMoveInterval);
}


class DanmakuController{
    constructor() {
        this.rails = Array(maxDanmakuRailCount).fill(false);
        this.danmakuQueue = [];
    }



    generateDanmaku(jsonMessage, thisRail) {
        if (!jsonMessage) {
            return;
        }
        let newHTMLNode = document.createElement("div");
        newHTMLNode.classList.add("comment");
        newHTMLNode.innerHTML = jsonMessage['Text'];
        newHTMLNode.style.left = innerWidth + 'px';
        newHTMLNode.style.fontSize = defaultSize + 'px';
        if (/[0-9a-fA-F]{6}/.test(jsonMessage['Color'])){
            newHTMLNode.style.color = '#' + jsonMessage['Color'];
        } else {
            newHTMLNode.style.color = '#FFFFFF';
        }
        newHTMLNode.style.top = defaultSize * thisRail + 'px';
        danmakuDiv.appendChild(newHTMLNode);
        moveDanmaku(newHTMLNode, thisRail);
    }

    processDanmaku(jsonMessage) {
        if (keywordFilter.checkCommentBanned(jsonMessage)){
            return;
        }
        const rail = this.getAvailableRail();
        if (rail > -1) {
            this.generateDanmaku(jsonMessage, rail);
        } else {
            this.danmakuQueue.push(jsonMessage);
        }
    }

    getAvailableRail(){
        for (let i = 0; i < maxDanmakuRailCount; i++){
            if (!this.rails[i]){
                this.rails[i] = true;
                return i;
            }
        }
        // no rail available. push the danmaku into buffer.
        return -1;
    }

    releaseRail(rail) {
        if (this.danmakuQueue.length) {
            setTimeout(() => {
                console.log("buffer: ", this.danmakuQueue);
                const cachedDanmaku = this.danmakuQueue.shift();
                console.log("reading buffer", cachedDanmaku);
                this.generateDanmaku(cachedDanmaku, rail);
            }, continuousDanmakuWaitTime);
        } else {
            this.rails[rail] = false;
        }
    }

}

class PlaybackController{

    static getActiveVideoDom(){
        const activeDom = document.getElementsByClassName('active')[0].children[0];
        console.log(activeDom);
        if (activeDom.tagName === 'VIDEO') {
            return activeDom
        } else {
            return null
        }
    }

    static play() {
        const videoDom = this.getActiveVideoDom();
        console.log(videoDom);
        if (videoDom){
            videoDom.play();
        }
    }
    
    static pause() {
        const videoDom = this.getActiveVideoDom();
        if (videoDom){
            videoDom.pause();
        }
    }
}

const danmakuController = new DanmakuController();
const impressController = impress();


function generateTestDanmaku(text) {
    danmakuController.processDanmaku({
        'Text': text,
        'Color': 'FFFFFF'
    });
}

function generateMultipleTestDanmaku(text, count) {
    for (let i = 0; i < count; i ++) {
        generateTestDanmaku(text);
    }
}

async function processAdminCommand(operation, parameter) {
    switch (operation) {
        case "banKeyword" :
            keywordFilter.addBannedKeyword(parameter);
            break;
        case "openLottery" :
            await blackboardController.initializePrizeDraw();
            break;
        case "close" :
            blackboardController.closeBlackBoard();
            break;
        case "openCommentRanking" :
            await blackboardController.openCommentRanking();
            break;
        case "openLog" :
            displayLog();
            break;
        case "closeLog" :
            hideLog();
            break;
        case "prev":
            impressController.prev();
            break;
        case "next":
            impressController.next();
            PlaybackController.play();
            break;
        case "goto":
            impressController.goto(parseInt(parameter));
            break;
        case "play":
            PlaybackController.play();
            break;
        case "pause":
            PlaybackController.pause();
            break;
        case "refresh":
            location.reload();
            break;
        case "init":
            impressController.swipe()
    }
}


async function processWSMessage(message){
    printLog(message);
    let jsonMessage = JSON.parse(message);
    console.log(jsonMessage);
    switch (jsonMessage.MessageType) {
        case wsToken:
            connected = true;
            return;
        case "adminMessage":
            await processAdminCommand(jsonMessage.AdminOperation, jsonMessage.OperationParameter);
            return;
        case "danmaku":
            danmakuController.processDanmaku(jsonMessage);
            return;
        case "question":
            questionDisplaying.processMessage(jsonMessage);
            return;
    }

}


function objectMap(obj, call){
    Object.keys(obj).map(
        function (key){
            call(obj[key]);
        })
}

function configConnection(conn) {
    conn.onclose = function (evt) {
        printLog("connection closed. Reconnection");
        conn = new WebSocket(config.wsUrl);
        configConnection(conn);
    };
    conn.onmessage = async function (evt) {
        await processWSMessage(evt.data)
    };
}

function generateSlides() {
    let xStep = windowWidth;
    let yStep = windowWidth;
    let zStep = 100;
    let x = 0;
    let y = 0;
    let z = 0;
    let impressDiv = document.getElementById('impress');

    showList.map(
        el => {
            console.log(el);
            let fileType = el.match(/\..{3}$/)[0];
            let showDiv = document.createElement('div');
            showDiv.setAttribute('class', 'step slide');
            showDiv.setAttribute('data-x', x.toString());
            showDiv.setAttribute('data-y', y.toString());
            showDiv.setAttribute('data-z', z.toString());
            impressDiv.appendChild(showDiv);
            if (['.png', '.JPG', '.PNG', '.jpg', '.jpeg'].includes(fileType)) {
                let imgDiv = document.createElement('img');
                imgDiv.setAttribute('src', el);
                showDiv.appendChild(imgDiv);
            }
            if (['.mp4', '.flv'].includes(fileType)) {
                let videoDiv = document.createElement('video');
                videoDiv.setAttribute('controls', 'controls');
                let sourceDiv = document.createElement('source');
                sourceDiv.setAttribute('src', el);
                sourceDiv.setAttribute('type', 'video/mp4');
                videoDiv.appendChild(sourceDiv);
                showDiv.appendChild(videoDiv);
            }
            x += xStep;
            y += yStep;
            z += zStep;
        })
}

function showInit(){
    generateSlides();
    objectMap(document.getElementsByTagName('img'), function setWidth(el) {
        el.style.height = windowHeight + 'px';
    });
    objectMap(document.getElementsByTagName('video'), function setWidth(el) {
        el.height = windowHeight;
    });
}

function displayLog() {
    logDiv.style.display = 'block';
}

function hideLog() {
    logDiv.style.display = 'none';
}

function printLog(message) {
    let item = document.createElement("div");
    item.innerHTML = message;
    logDiv.appendChild(item);
    console.log(message)
}