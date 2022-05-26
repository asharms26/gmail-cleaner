# Gmail Cleaner

A simple service written in Golang that cleans out your Gmail Inbox. 

## How-To

- Pull the repo
- Run the main.go driver
- Make a POST call to http://localhost:10000/flush

Request Payload

token: Google OAuth Token
keywords: Keywords in the "From" field of the email in question

`{
    "token":"<GOOGLE_OAUTH_TOKEN_HERE>",
    "keywords": ["HBO", "Herman Miller", "WeWork"]
}`


## Notes
You need an OAuth token that authorizes you to make the API call to the Gmail servers.

- Create API credentials in the Google Cloud Developer Console
- Use Postman to get OAuth token following this format to authenticate through a brower
