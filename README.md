# Face it instructions

## how to start development

```sh
# cp example env file and fill it
cp .env.example .env
# run app
make run-dev
```

that's all!

### An explanation of the choices taken and assumptions made during development - coding standarts

[Click me for coding standarts detail](https://github.com/safciplak/faceit/blob/master/coding-standarts.md)


## Possible extensions or improvements to the service (focusing on scalability and deployment to production

### Deployment

Single page app deployment On AWS, we first open a public S3 and enable the static website feature here.

Then we get the build of our app and assign the dist folder of this place to S3. Then we go to cloudfront and create a new distribution there and give this public S3 that I opened as origin. Here we turn on the cache optimized option in the behavior section and update the policy on the S3 to access this cloudfront distribution. Then we enter the route53 service and create a subdomain from our domain. If this page is not a static website and is a server-side or backend project, we get a docker build, we assign this build to a repo that we open to ecr. Then we create a cluster and service in ECS and add this ECR repo to the container of this service. (In the meantime, we set the load balancer and target group)

We handle the pipeline processes of these projects with GitHub actions or bitbucket pipeline. Here we make three different environments as development-stage-production. If necessary, we add unity test to the pipeline steps. For security, we can keep environment parameters in AWS parameter store or secret manager service, or we can manage them from bitbucket environment variables. In the meantime, we wake up each environment on a different VPC on AWS.

[Scalability details](https://www.pgs-soft.com/blog/scaling-containerised-applications-on-aws/)
 
