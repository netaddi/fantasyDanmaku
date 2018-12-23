package requestProcessors

import (
	"danmakuBackend/danmakuLib"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const optionCount = 4
const countdownSecond = 1
const questionMessageType = "question"

type Problem struct {
	Id int
	Question string
	Answers [optionCount]string
	correctAnswer string
}

type QuestionBasicOperation struct {
	MessageType string
	QuestionOperation string
}

type QuestionUpdateOperation struct {
	MessageType string
	QuestionOperation string
	Question string
	Answers [optionCount]string
	TimeLeft int
}

type CountdownUpdateOperation struct {
	MessageType string
	QuestionOperation string
	TimeLeft int
}

type PersonResult struct {
	UserId string
	Nickname string
	CorrectCount string
	Penalty string
}

type QuestionUserPair struct {
	QuestionId int
	UserId string
}

func (pair QuestionUserPair) String() string {
	return pair.UserId + "_" + strconv.Itoa(pair.QuestionId)
}

type SingleAnswerRecord struct {
	Answer string
	Penalty int64
}

var problemSet []Problem
var currentProblemId int
var currentProblemIndex int
var currentProblemStartTime int64
var answeringStarted bool
//var answerRecord map[string]bool
var answerRecords sync.Map
var questionInitialized = false

func getMillisecondTime() int64 {
	return time.Now().UnixNano() / 1000000
}

// load all the problems from database to var problemSet
func InitializeProblemSet(){

	if questionInitialized {
		return
	}
	config := danmakuLib.GetConfig()

	db, _ := sql.Open("mysql", config.DBsource)
	defer db.Close()

	dbQuery := "SELECT id, question, answer1, answer2, answer3, answer4, correct_answer FROM problem_set;"
	rows, _ := db.Query(dbQuery)
	defer rows.Close()
	var thisProblem Problem
	for rows.Next(){
		//_ := rows.Scan(&thisProblem.Id,
		_ = rows.Scan(&thisProblem.Id,
			&thisProblem.Question,
			&thisProblem.Answers[0],
			&thisProblem.Answers[1],
			&thisProblem.Answers[2],
			&thisProblem.Answers[3],
			&thisProblem.correctAnswer,
		)
		problemSet = append(problemSet, thisProblem)
	}
	questionInitialized = true
}

func prepareAnswering(){
	answeringStarted = false
	//answerRecord = make(map[string]bool)
	InitializeProblemSet()
	initializeMessage := &QuestionBasicOperation{
		MessageType: questionMessageType,
		QuestionOperation: "prepare",
	}
	Frontend.SendMessage(danmakuLib.GetJSON(initializeMessage))
}

func sendRanking(){
	initializeMessage := &QuestionBasicOperation{
		MessageType: questionMessageType,
		QuestionOperation: "ranking",
	}
	Frontend.SendMessage(danmakuLib.GetJSON(initializeMessage))
}

func endAnswering(){
	initializeMessage := &QuestionBasicOperation{
		MessageType: questionMessageType,
		QuestionOperation: "end",
	}
	Frontend.SendMessage(danmakuLib.GetJSON(initializeMessage))
}

func startAnswering(){
	//InitializeProblemSet()
	if answeringStarted {
		return
	}
	answeringStarted = true
	for index, currentProblem := range problemSet{

		secondLeft := countdownSecond
		questionMessage := &QuestionUpdateOperation{
			MessageType:       questionMessageType,
			QuestionOperation: "updateQuestion",
			Question:          currentProblem.Question,
			Answers:           currentProblem.Answers,
			TimeLeft:          secondLeft,
		}

		Frontend.SendMessage(danmakuLib.GetJSON(questionMessage))
		currentProblemId = currentProblem.Id
		currentProblemStartTime = getMillisecondTime()
		currentProblemIndex = index + 1
		for secondLeft >= 0{
			secondLeft--
			time.Sleep(time.Second)
			countdownMessage := &CountdownUpdateOperation {
				MessageType:       questionMessageType,
				QuestionOperation: "updateCountdown",
				TimeLeft:          secondLeft,
			}
			Frontend.SendMessage(danmakuLib.GetJSON(countdownMessage))
		}

	}

	answeringStarted = false
	endMessage := &QuestionUpdateOperation {
		MessageType:       questionMessageType,
		QuestionOperation: "updateQuestion",
		Question:          "答题结束！正在写入数据库……",
	}
	Frontend.SendMessage(danmakuLib.GetJSON(endMessage))

	dumpAnswerRecords(answerRecords)

	answeringStarted = false
	endMessage.Question = "数据库写入完成！"
	Frontend.SendMessage(danmakuLib.GetJSON(endMessage))

}

func dumpAnswerRecords(answerRecords sync.Map) {

	config := danmakuLib.GetConfig()
	db, _ := sql.Open("mysql", config.DBsource)
	defer db.Close()

	resetStmt, _ := db.Prepare("DELETE FROM user_answer WHERE 1=1;")
	_, _ = resetStmt.Exec()
	resetStmt.Close()

	stmt, _ := db.Prepare(`INSERT INTO user_answer 
  						(user_id, question_id, answer, time) VALUES (?, ?, ?, ?);`)
	defer stmt.Close()

	answerRecords.Range(func(key, value interface{}) bool {

		questionId, userId := key.(QuestionUserPair).QuestionId, key.(QuestionUserPair).UserId
		answer, penalty := value.(SingleAnswerRecord).Answer, value.(SingleAnswerRecord).Penalty

		result, _ := stmt.Exec(userId, questionId, answer, penalty)
		println(questionId, userId, answer, penalty)
		_, _ = result.RowsAffected()

		return true
	})
}

func ProcessAnswering(w http.ResponseWriter, r * http.Request){
	_ = r.ParseForm()
	danmakuLib.LogHTTPRequest(r)
	session := danmakuLib.GetSession(r, w)

	userId := session.Values["user"]

	if userId == nil {
		danmakuLib.DenyRequest(w, "请先登录再答题<a href=\\\"login.html\\\">点我登陆</a>")
		return
	}

	if !answeringStarted {
		danmakuLib.DenyRequest(w, "答题还没开始或已经结束，点我也没用哦！")
		return
	}

	questionUserPair := &QuestionUserPair{
		QuestionId: currentProblemId,
		UserId: userId.(string),
	}

	thisAnswer := r.Form.Get("answer")
	thisTime := getMillisecondTime() - currentProblemStartTime

	answerRecord := &SingleAnswerRecord{
		Answer: thisAnswer,
		Penalty: thisTime,
	}

	//_, questionAnswered := answerRecords.Load()
	answerRecords.Store(*questionUserPair, *answerRecord)

	danmakuLib.DenyRequest(w, "成功回答第" + strconv.Itoa(currentProblemIndex) + "题")

}

func GetQuestionResult(w http.ResponseWriter, r * http.Request){
	danmakuLib.LogHTTPRequest(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	config := danmakuLib.GetConfig()
	db, _ := sql.Open("mysql", config.DBsource)
	defer db.Close()

	dbQuery := `SELECT user_id, nickname, count(user_id), sum(time)
					FROM user_answer
						LEFT JOIN problem_set ON user_answer.question_id = problem_set.id
						LEFT JOIN users ON user_answer.user_id = users.reg_code
							WHERE user_answer.answer = problem_set.correct_answer
					GROUP BY user_id
					ORDER BY count(user_id) DESC , sum(time) ASC;`
	rows, _ := db.Query(dbQuery)

	defer rows.Close()
	results := make([]PersonResult, 0)
	for rows.Next(){
		var userResult PersonResult
		_ = rows.Scan(&userResult.UserId, &userResult.Nickname, &userResult.CorrectCount, &userResult.Penalty)
		results = append(results, userResult)
	}
	jsonData, _ := json.Marshal(results)
	_, _ = io.WriteString(w, string(jsonData))

}