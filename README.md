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

***__Progress 6️⃣:__*** 

 In the sixth progression, we created validation - firstname, lastname, password, email and hashing password.
 - [x] Validation for Register(Upper,small letters, numbers, min 8 symbols) 
 - [x] Hashing password
 
 ***__Progress 7️⃣:__*** 
 - [x] Filter prices(max, min)
 - [x] Logout session, cookie(without loggin, couldn't open home page even with link)
 - [x] User data in home page(shows after login, and his fname, lname)
 - [x] Filter prices(max, min)
 - [x] Fixed issues

PostgreSQl DB,  `Table: products`
```sql
+-------------+------------+--------+-----+------------+---------+
| Name        | Data type  | Length | Key | Not Null?  | Default |
+-------------+------------+--------+-----+------------+---------+
| id          | integer    |        |     | Yes        |         |
| name        | varchar    | 50     |     | Yes        |         |
| description | varchar    | 355    |     | Yes        |         |
| price       | integer    |        |     | Yes        |         |
+-------------+------------+--------+-----+------------+---------+
```
home_page

<img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/home_page.jpeg?raw=true" width="500px" height="auto" />
products

<img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/items.jpeg?raw=true" width="500px" height="auto" />
search and filer

<img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/search.jpeg?raw=true" width="500px" height="auto" />
searching_results

<img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/searching_results.jpeg?raw=true" width="500px" height="auto" />



***__Progress 8️⃣:__*** 

In the eighth progress, we created user profile, and add profile menu,  edit, update profile functions
 
|  ---                        |      __User profile__  |       ---               |
|     ---                     |       ---              |       ---               |   
| `/profile `                 | `/edit_profile`        | `/update_profile`       | 
| view information about user | change information     | update inside database  | 



 - [x] Nav profile menu (Profile page, setting, logout,)  
            <img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/menu.jpeg?raw=true" width="150px" height="auto" />
 
 - [x] View all information about user(fname, lname, email,)
 - [x] Edit all information about user(fname, lname, email, password)
 - [x] Updated scripts( alerts for info such as successfully register, change information)
 - [x] minimalize & clear  code, fixed issues

<img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/profile.jpeg?raw=true" width="500px" height="auto" />
<img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/edir_profile.jpeg?raw=true" width="500px" height="auto" />
<img src="https://github.com/galymzhantolepbergen/galymzhantolepbergen/blob/main/jobs/update_profile.jpeg?raw=true" width="500px" height="auto" />

***__Progress 9:__***
In the ninth progression, we created product detail page and a rating for products. Users can give a rating for each product, and the average product rating will be displayed on the page.

![image](https://user-images.githubusercontent.com/82767082/232326055-c57b759c-3cee-4ec1-b7a8-4ccc1030270f.png)

***__Progress 10:__***
In the tenth progression, we created comments, so users can add comments for products

***__Midterm 2:__***
 We made a page to add and improve comments function, filtering by price and product rating. In addition to comments, the user can give ratings to books.
