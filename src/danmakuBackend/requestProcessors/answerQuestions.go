package requestProcessors

import (
	"database/sql"
	"danmakuBackend/danmakuLib"
	"time"
	"net/http"
	"strconv"
	"io"
	"encoding/json"
)

const optionCount = 4
const countdownSecond = 10
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

var problemSet []Problem
var currentProblemId int
var currentProblemIndex int
var currentProblemStartTime int64
var answeringStarted bool
var answerRecord map[string]bool
var questionInitialized = false

func getMilisecondTime() int64 {
	return time.Now().UnixNano() / 1000000
}

// load all the prolems from database to var problemSet
func InitializeProlemSet(){

	if questionInitialized {
		return
	}

	config := danmakuLib.GetConfig()

	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database: ", err.Error())
		db.Close()
		return
	}
	defer db.Close()

	dbQuery := "SELECT id, question, answer1, answer2, answer3, answer4, correct_answer FROM problem_set;"
	rows, err := db.Query(dbQuery)
	if err != nil {
		println("failed to query database.: ", err.Error())
	} else {
		defer rows.Close()
		var thisProblem Problem
		for rows.Next(){
			err := rows.Scan(&thisProblem.Id,
							&thisProblem.Question,
							&thisProblem.Answers[0],
							&thisProblem.Answers[1],
							&thisProblem.Answers[2],
							&thisProblem.Answers[3],
							&thisProblem.correctAnswer)

			if err != nil {
				println("question error: ", err.Error())
			}

			problemSet = append(problemSet, thisProblem)
		}
	}
	questionInitialized = true
}

func prepareAnswering(){
	answeringStarted = false
	answerRecord = make(map[string]bool)
	InitializeProlemSet()
	initializeMessage := &QuestionBasicOperation{questionMessageType, "prepare"}
	Frontend.SendMessage(danmakuLib.GetJSON(initializeMessage))
}

func sendRanking(){
	initializeMessage := &QuestionBasicOperation{questionMessageType, "ranking"}
	Frontend.SendMessage(danmakuLib.GetJSON(initializeMessage))
}

func endAnswering(){
	initializeMessage := &QuestionBasicOperation{questionMessageType, "end"}
	Frontend.SendMessage(danmakuLib.GetJSON(initializeMessage))
}

func startAnswering(){
	//InitializeProlemSet()
	if answeringStarted {
		return
	}
	answeringStarted = true
	for index, currentProblem := range problemSet{

		secondLeft := countdownSecond
		questionMEssage := &QuestionUpdateOperation{
			questionMessageType,
			"updateQuestion",
			currentProblem.Question,
			currentProblem.Answers,
			secondLeft }

		Frontend.SendMessage(danmakuLib.GetJSON(questionMEssage))
		currentProblemId = currentProblem.Id
		currentProblemStartTime = getMilisecondTime()
		currentProblemIndex = index + 1
		for secondLeft >= 0{
			secondLeft--
			time.Sleep(time.Second)
			countdownMessage := &CountdownUpdateOperation{
				questionMessageType,
				"updateCountdown",
				secondLeft }
			Frontend.SendMessage(danmakuLib.GetJSON(countdownMessage))
		}

	}
	endMessage := &QuestionUpdateOperation{
		questionMessageType,
		"updateQuestion",
		"答题结束！",
		[optionCount]string{},
		0 }
	Frontend.SendMessage(danmakuLib.GetJSON(endMessage))

	answeringStarted = false

}

func ProcessAnswering(w http.ResponseWriter, r * http.Request){
	r.ParseForm()
	danmakuLib.LogHTTPRequest(r)
	session := danmakuLib.GetSession(r, w)

	username := session.Values["user"]

	if username == nil {
		danmakuLib.DenyRequest(w, "请先登录再答题<a href=\\\"login.html\\\">点我登陆</a>")
		return
	}

	if !answeringStarted {
		danmakuLib.DenyRequest(w, "答题还没开始，点我也没用哦！")
		return
	}

	thisAnswer := r.Form.Get("answer")
	thisTime := getMilisecondTime() - currentProblemStartTime

	config := danmakuLib.GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database.")
		danmakuLib.DenyRequest(w, "failed to connect database.")
		db.Close()
		return
	}
	defer db.Close()

	_, recordFound := answerRecord[username.(string) + "__" + strconv.Itoa(currentProblemId)]
	if recordFound {
		// this user has answered this question before
		stmt, err := db.Prepare("UPDATE user_answer SET answer=? , time=? WHERE user_id=? and question_id=?;")
		//println(thisAnswer, thisTime, username, currentProblemId)
		result, err := stmt.Exec(thisAnswer, thisTime, username, currentProblemId)
		if err != nil {
			danmakuLib.DenyRequest(w, "failed to write database.")
			println("err: ", err.Error())
		}
		affect, err := result.RowsAffected()
		defer stmt.Close()
		if err != nil {
			danmakuLib.DenyRequest(w, "failed to write database.")
			println("err: ", err.Error())
		}
		if affect == 1 {
			danmakuLib.DenyRequest(w, "成功修改第" + strconv.Itoa(currentProblemIndex) + "题的回答")
			//danmakuLib.AcceptRequest(w)
		} else {
			danmakuLib.DenyRequest(w, "failed to write database. no update")
			println("err: no row affected. ")
		}
	} else {
		answerRecord[username.(string) + "__" + strconv.Itoa(currentProblemId)] = true

		stmt, err := db.Prepare("INSERT INTO user_answer (user_id, question_id, answer, time) " +
										"VALUES (?, ?, ?, ?);;")
		defer stmt.Close()
		if err != nil {
			println("error: ", err.Error())
			danmakuLib.DenyRequest(w, "database error. ")
		}
		result, err := stmt.Exec(username, currentProblemId, thisAnswer, thisTime)
		if err != nil {
			println("error: ", err.Error())
			danmakuLib.DenyRequest(w, "database error. ")
		}
		affect, err := result.RowsAffected()
		if err != nil {
			println("error: ", err.Error())
			danmakuLib.DenyRequest(w, "database error. ")
		}
		if affect == 1{
			danmakuLib.DenyRequest(w, "成功回答第" + strconv.Itoa(currentProblemIndex) + "题")
			//danmakuLib.AcceptRequest(w)
		} else {
			danmakuLib.DenyRequest(w, "数据库写入失败")
		}
	}

	if err != nil {
		println("failed to query database.: ", err.Error())
		danmakuLib.DenyRequest(w, "failed to query database")
		//db.Close()
		return
	}
}

func GetQuestionResult(w http.ResponseWriter, r * http.Request){
	danmakuLib.LogHTTPRequest(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	config := danmakuLib.GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database: ", err.Error())
		io.WriteString(w, "{}")
		db.Close()
		return
	}
	defer db.Close()

	dbQuery := `SELECT user_id, nickname, count(user_id), sum(time)
					FROM user_answer
						LEFT JOIN problem_set ON user_answer.question_id = problem_set.id
						LEFT JOIN users ON user_answer.user_id = users.reg_code
							WHERE user_answer.answer = problem_set.correct_answer
					GROUP BY user_id
					ORDER BY count(user_id) DESC , sum(time) ASC;`
	rows, err := db.Query(dbQuery)
	if err != nil {
		println("failed to query database.: ", err.Error())
		io.WriteString(w, "{}")
	} else {
		defer rows.Close()
		results := make([]PersonResult, 0)
		for rows.Next(){
			var userResult PersonResult
			rows.Scan(&userResult.UserId, &userResult.Nickname, &userResult.CorrectCount, &userResult.Penalty)
			results = append(results, userResult)
		}
		jsonData, _ := json.Marshal(results)
		io.WriteString(w, string(jsonData))
	}

}