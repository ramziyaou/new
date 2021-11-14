homework-10
Нужно ваш сервис по майнингу переместить в Docker
Создать Feeder для БД – (насыщение БД пользователями и кошельками)
Система должна полностью и самостоятельно подниматься на любой машине с докером после docker compose up
Предоставить http запросы для проверки (через Postman collections)
Работа оценивается зачет/незачет

# WEB SERVICE FOR CRYPTO MINING

## Instructions
### Step 1. 
Run in command line:

```
docker compose up -d
``` 
where -d is for detached mode and optional

| ⚠️WARNING: Wait until you see the following and only then go to localhost: |
| --- |
```
homework-10-ramziyaou-app-1  | 2021/11/14 23:30:40 stderr: 2021/11/14 23:30:40 Connection was successful!!
homework-10-ramziyaou-app-1  | 2021/11/14 23:30:40 stderr: 2021/11/14 23:30:40 Table migrated
homework-10-ramziyaou-app-1  | 2021/11/14 23:30:40 stderr: 2021/11/14 23:30:40 Starting the HTTP server on port 8080
```

### Step 2. 
Use ```Crypto Collection.postman_collection``` and run the collection on POSTMAN.
<br /><br /><br />

## Description 
<br /><br />
The web service supports the following endpoints for ``` http:localhost:8080/```

The following endpoint supports GET and POST methods for authorized users to retrieve data on user under given ID or create a new user, accepting numeric symbols (0-9) for ID.

| ⚠️WARNING: Note that ID of 0 is not supported by the program. |
| --- |

The corresponding web service methods are GetUser and SaveUser.
```
/app/user/$id
```
The endpoint below accepts GET and POST methods to get wallet amount if such exists or register a new one under given title. 

It helps implement GetWallet and SaveWallet methods of the web service, respectively.
```
/app/wallet/$name
```
Finally, the web service provides mining methods via HTTP OPTIONS method.

The following implement StartMining and StopMining methods to perform respective operations.
```
/app/wallet/$name/start
```
```
/app/wallet/$name/stop
```

To start the web service, run 
```
go run .
```
Next, using curl (see ```curl.txt```) or POSTMAN, authorize as admin under username "Ramziya" with password "1234".

Respective HTTP methods are specified in parentheses.

To get info on admin (GET), go to path below
```
http:localhost:8080/app/user/0
```
This will list username, ID and wallets (names only) for admin.

To create a new user (POST), follow path below providing params for username and password
```
http:localhost:8080/app/user/1?username=USERNAME&password=PASSWORD
```
This will register a user USERNAME under ID 1. ID cannot clash with those of existing users.
To save a new wallet, an authorized user must provide new wallet name, different from the ones they already have, if any, in the following format:
```
http:localhost:8080/app/wallet/WALLETNAME
```
WALLETNAME accepts alphabetical characters (a-zA-z) only. 


To start or stop mining (OPTIONS), a registered user should select one of their existing wallets and trigger the corresponding method by:
```
http:localhost:8080/app/wallet/WALLETNAME/start
```
or
```
http:localhost:8080/app/wallet/WALLETNAME/stop
```

In the meantime, a user may check amount of their particular wallet via (GET):
```
http:localhost:8080/app/WALLETNAME
```
This returns 404 if wallet is not found.

# Creating a database

Using CMD, go to  
```
mysql -u root -p
```
Type your password, which you should also update on Line 62 of ``` main.go ```

To create a database, run
```
create database test_two;
```
Similarly, run the following command if you want to renew your database:
```
drop database test_two;
```
This will delete the database, you will have to recreate it by using the previous command

To get database contents, run
```
use test_two;
```
followed by
```
show tables;
```
to display all tables contained in the database

The output for our program should be:
```
mysql> show tables;
+--------------------+
| Tables_in_test_two |
+--------------------+
| crypto_wallets     |
| start_stop_checks  |
| users              |
+--------------------+
3 rows in set (0.00 sec)

mysql>
```
Three useful commands to check table contents after performing http requests:
```
select * from users;
select * from crypto_wallets;
select * from start_stop_checks;
```
These will display all users in the ```users``` table, all cryptowallets in the ```crypto_wallets``` table, and mining statuses for all existing wallets in the ```start_stop_checks``` table

Have fun :)