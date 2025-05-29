# Social Network API v1 Documentation

## Base URL: /api/v1/
API Endpoints Overview
WebSocket

    GET /api/v1/ws
    Establish a WebSocket connection for real-time features like chats.
    Requirements:

        User must be logged in

        Additional requirements TBD

## Authentication (/api/v1/auth)
Method	|Endpoint	Description	Requirements / Payload
| Method | Endpoint         | Description         | Requirements / Payload                                                                                     |
|--------|------------------|---------------------|------------------------------------------------------------------------------------------------------------|
| GET    | /auth            | Check login status  | Returns `{ "status": true }`                                                                                |
| POST   | /auth/register   | Register new user   | Not logged in. Payload JSON:<br> `{ "username": "", "email": "", "birthdate": "YYYY-MM-DD",<br> "fname": "", "lname": "", "password": "", "gender": "" }` |
| POST   | /auth/login      | User login          | Not logged in. Payload: <br>`{ "login": "username/email",`<br>` "password": "something" }`                                     |
| POST   | /auth/logout     | Logout user         | No payload, just POST to log out                                                                             |
## GET Endpoints (/api/v1/get/)

| Method | Endpoint             | Description              | Notes / Requirements  |
|--------|----------------------|--------------------------|----------------------|
| GET    | /get/profile         | Retrieve user profile    | To be defined         |
| GET    | /get/posts           | Fetch posts              | To be defined         |
| GET    | /get/categories      | List all categories      | To be defined         |
| GET    | /get/conversations   | List user conversations  | To be defined         |
| GET    | /get/messages        | Get messages from a conversation | To be defined |

### GET /get/ ( `profile` ,`posts` ,`comments` ,`categories` ,`conversations` ,`messages` )
- profile TODO
    >{"id":"","display_name":"azer","avatar": "TODO:","description": "azer",<br>"is_public": true/false,"is_person": true/false,"created_at":"RFC3339"}

    - logged in
        > add adds "isliked" values either (null/true/false) or( 0/1/2) or (undefined/true/false) 
        > - TODO:decide
- posts 
    >{"post_id" : 123,"creator_id" : 123,"group_id" : 123,<br>"creator_name" : "azer","title": "todo","body" : "todo","image_path" : "todo",<br>"privacy" : "public/follower/private","creationtime" : "RFC3339"}
    
    - if logged in

        > adds "isliked" values either (null/true/false) or (0/1/2) or (undefined/true/false) 
        > - TODO:decide
- comments
    > { "comment_id" : 123,"creator_id" : 123,<br>"creator_name" : "azer" , "body" : "todo" ,<br>"image_path" : "todo","creationtime" : "RFC3339" }
- categories
    > [ {"id" : 123 , "name":"azer" } , {"id" : 123 , "name":"azer" } , ... ]
- conversations
    > TODO decide access managment <br> probably :
    >> { profile_id : 123 , profile_name : "azer", is_group : true/false }
- messages
    > TODO decide <br> probably :
    >> { message_id : 132 , sender_id : 123 , receiver_id: 123 , sender_name : "azer" , receiver_name: "azer" , message : "azer" , is_read : true/false , sent_at : "RFC3339" }
- notifications
    > { id : 123 , title : "azer" , "body" : TODO to be decided}

## SET Endpoints (/api/v1/set/)

| Method | Endpoint          | Description            | Notes / Requirements  |
|--------|-------------------|------------------------|----------------------|
| POST   | /set/profile      | Update user profile    | To be defined         |
| POST   | /set/posts        | Create or update posts | To be defined         |
| POST   | /set/comments     | Add a comment          | To be defined         |
| POST   | /set/follow       | Follow or unfollow a user | To be defined       |

- idk u suggest