# Random programmig joke slack bot


This slack bot will send a random joke into a slack channel according to the configured time interval

Fetches the programming joke from https://v2.jokeapi.dev and publish in your slack channel

## Configuration

Create .env file into your .env folder
Example:

```
SLACK_BOT_TOKEN=<your slack authentication token>
SLACK_CHANNEL_ID=<your slack channel ID C**********>

# uses linux crontab format
# timing examples
# 0 12 1 * * - on 1st day of month
# * * * * *  - every minute
# 0 12 * * *  - noon lauch
# 0 0 * * 1,2 - Monday and Tuesday midnight
# */5 * * * * -every five min"

CRON=*/5 * * * *
SKIP_BEFORE=9:30
SKIP_AFTER=17:30
```

The SKIP_BEFORE, and SKIP_AFTER fields are not mandatory, if not set then it will run 0->24.
Setting the values to 0 also result sending messages 0-24
You can set individually any of them, the other will be ignored. For example if you set SKIP_BEFORE=9 and you don't set SKIP_AFTER then jokes will be delivered 9->23:59


## Make targets

```
make build
make run
run-background
```


With you run the "make run-background", the code will be compiled and will resign in the backround.
If you would like to kill it:

```
ps -ax | grep jokebot

Get the process ID from the output:
Example:
52358 pts/0    Sl     0:00 ./jokebot
50225 pts/0    S+     0:00 grep --color=auto jokebot

Your id is "./jokebot"
kill 52358
```

## What's coming

- tests
- multiple channels

