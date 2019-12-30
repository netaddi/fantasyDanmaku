
async function getJSON(jsonUrl) {
    return new Promise(function (resolve) {
        let request = new XMLHttpRequest();
        request.open('GET', jsonUrl, true);
        request.onload = async function() {
            if (this.status >= 200 && this.status < 400) {
                resolve(JSON.parse(this.response));
            }
        };
        request.send();
    });
}

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

class BlackBoardController {

    constructor(){
        this.blackboard = document.getElementById('blackboard');
        this.board = document.getElementById('board');
        this.blackBoardCSSHeight = '84%';
        this.userList = [];
        this.prizeDrawing = false;
        this.prizeInterval = null;
    }

    openBlackboard() {
        this.blackboard.style.height = this.blackBoardCSSHeight;
        this.board.style.display = 'block';
    }

    closeBlackboard() {
        this.blackboard.style.height = '0';
        this.board.style.display = 'none';
    }

    async initializePrizeDraw(){
        this.userList = await getJSON(userListUrl);
        this.openBlackboard();
        document.getElementById('lottery').style.display = 'block';
    }

    closeBlackBoard() {
        document.getElementById('lottery').style.display = 'none';
        document.getElementById('ranking').style.display = 'none';
        this.closeBlackboard();
    }

    async openCommentRanking(){
        const rankingList = await getJSON(commentRankingUrl);
        const rankingDiv = document.getElementById('ranking-table');
        const displayCount = 15;
        this.openBlackboard();

        rankingList.slice(0, displayCount).forEach(item => {
            rankingDiv.innerHTML += `<tr>
                    <td>${item.nickname}</td>
                    <td>${item.count}</td>
                </tr>`
        });
        document.getElementById('ranking').style.display = 'block';

    }

    drawPrize() {
        if (this.prizeDrawing){
            this.prizeDrawing = false;
            clearInterval(this.prizeInterval);
        } else {
            this.prizeDrawing = true;
            this.setPrizeInterval();
        }
    }

    setPrizeInterval(){
        let userList = this.userList;
        let userListLength = userList.length;
        let waitMS = 50;
        this.prizeInterval = setInterval(function (){
            let id = getRandomInt(0, userListLength - 1);
            document.getElementById('prizeRegCode').innerHTML = userList[id].regCode.toString();
            document.getElementById('prizeNickname').innerHTML = userList[id].nickname.toString();
        }, waitMS);
    }
}

let blackboardController = new BlackBoardController();

class QuestionDisplay {
    constructor() {
        this.problemId = 'problem';
        this.answerIdPrefix = 'answer';
        this.countdownId = 'countdown';
        this.questionId = 'question-div';
        this.rankingId = 'standing';

        this.problemDiv = document.getElementById(this.problemId);
        this.answerDivs = [];
        for (let i = 1; i <= 4; ++i){
            this.answerDivs.push(document.getElementById(this.answerIdPrefix + i));
        }
        this.countdownDiv = document.getElementById(this.countdownId);
        this.questionDiv = document.getElementById(this.questionId);
        this.rankingDiv = document.getElementById(this.rankingId)
    }

    processMessage(jsonMessage){
        switch (jsonMessage.QuestionOperation){
            case "prepare":
                this.initialize();
                break;
            case "ranking":
                this.stop();
                break;
            case "updateQuestion":
                this.updateQuestion(jsonMessage.Question, jsonMessage.Answers, jsonMessage.TimeLeft);
                break;
            case "updateCountdown":
                this.updateCountdown(jsonMessage.TimeLeft);
                break;
            case "end":
                this.end();
                break;
        }
    }

    initialize() {
        blackboardController.openBlackboard();
        this.questionDiv.style.display = 'block';
    }

    async stop() {
        this.questionDiv.style.display = 'none';
        await QuestionDisplay.displayRanking();
    }


    updateQuestion(question, answers, timeLeft) {
        this.problemDiv.innerHTML = question;
        for (let i = 0; i <= 3; ++i){
            this.answerDivs[i].innerHTML = answers[i];
        }
        this.updateCountdown(timeLeft)
    }

    updateCountdown(timeLeft) {
        this.countdownDiv.innerHTML = timeLeft
    }

    static async displayRanking(){
        const questionResult = await getJSON(questionResultUrl);
        document.getElementById('standing').style.display = 'block';
        if (questionResult.length > 20) {
            questionResult.length = 20;
        }
        for (let userResult of questionResult){
            let thArray = [];
            let trDiv = document.createElement('tr');
            for (let i = 0; i < 4; i++ ) {
                thArray.push(document.createElement('th'));
            }
            thArray[0].innerHTML = userResult.UserId;
            thArray[1].innerHTML = userResult.Nickname;
            thArray[2].innerHTML = userResult.CorrectCount;
            thArray[3].innerHTML = userResult.Penalty;
            for (let el of thArray){
                trDiv.appendChild(el);
            }
            document.getElementById('standing-table').appendChild(trDiv);
        }
    }

    end(){
        this.rankingDiv.style.display = 'none';
        blackboardController.closeBlackboard();
    }

}

const questionDisplaying = new QuestionDisplay();