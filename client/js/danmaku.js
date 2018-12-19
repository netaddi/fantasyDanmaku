
const danmakuDivId = 'danmaku';
const logDivId = 'log';
const logDiv = document.getElementById(logDivId);
// const log = document.getElementById("log");

const windowWidth = window.innerWidth;
const windowHeight = window.innerHeight;
let connected;
// rails for danmaku.
// let rails = [0];


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

class DanmakuRail{
    constructor() {
        this.rails = [0];
    }

    pushIntoRail(){
        let i = 0;
        for (; i < this.rails.length; i++){
            if (!this.rails[i]){
                this.rails[i] = 1;
                return i;
            }
        }
        this.rails.push(1);
        return i;
    }

    releaseRail(rail) {
        this.rails[rail] = 0;
    }

}


let rails = new DanmakuRail();
let keywordFilter = new KeywordFilter();


function processAdminCommand(operation, parameter) {
    switch (operation) {
        case "banKeyword" :
            keywordFilter.addBannedKeyword(parameter);
            break;
        case "openLottery" :
            initializePrizeDraw();
            break;
        case "closeLottery" :
            closePrizeDraw();
            break;
        case "openLog" :
            displayLog();
            break;
        case "closeLog" :
            hideLog();
            break;
    }
}

function processWSMessage(message){
    printLog(message);
    let jsonMessage = JSON.parse(message);
    console.log(jsonMessage);
    switch (jsonMessage.MessageType) {
        case wsToken:
            connected = true;
            return;
        case "adminMessage":
            processAdminCommand(jsonMessage.AdminOperation, jsonMessage.OperationParameter);
            return;
        case "danmaku":
            generateDanmaku(jsonMessage);
            return;
        case "question":
            questionDisplaying.processMessage(jsonMessage);
            return;
    }

}

function moveDanmaku(danmakuItem, rail){
    const freeLength = window.innerWidth / 10 ;
    const totalLength = window.innerWidth + danmakuItem.offsetWidth;
    let freeRailMoveLength = danmakuItem.offsetWidth + freeLength;
    let movedLength = 0;
    let railFreed = false;
    let timer = setInterval(function () {
        movedLength += speed;
        danmakuItem.style.transform = 'translateX(-' + movedLength + 'px)';
        if(!railFreed && movedLength > freeRailMoveLength){
            railFreed = true;
            // releaseRail(rail);
            rails.releaseRail(rail)
        }
        if(movedLength > totalLength){
            clearInterval(timer);
            danmakuItem.remove();
        }
    }, 10);

}

function generateDanmaku(jsonMessage) {
    let danmakuDiv = document.getElementById(danmakuDivId);
    let newHTMLNode = document.createElement("div");

    if (keywordFilter.checkCommentBanned(jsonMessage['Text'])) {
        return;
    }

    newHTMLNode.classList.add("comment");
    newHTMLNode.innerHTML = jsonMessage['Text'];
    newHTMLNode.style.left = innerWidth + 'px';
    newHTMLNode.style.fontSize = defaultSize + 'px';
    if (/#[0-9a-fA-F]{6}/.test(jsonMessage['Color'])){
        newHTMLNode.style.color = jsonMessage['Color'];
    } else {
        newHTMLNode.style.color = '#FFFFFF';
    }

    // let thisRail = pushIntoRail();
    let thisRail = rails.pushIntoRail();
    newHTMLNode.style.top = defaultSize * thisRail + 'px';
    danmakuDiv.appendChild(newHTMLNode);
    moveDanmaku(newHTMLNode, thisRail);
}


function objectMap(obj, call){
    Object.keys(obj).map(
        function (key){
            call(obj[key]);
        })
}

function configConnection(conn) {
    conn.onclose = function (evt) {
        printLog("connection closed.");
        alert("connection closed.");
        conn = new WebSocket(config.wsUrl);
        configConnection(conn);
    };
    conn.onmessage = function (evt) {
        processWSMessage(evt.data)
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
        function(el){
            let fileType = el.file.match(/\..+$/)[0];
            let showDiv = document.createElement('div');
            showDiv.setAttribute('class', 'step slide');
            showDiv.setAttribute('data-x', x.toString());
            showDiv.setAttribute('data-y', y.toString());
            showDiv.setAttribute('data-z', z.toString());
            impressDiv.appendChild(showDiv);
            if (1){
                let imgDiv = document.createElement('img');
                imgDiv.setAttribute('src', el.file);
                showDiv.appendChild(imgDiv);
            }
            if (['.mp4'].indexOf(fileType) > -1){
                let videoDiv = document.createElement('video');
                videoDiv.setAttribute('controls', 'controls');
                let sourceDiv = document.createElement('source');
                sourceDiv.setAttribute('src', el.file);
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
        // el.style.width = windowWidth + 'px';
        el.style.height = windowHeight + 'px';
        // el.style.width = 'auto';
    });
    objectMap(document.getElementsByTagName('video'), function setWidth(el) {
        // el.width = windowWidth;
        el.height = windowHeight;
    });
}

function displayLog() {
    logDiv.style.display = 'block';
}

function hideLog() {
    logDiv.style.display = 'none';
}

function clearLog() {
    logDiv.innerHTML = '';
}

function printLog(message) {
    let item = document.createElement("div");
    item.innerHTML = message;
    logDiv.appendChild(item);
}