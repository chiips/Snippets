# API

This API folder contains sample snippets for developing production-quality RESTful JSON APIs in Go.

## Organization
The API folder is organized into the following files and folders.

### Main.go
Main.go is the entry point for the program and initializes every component of the API from the database to the request rate limiter. The server is set for deployment over HTTP for development purposes but includes, commented out, the necessary adjustments for HTTPS.

### App
The app folder contains the bulk of the API's files.

The server struct is defined for dependency injection in server.go. All handlers hang off the server to access dependencies (e.g., database, logger). Sample handlers are included that relate to user's creating and managing accounts, as well as reading, creating, updating, and deleting posts. The routes.go file defines all routes and corresponding handlers.

The folders contains middleware defined for the server's router, and therefore all requeusts, as well as middleware defined for specific handlers, namely authentication middleware for protected routes. JSON Web Token (JWT) authentication is used. Auth.go includes the code for administering JWTs on successful login.

The folder also contains sample tests for the user and post handlers, supplemented with the test-setup.go file.

### Logs
The logs folder contains a log.go file that creates a new logger using logrus (https://github.com/Sirupsen/logrus) and a log.txt file which can serve as the destination for logs if chosen. Choose to log to a file or the terminal.

### Models
The models folder contains files to define the API's datastore and database methods, as well as to establish a connection with PostgreSQL.

### Private
The private folder is where users' file uploads would be stored as referenced in the app/handlers-users.go EditProfilePhoto function.

## Improvements

This repository would benefit from several improvements.

### Sign-Up
The signup handler currently logs in users too by administering a JWT. This handler should instead send a verification email which redirects users to log in; log-in should then be handled by a separate handler.

This improvement would require the addition of an email client, ideally included in the server struct for the signup handler to easily access.

An email client would also be necessary for handling forgotten passwords with a similar workflow as signing up.

### JWT
The API's JWTs are currently created using an HMAC signing method which requires a single key for both signing and verifying JWTs.

JWT security can be improved by changing the signing method to RSA using private and public keys. With this signing method the middleware code that authenticates JWTs only has access to the public key for verification and not the private key used for administering JWTs. This separation is particularly useful if the API is moved to a microservices architecture and authentication of users is isolated from resource access.

The API currently administers only access JWTs which expire in 5 minutes. Users are forced to log in repeatedly if they wish to access protected routes (e.g., submitting posts or editing their profile).

The API should also administer refresh tokens which serve to refresh access tokens within defined limits of time (e.g., 7 or 30 days). Refresh tokens need not be JWTs—random strings such as UUIDs would suffice—and can be cached for checking when requests come in with expired access tokens. Alternatively access or refresh tokens can be blacklisted by storing in a cache once used; the cache can then be cleared repeatedly as tokens reach their expiration dates. A limit with this approach is the introduction of server-side state which is contrary to RESTful principles.

Another option is introducing standard sessions since cookies are already used. Simply create UUIDs for each user, store them in cookies, and on every request check the session cookies against the cache. Again the limit here is the introduction of server-side state.

### Separate Authentication and Resources
Currently authentication and resource access are all part of the same components in one API. In a microservices architecture the authentication of users can be separated from resource access: a user connects to an authentication API, receive an access token, and uses that to request resources from the resource API. This architecture provides useful modularity for scaling and adding new functionality. The additional work, however, may not be worth the cost for certain apps.

The current API can benefit nonetheless by separating authentication and resource access into distinct packages.

### Refactor
Certain checks are repeated across several handlers, such as user and paramater checks. These checks could be factored out into distinct functions to avoid repetition. Any changes to these checks could be made in one function and update all handlers. This change, however, would need to be balanced with ease of understanding and use. Different handlers need different checks so being explicit within each handler may be ideal.

### Overall Extension of the API
As mentioned above, the API is not complete. Many more handlers are needed for a full app, and tests are needed for every handler.
