## REST API

### Basic

1. GET /
   - response body
     ```json
     Hello World
     ```
2. GET /bar?name={name}
   - response body
     ```json
     Hello {name}
     // If there is no query-string(name), response "Hello Bar".
     ```
3. POST /foo
   - request body
     ```json
     {
       "first_name": "mg",
       "last_name": "kim",
       "email": "wowo0201@gmail.com"
     }
     ```
   - response body
     ```json
     {
       "id": 1,
       "first_name": "mg",
       "last_name": "kim",
       "email": "wowo0201@gmail.com",
       "created_at": "2021-01-27T22:19:40.1843404+09:00"
     }
     // If request body is not user interface structure, response 400 bad request.
     ```

---

### Upload and get static resource

1. GET /uploads
   - upload file
1. GET /static/
   - get index.html (upload form)
1. GET /static/{upload_file_name}
   - get static resouce file in public dicrectory
   - If upload_file_name is not exist on public dicrectory, response 404 not found.

---

### CRUD

1. GET /users
   - response body
     ```json
     [
       {
         "id": 1,
         "first_name": "Myung-gwan22",
         "last_name": "Kim22",
         "email": "cartopia9225@naver.com",
         "created_at": "2021-01-27T22:39:53.1916656+09:00"
       },
       {
         "id": 2,
         "first_name": "Myung-gwan33",
         "last_name": "Kim33",
         "email": "cartopia9335@naver.com",
         "created_at": "2021-01-27T22:40:07.0881323+09:00"
       }
     ]
     // If there are no users, response "No Users" (404 not found).
     ```
2. GET /user/{user_id}
   - response body
     ```json
     {
       "id": 1,
       "first_name": "Myung-gwan22",
       "last_name": "Kim22",
       "email": "cartopia9225@naver.com",
       "created_at": "2021-01-27T22:39:53.1916656+09:00"
     }
     // If user id is not valid or not exist, response "No User Id: {user_id}"  (404 not found).
     ```
3. POST /user
   - request body
     ```json
     {
       "first_name": "mg",
       "last_name": "kim",
       "email": "wowo0201@gmail.com"
     }
     ```
   - response body
     ```json
     {
       "id": 1,
       "first_name": "mg",
       "last_name": "kim",
       "email": "wowo0201@gmail.com",
       "created_at": "2021-01-27T22:19:40.1843404+09:00"
     }
     // If request body is not user interface structure, response 400 bad request
     ```
4. PUT /user

   - request body
     ```json
     {
       "id": 2,
       "first_name": "Myung-gwan (UPDATED)",
       "last_name": "Kim33 (UPDATED)",
       "email": "cartopia9335@naver.com (UPDATED)"
     }
     ```
   - response body
     ```json
     {
       "id": 2,
       "first_name": "Myung-gwan (UPDATED)",
       "last_name": "Kim33 (UPDATED)",
       "email": "cartopia9335@naver.com (UPDATED)",
       "created_at": "2021-01-27T22:19:40.1843404+09:00"
     }
     // If request body is not user interface structure, response 400 bad request
     // If user id is not valid or not exist, response "No User Id: {user_id}"  (404 not found)
     ```

5. DELETE /user/{user_id}
   - response body
   ```json
   Deleted User Id: 1
   // If user id is not valid or not exist, response "No User Id: {user_id}" (404 not found)
   ```
