
use danmaku;

select * from comments order by time desc ;

select *
from users
where enrolled = 1 and nickname like '%巴%';


select nickname, count(nickname)
from comments
  left join users
    on comments.user = users.reg_code
where nickname != ''
group by nickname
order by count(nickname) desc;

select nickname, count(nickname)
from comments
  inner join users
    on content
      like concat('%', users.nickname, '%')
group by nickname;


select content, nickname, time, color
from comments
  left join users
    on comments.user = users.reg_code
where nickname = '无敌高_神';

select content, nickname, time, color
from comments
  left join users
    on comments.user = users.reg_code
where content like '%林渊%';


# answer ranking
SELECT user_id, nickname, count(user_id), sum(time)
FROM user_answer
  LEFT JOIN problem_set
    ON user_answer.question_id = problem_set.id
  LEFT OUTER JOIN users
    ON user_answer.user_id = users.reg_code
WHERE user_answer.answer = problem_set.correct_answer
GROUP BY user_id
ORDER BY count(user_id) DESC , sum(time) ASC;


# query user answer
SELECT question, answer1, answer2, answer3, answer4, answer, correct_answer
FROM user_answer
  LEFT OUTER JOIN problem_set
    ON user_answer.question_id = problem_set.id
  LEFT JOIN users
    ON user_answer.user_id = users.reg_code
WHERE nickname = '上帝';

# question correct answer rate ranking
SELECT question, answer1, answer2, answer3, answer4, answer, count(*) / total_answer_count, count(*), total_answer_count
FROM user_answer
  LEFT OUTER JOIN problem_set
    ON user_answer.question_id = problem_set.id
  LEFT JOIN users
    ON user_answer.user_id = users.reg_code
  LEFT JOIN (
    SELECT problem_set.id as id, count(*) as total_answer_count
    FROM user_answer
      LEFT JOIN problem_set
        ON problem_set.id = user_answer.question_id
    GROUP BY problem_set.id
  ) total_answer_table
    ON user_answer.question_id = total_answer_table.id
WHERE user_answer.answer = problem_set.correct_answer
GROUP BY problem_set.id, total_answer_count;



# incorrect option ranking
SELECT question, count(answer), answer,
       CASE
         WHEN answer = 'A' THEN answer1
         WHEN answer = 'B' THEN answer2
         WHEN answer = 'C' THEN answer3
         WHEN answer = 'D' THEN answer4
       END as term
FROM user_answer
  LEFT OUTER JOIN problem_set
    ON user_answer.question_id = problem_set.id
  LEFT JOIN users
    ON user_answer.user_id = users.reg_code
WHERE user_answer.answer != problem_set.correct_answer
GROUP BY question_id, answer
ORDER BY count(*) DESC ;

