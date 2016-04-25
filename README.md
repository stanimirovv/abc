# ABC

An anonymous image board (think 4chan/2chan) which doesn't store any personal information and the threads and attachments are deleted after a certain period. The design allows several image boards to be hosted on the same backend.

## Installation

### Packages

In order to run ABC you need to have the following packages installed:

PostgreSQL
Go(lang) (Version 1.5+)
Although the names of the packages may vary from distribution to distribution for ubuntu/debian you can go with
```
 sudo apt-get install postgresql postgresql-contrib
 sudo apt-get install golang
 ```

### Golang Libraries

 All go libraries come with this source, you don't have to install anything extra.


## Build
 Inside the project root, run ```build_server.sh``` It will compile and buld the binary. After that it will run the automated tests. If the build fails it means that something in the environment (such as a database is not configured)

## Creating a database
 Run the 	```create_database.sh``` create database script to create a database. You should change the name of the database, the db user and it's password!

## Run

Run the script ``` run_server.sh```  to run the compiled binary. Before that you should check if the configuration of the run_server.sh is correct. It is  done via environment variables. The default ones should work if you have created the database with the default settings and your default ports are free.

## Configuration

todo

## Code structure

todo

## REST API
todo

## REST API Examples

todo
