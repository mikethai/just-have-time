# just-have-time

## Purpose
The goal of this project is to make the [KKBOX](https://www.kkbox.com/) has more fun feature in the APP.
We create the story card that contains music songs with the tag.
And those cards are shared by your followers who will talk about their daily story all in the story card by songs.

## :memo: Build a RESTful API in Go
- Language: Golang.
- Framework: Fiber
- Database: Cloud SQL (Postgres), GCP FireStore
- Server: Cloud Run ＆ Cloud Scheduler

Find the full article to build this API [here](https://documenter.getpostman.com/view/15733862/2s8Z6u6bDA)

![image server-architecture](https://just-have-time-tcj2k7lwbq-de.a.run.app/img/server-architecture.png)
## Get It Started -

- Create a database just-have-time in your Postgres local instance
- Rename `.env.exmaple` to `.env`.
- Add in the database user and password for your postgres instance containing the database just-have-time.
- In the root folder run `make dev`.
- Get the API Server at http://localhost

## What the KKBOX API We using?
- https://api.kkbox.com/v1.1/charts
- https://api.kkbox.com/v1.1/search
- https://api.kkbox.com/v1.1/tracks
- https://api-listen-with.kkbox.com.tw/v3/users
- https://api-listen-with.kkbox.com.tw/v3/users/{{userID}}/following
- https://api-listen-with.kkbox.com.tw/v3/encrypt
- https://api-listen-with.kkbox.com.tw/v3/decrypt