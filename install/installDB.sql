use danmaku;

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `comments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` varchar(128) DEFAULT NULL,
  `content` varchar(1024) DEFAULT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `color` char(16) DEFAULT NULL,
  PRIMARY KEY (`id`)
) CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `problem_set`
--

DROP TABLE IF EXISTS `problem_set`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `problem_set` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `question` varchar(512) DEFAULT NULL,
  `answer1` varchar(32) DEFAULT NULL,
  `answer2` varchar(32) DEFAULT NULL,
  `answer3` varchar(32) DEFAULT NULL,
  `answer4` varchar(32) DEFAULT NULL,
  `correct_answer` char(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_answer`
--

DROP TABLE IF EXISTS `user_answer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_answer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(128) NOT NULL,
  `question_id` int(11) NOT NULL,
  `answer` char(1) DEFAULT NULL,
  `time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;


DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `password` varchar(128) DEFAULT NULL,
  `nickname` varchar(128) DEFAULT NULL,
  `reg_code` varchar(128) NOT NULL,
  `permission` int(11) DEFAULT NULL,
  `enrolled` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `reg_code` (`reg_code`)
) CHARSET=utf8mb4;

