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

After looking through the various providers, I've decided to go with AWS, because:

1. They are probably the most widely used cloud provider, so having skills there will help

2. I don't get to use AWS in my day job (we're a GCP shop)

3. It's going to be _relatively_ cheap. DynamoDB is basically free at my intended level of usage, and Fargate isn't too expensive either. ECR is like 10c per month.

4. (Sort of) Fits the use case. 

In terms of cloud products:

1. I'm going to be using S3 to host my resume. Nothing too controversial here.

2. Use Fargate pulling images from ECR. I'm going to Dockerise the application. I'll get some exposure with ECS as well. While I want to get exposure to AWS products, I think manually configuring an EC2 instance is probably something for the future.

3. Ideally I'd like to use a relational database (like RDS) but that is _expensive_. So I'll use a non-relational database in the interest of cost.

### Future Enhancements

1. Get a nice CI/CD pipeline going. I'm keen to get more exposure to Github Actions, so keen to implement these as part of deployments.

2. Use Infrastructure as Code. I am unashamedly going to click-ops my way through everything initially. But that is fraught with danger and unpleasantness (having to navigate AWS' UI). I'll get some nice experience with Terraform.

3. Creating a nice-ish frontend. I don't particularly enjoy styling or visual design. But that might be because I am bad at both. It's a useful skill to have.
