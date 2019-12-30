
use danmaku;

select users.nickname, content, time, color from comments
    left join users
        on comments.user = users.reg_code;


select *
from users
where enrolled = 1 and nickname like '%%';


select nickname, count(nickname)
from comments
  left join users
    on comments.user = users.reg_code
where nickname != ''
group by nickname
order by count(nickname) desc;

select citer_user.nickname, cited_user.nickname, count(*)
from comments
  inner join users as cited_user
    on content
      like concat('%', cited_user.nickname, '%')
  left join users as citer_user
    on citer_user.reg_code = comments.user
group by cited_user.nickname, citer_user.nickname;

delete from users where nickname = '';

select content, nickname, time, color
from comments
  left join users
    on comments.user = users.reg_code
where nickname = '';

select content, nickname, time, color
from comments
  left join users
    on comments.user = users.reg_code
where content like '%%';


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
WHERE nickname = '';

# question correct answer rate ranking
SELECT question                                                          as 题目,
       answer1                                                           as A,
       count(IF(answer = 'A', 1, null)) / total_answer_count             as A选择率,
       answer2                                                           as B,
       count(IF(answer = 'B', 1, null)) / total_answer_count             as B选择率,
       answer3                                                           as C,
       count(IF(answer = 'C', 1, null)) / total_answer_count             as C选择率,
       answer4                                                           as D,
       count(IF(answer = 'D', 1, null)) / total_answer_count             as D选择率,
       correct_answer                                                    as 答案,
       count(IF(answer = correct_answer, 1, null)) / total_answer_count  as 正答率,
       total_answer_count                                                as 作答人数
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
GROUP BY problem_set.id;



# incorrect option ranking
SELECT question as 问题, count(answer) as 错答人数, answer as 错答选项,
       CASE
         WHEN answer = 'A' THEN answer1
         WHEN answer = 'B' THEN answer2
         WHEN answer = 'C' THEN answer3
         WHEN answer = 'D' THEN answer4
       END as 选项内容
FROM user_answer
  LEFT OUTER JOIN problem_set
    ON user_answer.question_id = problem_set.id
  LEFT JOIN users
    ON user_answer.user_id = users.reg_code
WHERE user_answer.answer != problem_set.correct_answer
GROUP BY question_id, answer
ORDER BY count(*) DESC ;

