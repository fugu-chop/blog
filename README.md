# README.md

## Dean's Blog
This is a repo to house code and instructions for my personal blog. 

This project is intended for me to learn a number of things:

1. Spinning up a production standard Golang server

2. Structuring an app and writing more Go code

3. Practice containerising the app

4. Deploying this into production for the internet to consume

5. Get more practice writing (I'd like to capture my thoughts and make them understandable)

6. Creating a nice-ish frontend (maybe)

### Design Decisions
Based on the above use cases, I could probably get away with just making a static site using Hugo or something.

However, I'm interested in learning 'the hard way', so my original intention was to use as few libraries outside of the standard Go library as possible.

#### Routing
I immediately gave in to using an external framework for routing after evaluating the [routing updates](https://go.dev/blog/routing-enhancements) in the `net/http` package in Go 1.22. 

While the addition of HTTP verbs and extraction of parts of the path are wonderful, applying _middleware_ is still a pain - AFAICT I still have to manually apply middleware to each and every `http.Handler`. That is a lot of error-prone repetition. 

I've settled on using [`chi`](https://github.com/go-chi/chi), which offers a much more granular way to apply middleware to `http.Handler`s.
