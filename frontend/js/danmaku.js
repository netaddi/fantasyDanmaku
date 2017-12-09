
const danmakuDivId = 'danmaku';
const logDivId = 'log';
const logDiv = document.getElementById(logDivId);

var windowWidth = window.innerWidth;

// rails for danmaku.
var rails = [0];

function pushIntoRail(){
    var i = 0;
    for (; i < rails.length; i++){
        if (!rails[i]){
            rails[i] = 1;
            return i;
        }
    }
    rails.push(1);
    return i;
}

function releaseRail(rail) {
    rails[rail] = 0;
}

function moveDanmaku(danmakuItem, rail){
    const freeLength = window.innerWidth / 10 ;
    const totalLength = window.innerWidth + danmakuItem.offsetWidth;
    var freeRailMoveLength = danmakuItem.offsetWidth + freeLength;
    var movedLength = 0;
    var railFreed = false;
    var timer = setInterval(function () {
        movedLength += speed;
        danmakuItem.style.transform = 'translateX(-' + movedLength + 'px)';
        if(!railFreed && movedLength > freeRailMoveLength){
            railFreed = true;
            releaseRail(rail);
        }
        if(movedLength > totalLength){
            clearInterval(timer);
            danmakuItem.remove();
        }
    }, 10);

}

function generateDanmaku(jsonMessage) {
    var danmakuDiv = document.getElementById(danmakuDivId);
    var newHTMLNode = document.createElement("div");
    newHTMLNode.classList.add("comment");
    newHTMLNode.innerHTML = jsonMessage['Text'];
    newHTMLNode.style.color = jsonMessage['Color'];
    newHTMLNode.style.left = innerWidth + 'px';
    newHTMLNode.style.fontSize = defaultSize + 'px';
    var thisRail = pushIntoRail();
    newHTMLNode.style.top = defaultSize * thisRail + 'px';
    danmakuDiv.appendChild(newHTMLNode);
    moveDanmaku(newHTMLNode, thisRail);
}

function openLog() {
    logDiv.style.display = 'block';
}
function closeLog() {
    logDiv.style.display = 'none';
}

function objectMap(obj, call){
    Object.keys(obj).map(
        function (key){
            call(obj[key]);
        })
}

function generateSlides() {
    var xStep = windowWidth;
    var yStep = windowWidth;
    var zStep = 100;
    var x = 0;
    var y = 0;
    var z = 0;
    var impressDiv = document.getElementById('impress');

    showList.map(
        function(el){
            var fileType = el.file.match(/\..+$/)[0];
            var divHtml;
            if (['.jpg', '.png', '.gif'].indexOf(fileType) > -1){
                divHtml = '<div class="step slide" data-x="' + x + '" data-y="' + y + '" data-z="' + z + '">\n' +
                    '    <img src="' + el.file + '">\n' +
                    '</div>';
            }
            if (['.mp4'].indexOf(fileType) > -1){
                divHtml = '<div class="step slide" data-x="' + x + '" data-y="' + y + '" data-z="' + z + '">\n' +
                    '                <video controls="controls">\n' +
                    '                    <source src="' + el.file + '" type="video/mp4" />\n' +
                    '                </video>' +
                    '</div>';
            }
            impressDiv.innerHTML += divHtml;
            x += xStep;
            y += yStep;
            z += zStep;
        })
}

function frontInit(){
    generateSlides();
    objectMap(document.getElementsByTagName('img'), function setWidth(el) {
        el.style.width = windowWidth + 'px';
    });
    objectMap(document.getElementsByTagName('video'), function setWidth(el) {
        el.width = windowWidth;
    });
}