
# CRUD Mongo Test

Technical test for Software Developer job.

The project consists on an application wich has a CRUD for users and a CRUD for posts wich are vinculated to a user.
Also the application has authentication (register and log in) and protection to it endpoints by JWT.

## Run Locally

Clone the project.

```bash
  git clone https://github.com/RamiroCuenca/crud-mongo-test.git
```

Go to the project directory.

```bash
  cd crud-mongo-test
```

Build images from application and database.

```bash
  docker-compose build
```

Run the generated images.

```bash
  docker-compose up
```

Run the unit tests.

```bash
  go test ./... -cover 
```

  
## API Reference

### Users CRUD

#### Register / Sign In a new user

```http
  POST localhost:8080/users/register
```

| Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | `string` | **Required** |
| `password` | `string` | **Required** |

Returns a JWT (located on the headers) wich is necessary to access to access protected endpoints.


#### Log In with an existing user

```http
  POST localhost:8080/users/login
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username` | `string` | **Required** |
| `password` | `string` | **Required** |

Returns a JWT (located on the headers) wich is necessary to access to access protected endpoints.

  
#### Update password from an existing user

```http
  PUT localhost:8080/users/updatebyid?{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `password` | `string` | **Required** - New password|

Endpoint protected with JWT. To access it should send JWT through headers as "Authorization".
  

#### Delete an existing user

```http
  DELETE localhost:8080/users/deletebyid?{id}
```

Endpoint protected with JWT. To access it should send JWT through headers as "Authorization".


### Posts CRUD

#### Create a new post

```http
  POST localhost:8080/posts/create
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `user_id` | `string` | **Required**|
| `title` | `string` | **Required**|
| `description` | `string` | **Required**|

Endpoint protected with JWT. To access it should send JWT through headers as "Authorization".
  

#### Get a post by it ObjectId

```http
  GET localhost:8080/posts/getbyid?{id}
```

Endpoint protected with JWT. To access it should send JWT through headers as "Authorization".
  
#### Get all post posted by a user

```http
  GET localhost:8080/posts/getallfromuserid?{id}
```

Endpoint protected with JWT. To access it should send JWT through headers as "Authorization".

#### Update title and description from an existing post

```http
  PUT localhost:8080/posts/updatebyid?{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `title` | `string` | **Required**|
| `description` | `string` | **Required**|

Endpoint protected with JWT. To access it should send JWT through headers as "Authorization".

#### Delete an existing post

```http
  DELETE localhost:8080/posts/deletebyid?{id}
```

Endpoint protected with JWT. To access it should send JWT through headers as "Authorization".
  
## Tech Stack

**Language:**
- Go

**Packages:** 
- go.mongodb.org/mongo-driver
- gorilla/mux
- golang-jwt/jwt/v4


**Database:**
- MongoDB

**Others:**
- Docker Compose

  
## Author

- [@RamiroCuenca](https://www.linkedin.com/in/ramiro-cuenca-salinas-749a2020a/)

  
## Appendix

- I did not add encryption to passwords as i wasn't allowed to use any package external to the Standard Library (With a few exceptions as mongo drivers or gorilla/mux). 
 
