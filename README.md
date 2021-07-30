# truvest-identity-management
An RBAC based Full fledge User Management Service written in Golang. This service has following features:
- Users
	* User Signup
	* Create New User
	* Send Email Link to set password
	* Resend Email link to set password
	* User Login
	* User Login Refresh
	* Reset Password
	* Logged-in User API
	* User Logout
  	* Update User
  	* Delete User
- Roles
	* Create new Roles
	* Add User to Role
	* Add Permission to Roles
	* Get Roles
	* Delete Roles
	* Remove Users from Role
	* Remove Permission from Role
  	* Seed some basic Roles
- Permissions
	* Seed Default Permissions
	* Get Permissions
- Token
	* Generate JWT Token
	* Verify JWT token
- OAuth
	* Multi-Provider OAuth support
	* Initially supported with Google and Github
	* Simple OAuth Redirect UI to demo functionality
	* OAuth Logout
- Handle RBAC
	* Handle multiple RBAC scenarios in the service based upon the Permission a user is tagged to
- Swagger documentation
	* Swagger doc for each Rest API
	* Integrate Swagger doc with jwt token(ApiKeyAuth)
- Dockerize the App
	* Run the PostgreSQL databse
	* Run the app by defining custom environment variables


Open the file docker-compose.yaml and edit the environment variables as per your need. You can run the overall app by running:
```
docker-compose up --build
```

Then you can access the swagger URL by calling https://localhost:9191/swagger/index.html

You can find more details about each API usage in Swagger doc.
Once the Swagger URL is up, hit the Login API using the default System Admin details that has been passed in docker-compose.yaml.
You will receive the JWT token as string. You need to take that jwt token and append with Bearer string and pass it in Authorize option in Swagger doc.
For example, Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJlbWFpbCI6InN5c3RlbUBnbWFpbC5jb20iLCJleHAiOjE2MDgwMjYxNDAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoiYWRtaW4ifQ.-Uzikc-jmW6io51IkuF29xQ9LsKzqLAZsVTlsRfEC4EyfMANg-lZEWubNvvOk3UbliCMf9VUImv1_J7LnLhznw

One sample email template is also being bundled under html folder in case someone wants to try out "Send Email" through SMTP server to alert user about its credentials or "Forget Password". This can be modified as per the usage.

This App is configured with OAuth, presently supporting just Google and GitHub but it can be extended with a pretty wide list from below:

* Amazon
* Apple
* Auth0
* Azure AD
* Battle.net
* Bitbucket
* Box
* Cloud Foundry
* Dailymotion
* Deezer
* DigitalOcean
* Discord
* Dropbox
* Eve Online
* Facebook
* Fitbit
* Gitea
* GitHub
* Gitlab
* Google
* Heroku
* InfluxCloud
* Instagram
* Intercom
* Kakao
* Lastfm
* Linkedin
* LINE
* Mailru
* Meetup
* MicrosoftOnline
* Naver
* Nextcloud
* Okta
* OneDrive
* OpenID Connect (auto discovery)
* Oura
* Paypal
* SalesForce
* Shopify
* Slack
* Soundcloud
* Spotify
* Steam
* Strava
* Stripe
* Tumblr
* Twitch
* Twitter
* Typetalk
* Uber
* VK
* Wepay
* Xero
* Yahoo
* Yammer
* Yandex
