# Coding Standarts & Guidelines

While developing this project, the following principles and principles are followed.

## General

- SOLID principles should be applied whenever possible.
- [Clean Code Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) principles should be applied whenever possible.
- While applying the DRY principle, "less copy paste better than dependency" should not be forgotten.

## Naming Convention

- REST API resource names must be plural. example: /users, /users/:ID
- Struct names that we bind HTTP Request payloads to must end with \*Request. example: LoginRequest
- The name of the struct containing the parameters passed to the Service Methods must end as \*Params. example: LoginParams

## Version Control (GIT)

In this project, github-flow is used instead of git-flow. github-flow differs from git-flow;

There is no branch called development in this flow. There is only one living branch, "master", every development is removed from this branch and merged into this branch.

![x](./githubflow.png)

https://docs.github.com/en/get-started/quickstart/github-flow

- A branch should be opened for each job.
- FEATURE-1150-add-forgot-password-api with branch name issue number and job summary
- Each branch should contain a single commit summarizing the issue number and the work. "FEATURE-1150 add forgot password api"
- commit messages must contain imperative mode. Instead of added forgot password api -> add forgot password api

## Repository

- When using GORM, each model must have a method called `TableName()` that explicitly freezes the table name.
