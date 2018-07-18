# fantasyDanmaku

A real time user comment display system for live performance.

# Basic structure

This frontend is divided into to parts: The user frontend and performance frontend.

The user frontend provides basic functions (login, register, send danmaku, etc) for users,
and provide other functions (show recent user list, ban user, etc) for administrators.

The performance frontend includes three modules:
 - Playing the background video / images
 - Playing danmuaku
 - A blackboard for displaying other applications (prize drawing, statistics, question standing, etc)


# deployment

## database

To deploy it on the server, first you have to set up the database (MySQL).
You may find the schemes for database (the .sql file) at `/install/installDB.sql`.

Table `users` stores the information of users. 
To create a new user, create a new null tuple on this table and fill the attribute `reg_code` with a random serial number.
if a new users knows the serail number, he is able to register and specify his own username and password.
This registration will fill the rest of attribute of the new record.


## backend configuration

You have to specify the backend configurations.

The configuration file is in json format,
and you may find it at `/src/danmakuBackend/config/config.json`.

The backend configuration is mainly about database.
Besides, the configuration also includes the port to run the app,
and a token used to indicate the successful establishment of websocket connection.

## frontend configuration

You have to specify the performance frontend configuration.

At performance frontend, there are two configuration files.
One is located at `/frontend/js/config.js`.
This config file is mainly about server information.
Another is located at `/frontend/js/showList.js`.
This configuration file is about the performance information.
The administrator use this file to specify each show and their background image / video file.
For show with separate background image and background music,
controller has to play the background music manually.

## compile it

use this command :
```
export GOPATH=`pwd`
```
at the root dir of the project.

Then, go to `/src/danmakuBackend` folder, and run these commands:
```
go get github.com/go-sql-driver/mysql
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/gorilla/websocket
go build
```

and finally, run the compiled executable file, and open performance frontend to connect.


TODO:
- question answering and standings
- danmaku statistics
- color pick up
