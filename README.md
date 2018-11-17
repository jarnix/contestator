# Go

- requires go 1.1+
- runs on Windows, and Linux, maybe on OS X

# First steps

### .env

Create file called .env that will contain the credentials (not to be commited), modify the example .env.example and rename it to .env

### database (mysql)

Create a database, a user with password, and the SELECT, INSERT, UDPATE, DELETE rights on this database.

```
CREATE DATABASE <db-name> CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'db-user'@'YOUR.IP' IDENTIFIED BY 'db-password';
```

Then launch the contestator.sql file to create the required tables.

# Available actions
| todo  | action   |
|---|---|
| downloadforindex  | Launch the indexing of some articles for generating random shit with markov chains for tweeting like a real huma... robot  |

# Use realize to develop

```
go get github.com/oxequa/realize
```








