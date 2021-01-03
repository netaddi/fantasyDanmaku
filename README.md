# fantasyDanmaku

A real time user comment display system for live performance.

# Basic structure

This frontend is divided into to parts: The user frontend and performance frontend.

The user frontend provides basic functions (login, register, send danmaku, etc) for users,
and provide other functions (show recent user list, ban user, etc) for administrators.

The performance frontend includes three modules:
 - Play the background video / images
 - Play danmuaku
 - A blackboard for displaying other applications (prize drawing, statistics, question standing, etc)


# deployment

## database

To deploy it on the server, first you have to set up the database (MySQL).
You may find the schemes for database (the .sql file) at `/install/installDB.sql`.

At the user frontend, the system is not opened for users to register.
Instead, you are supposed to generate serial numbers randomly,
and add these random serial numbers to attribute `reg_code` of table `users`.

For users of the system, each of them should have one serial number,
and register the system with their unique serial numbers,
along with their self-specified nicknames and passwords.

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
cd server
export GOPATH=`pwd`
```
at the root dir of the project.

Then, go to `src/danmakuBackend` folder, and run these commands:
```
go get github.com/go-sql-driver/mysql
go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/gorilla/websocket
go build
```

and finally, run the compiled executable file, and open performance frontend to connect.

### RegEx for matching problem
`(.+?)\s?.+\n(.+)\s(.+)\s(.+)\s(.+)\n(.)`
into
`('$1', '$2', '$3', '$4', '$5', '$6'),`

TODO:
- admin page: user comment list, ban user
- blackboard: display user comment rankings

temp command
```
PlaybackController.play()
```