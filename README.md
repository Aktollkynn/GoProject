# GoProject


***__Progress 1️⃣:__***

We created roadmap for our project:
Link to roadmap
https://www.notion.so/98ade335d571433d97e999790c7ad683?v=f54791d06d81453da286384a46619775&pvs=4

***__Progress 2️⃣:__***

In the second progress, we have connected the project to the database. For our project, we use a PostgreSQL database. Then we created models in the project: users, product, order, category, payment, section and so on. 

***__Progress 3️⃣:__***

In the third progress, our team created the main page, the product page, which displays a list of store products from the database, also the registration page.

***__Progress 4️⃣:__***

In fourth progress, our team created registration, users can register in the website and their data will appeared and stored in database
we also have a branch here - https://github.com/Aktollkynn/GoProject/tree/Report4
  
PostgreSQl DB,  `Table: users`
```sql
+------------+------------+--------+-----+------------+----------------------------------+
| Name       | Data type  | Length | Key | Not Null?  |  Default                         |
+------------+------------+--------+-----+------------+----------------------------------+
| id         | integer    |        | PRI | Yes        |nextval('users_id_seq'::regclass) |
| first_name | varchar    | 50     |     | Yes        |                                  |
| last_name  | varchar    | 50     |     | Yes        |                                  |
| email      | varchar    | 355    |     | Yes        |                                  |
| password   | varchar    | 50     |     | Yes        |                                  |
+----------+--------------+------+-----+-------------------+-----------------------------+
```
***__Progress 5️⃣:__*** 


| __How Authorization work__|---|---|---|---| ---|
|     ---        |       ---              |       ---  |      ---                                |    ---       | ---        |
| `/register `   | `/registerauth`        | `/login`   | `/loginauth`                            | `/home_page` | `/logout`  | 
| register first | _user_ is  registered? | login user | checks _email_ & _password_ to correct  | *Welcome!*   | end session|

**Overview:** *In the fifth progression, we created a user login and a product search by name. This way, users who have registered will be able to login to our website. Also on the home page appears information about products directly from the database, namely books; users can search for books by name and see the information.*

- [x] Registeration-> Insert data to DB
- [x] Login System -> Users from DB(if exist, otherwise try again)
- [x] Logout -> After session end return to Login page
- [x] Searching -> Any information
