# README.md

## Dean's Blog

This is a repo to house code and instructions for my personal website. 

This project is intended for me to learn a number of things:

1. Spinning up a non-terrible Golang server

2. Structuring an app and writing more Go code

3. Practice containerising the app

4. Deploying this into production for the internet to consume

5. Get more practice writing (I'd like to capture my thoughts and make them understandable)

6. Get a handle on using some cloud products

### Design Decisions

Based on the above use cases, I could probably get away with just making a static site using Hugo or something.

However, because I'm interested in learning about deployment and cloud products, I'd like to learn and practice those skills. 

My original intention was to use as few libraries outside of the standard Go library as possible. There are exceptions, however:

#### Routing
After evaluating the [routing updates](https://go.dev/blog/routing-enhancements) in the `net/http` package in Go 1.22, I decided using an external library would be a better fit/less tedious.

While the addition of HTTP verbs and extraction of parts of the path are wonderful, applying _middleware_ is still a pain - AFAICT I still have to manually apply middleware to each and every `http.Handler`. That is a lot of error-prone repetition. 

I've settled on using [`chi`](https://github.com/go-chi/chi), which offers a much more granular way to apply middleware to `http.Handler`s.

### Cloud Architecture

After initially looking at AWS, I've decided to go with `fly.io` for simplicity.

#### Why not AWS?
While AWS is very widely used, it seems like it has a bunch of different products that need to be set up and coordinated (i.e. there is just a lot of set up and coordination required).

I had originally decided to use a Fargate task as a serverless setup for the app.

However, there's a bunch of _stuff_ and/or tradeoffs that needs to be made for the app:

I would need to use DynamoDB (noSQL) due to it's 'free' status. I prefer Postgres (relational) because it still suits the blog use case and is easier for me to work with. Using AWS would require me to use a managed service (which is _very_ expensive) or spin up a virtual machine that runs a Postgres docker image. Hassle.

Just the sheer amount of stuff to manage:
- I have to use an application load balancer and likely Route53 just to use a custom domain for the ECS cluster (instead of using another provider like Cloudflare).

- Images for Fargate tasks have to be stored ECR (AFAICT)

- Usage of ALB and Route53 to have a custom domain for the app.

- AWS Certificate Manager for SSL/TLS certificates

Whilst almost all of these products would be free for my use case on the AWS free tier, having to use and 'synchronise' all these products creates monitoring and management hassle.

Also, AWS documentation quality is inconsistent. Sometimes it's a bit shit, quite frankly.

### Future Enhancements

1. Get a nice CI/CD pipeline going. I'm keen to get more exposure to Github Actions, so keen to implement these as part of deployments.

2. Creating a nice-ish frontend. I don't particularly enjoy styling or visual design. But that might be because I am bad at both. It's a useful skill to have.
