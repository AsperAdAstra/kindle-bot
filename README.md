# TgBot Send To Kindle

This is a telegram bot that sends documents to your Kindle [email address](https://www.amazon.com/sendtokindle/email).

## This bot is meant to be used as a self-hosted solution.
Bot can be hosted anywhere, it does not rely on webhooks.

Bot requires your own Bot and your user ID in order to work properly. 
As bots are always public, it is a part of a bot functionality to accept incoming
messages from a single user only.

## Configuration
For reference, [see example.env](./example.env) file.

### Bot
```BOT_TOKEN```
You need to create your own bot and get its token. 

```USER_ID```
You need to figure out and set your user id. Everyone else will be ignored.

### Kindle
[Send To Kindle by Email](https://www.amazon.com/sendtokindle/email)

```MAIL_FROM```
Email address that will be used as a sender. It must be whitelisted in Amazon.

```MAIL_TO```
Kindle email address. You can check it in your Kindle app.


### SMTP transport
```SMTP_HOST``` 
SMTP host. For Gmail it is ```smtp.gmail.com```.

```SMTP_USERNAME```
SMTP username. Usually it is an email address.

```SMTP_PASSWORD```
SMTP password. For Gmail you must create [App Password](https://support.google.com/accounts/answer/185833?hl=en).

```SMTP_PORT=587```
Default port for SMTP (587).

## Running
`go run cmd/bot/main.go`

## Bugs, issues, feature requests
Please use GitHub issues to report bugs, issues or feature requests.