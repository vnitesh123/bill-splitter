# Bill Splitter

This service is to track the shared expenses among friends which works as a bill splitter, similar to splitwise with minimal supported operations build on Golang.

## Requirements
Dependencies are manages by go modules

Go verstion - v1.15 or higher
DataBase - MySQL

## Execution

Build the project -  

```
go build -o splitwise main.go
```

Run the project - 

```
./splitwise
``` 

## APIs (operations)

### User Register - POST /register
example request body 
{
    "username":"nitesh",
    "phone":"9999999999",
    "email":"abc@xyz.com"
}

### Add Expense - POST /expense
example request body :
{
    "name":"house rent",
    "createdBy":5,
    "paidBy":{
        "6":6.00
    },
    "paidFor":{
        "5":1,
        "6":2,
        "7":3
    },
    "amount":6.00,
    "splitType":"EXACT"
} 

### Get Balances
- get balances of how much you owe other and how much you get from others
- GET /{userId}/balances

### Get Expenses
- get the list of expenses
- GET /expenses

## Database and Table Schemas 
Tables - Users, Balances, Expenses
Note : add your DB creds in repository/connection.go

Create following tables

Users table :
CREATE TABLE `Users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL DEFAULT '',
  `addedOn` datetime NOT NULL,
  `email` varchar(200) NOT NULL DEFAULT '',
  `phone` varchar(20) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8

Balances Table :
CREATE TABLE `Balances` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `paidBy` int(11) NOT NULL DEFAULT '0',
  `paidFor` int(11) NOT NULL DEFAULT '0',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `expenseId` int(11) NOT NULL DEFAULT '0',
  `addedOn` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8

Expenses Table :
CREATE TABLE `Expenses` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL DEFAULT '',
  `split` varchar(40) NOT NULL DEFAULT '',
  `amount` decimal(12,2) NOT NULL DEFAULT '0.00',
  `paidBy` text NOT NULL,
  `paidFor` text NOT NULL,
  `createdOn` datetime NOT NULL,
  `createdBy` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8
