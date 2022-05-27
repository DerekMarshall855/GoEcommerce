# Description

Small ecommerce project to improve at Go. Project start date: May 11 2022
Project finished, just the backend API made. Endpoint descriptions included below.
Finish date: May 27 2022

## How to run

With docker installed:
run 'docker-compose up -d' to create and run images
run 'go mod tidy' to get packages
run 'go run main.go' to start server

From here you can use any test api query software you prefer (I use postman) to test out the endpoints

# Endpoints

## POST

### /users/signup

    - Allows users to signup, use firstName, lastName, email, password, and phone

### /users/login

    - Login endpoint, use email, password

### /admin/addproduct

    - Allows you to add products (currently anyone can use doesn't require admin account)
    - Use name, price, rating, image

## GET

### /users/productview

    - Returns product view

### /users/search?name='name'

    - Returns product information using 'name' query parameter

### /addtocart?id='id'&userId='userId'

    - Add product to userId cart via product id query parameter

### /removeitem?id='id'&userId='userId'

    - Removes item to userId cart via product id query parameter

### /checkout?userId='userId'

    - Tallys price + checkout for user via user userId in query parameter

### /instantbuy?id='id'&userId='userId'

    - Instant buys it item of productId id for user with userId
