# Requirements

- requires go 1.1+
- runs on Windows, and Linux, maybe on OS X

# First steps

### .env

Create file called .env that will contain the credentials (not to be commited), modify the example .env.example and rename it to .env

# Available actions
| todo  | action   |
|---|---|
| downloadformarkov  | Launch the indexing of some articles for generating random shit with markov chains |
| tweetmarkov  | Tweet some markov generated text |
| downloadforretweet  | Launch the indexing of some tweets for retweeting |
| tweetretweet  | Retweet the previously downloaded tweets |
| tweetemojis  | Tweet a few emojis |
| contestplay  | Search for contests and plays |
	
# Use realize to develop or just go build

```
go get github.com/oxequa/realize
```

# Cronify/log all the things

### /etc/logrotate.d/contestator
```
/var/www/logs/contestator.log {
        daily
        missingok
        rotate 14
        compress
        delaycompress
        notifempty
        create 640 www-data www-data
        sharedscripts
}
```

### /etc/cron.d/contestator
```
5 5 * * *   www-data /var/www/contestator/contestator --todo downloadformarkov	 >> /var/www/logs/contestator.log 2>&1
10 6 * * *  www-data /var/www/contestator/contestator --todo downloadforretweet  >> /var/www/logs/contestator.log 2>&1
10 6 * * *  www-data /var/www/contestator/contestator --todo tweetretweet        >> /var/www/logs/contestator.log 2>&1
0 9 * * *   www-data /var/www/contestator/contestator --todo tweetmarkov         >> /var/www/logs/contestator.log 2>&1
30 13 * * * www-data /var/www/contestator/contestator --todo tweetemojis         >> /var/www/logs/contestator.log 2>&1
0 10 * * *  www-data /var/www/contestator/contestator --todo contestplay         >> /var/www/logs/contestator.log 2>&1
```









