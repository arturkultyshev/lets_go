# lets_go
## ✍️ Project Description
Our project is hotel booking. This is a web application created for convenient and reliable hotel reservations. Users can view available hotels, register, log in, leave reviews, book hotels and view their bookings.
### Basic functionality
#### For users
- **Registration/Login**: Users can register in the system or log in to their account.
- **View Hotels**: Users can view a list of available hotels with detailed information about each of them.
- **Hotel reservations**: Registered users can book hotels for certain dates of stay.
- **Viewing Orders**: Users have the ability to view and manage their current and previous orders.
- **Leaving reviews**: Users can share their experience of staying by leaving reviews and ratings of hotels.
#### For the administration
- **Hotel management**: The administrator can add, edit, and delete hotels from the system.
- **User Management**: The ability to view the list of users, lock/unlock accounts, view user information.
- **Review moderation**: The ability to view all reviews, delete inappropriate or offensive reviews.

##
## [![Typing SVG](https://readme-typing-svg.herokuapp.com?color=%=FFFFFF&lines=Booking+REST+API)](https://git.io/typing-svg)
```
POST /hotels - создать отель
GET /hotels/:id - один отель
GET /hotels - все отели
PUT /hotels/:id - обновить информацию по отелю
DELETE /hotels/:id - удалить отель

```
##
## [![Typing SVG](https://readme-typing-svg.herokuapp.com?color=%=FFFFFF&lines=DB+Structure)](https://git.io/typing-svg)
```
Table hotels {
  id serial [pk]
  name varchar(255) [not null]
  country varchar(255) [not null]
  city varchar(255) [not null]
  street varchar(255) [not null]
  rating decimal(3,2)
  capacity int
  cost int
  photo_url text
  additional_info text
}

Table users {
  id serial [pk]
  username varchar(255) [not null]
  password varchar(255) [not null]
  first_name varchar(255)
  last_name varchar(255)
  phone_number varchar(20)
  email varchar(255)
  admin boolean [default: false]
}

Table orders {
  id serial [pk]
  user_id int [ref: > users.id]
  hotel_id int [ref: > hotels.id]
  start_date date
  end_date date
  creation_date date
  additional_info text
}

Table reviews {
  id serial [pk]
  user_id int [ref: > users.id]
  hotel_id int [ref: > hotels.id]
  rating decimal(3,2)
  publication_date date
  comment text
}
```

## 🤝 Our team
- Kultyshev Artur 22B030554
- Yerzhanuly Adil 22B030535
- Yelkin Sergey 22B030534
