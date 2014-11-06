#### Yamr

Yamr is a chat application built with [golang](http://golang.org).
You can demo it by going to [yamr.net](http://yamr.net).
[Old node.js version moved here](https://github.com/poops/yamr-node).

#### Installation

Yamr uses packages from [npm](https://www.npmjs.org/) and [mysql](http://www.mysql.com). After cloning the repo, install dependencies with:

    npm install

Afterwards, you can use [grunt](http://gruntjs.com/) to compile assets:

    grunt

Create a config.go file with the following

    package main

    import (
      "code.google.com/p/gorilla/sessions"
    )

    // Change "change me" to a secret key used for authenticating sessions
    var store = sessions.NewCookieStore([]byte("change me"))

    // Database connection string
    var DSN = "root:@/yamr?charset=utf8"


Run the following to start the server

    ./yamr

#### Nginx

Here is a sample nginx configuration

    upstream app {
      server 127.0.0.1:8000;
    }

    server {
      listen 80;
      server_name yamr.net;
      root /path/to/compiled/public;

      location / {
        index index.html;

        if (-f $request_filename) {
          break;
        }

        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $http_host;
        proxy_set_header X-NginX-Proxy true;
        proxy_redirect off;
        proxy_pass http://app/;
      }

      location ~* ^.+.(jpg|jpeg|gif|css|png|js|ico|txt)$ {
        expires max;
        access_log off;
      }
    }

#### SQL

Here is the SQL for creating the mysql database

    CREATE DATABASE `yamr` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

    CREATE TABLE IF NOT EXISTS guests (
      `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
      `agent` varchar(255) DEFAULT NULL,
      `ip` char(15) DEFAULT NULL,
      `created` datetime DEFAULT NULL,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

    CREATE TABLE IF NOT EXISTS `messages` (
      `id` int(10) NOT NULL AUTO_INCREMENT,
      `room_id` int(10) unsigned NOT NULL,
      `user_id` int(10) unsigned NOT NULL,
      `guest_id` int(10) unsigned NOT NULL,
      `message` varchar(255) NOT NULL,
      `created` datetime NOT NULL,
      PRIMARY KEY (`id`),
      KEY `room_id` (`room_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

    CREATE TABLE IF NOT EXISTS `rooms` (
      `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
      `name` varchar(255) NOT NULL,
      PRIMARY KEY (`id`),
      UNIQUE KEY `name` (`name`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

    CREATE TABLE IF NOT EXISTS `users` (
      `id` int(10) NOT NULL AUTO_INCREMENT,
      `name` varchar(255) NOT NULL,
      `password` varchar(128) NOT NULL,
      `created` datetime NOT NULL,
      `ip` char(15) NOT NULL,
      `agent` varchar(255) NOT NULL,
      `headshot` int(10) unsigned NOT NULL DEFAULT '0',
      `last_login` date DEFAULT NULL,
      PRIMARY KEY (`id`),
      UNIQUE KEY `username` (`name`)
    ) ENGINE=InnoDB  DEFAULT CHARSET=utf8;
