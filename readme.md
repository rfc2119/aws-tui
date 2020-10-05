# AWS Terminal Interface

[TOC]

Unofficial terminal interface for AWS. This is still a work in progress

## Download

Go to the [releases](https://github.com/rfc2119/aws-tui/releases) page and grab the latest version. Additionally, since this interface is based on the AWS Go SDK, you can compile it from source by cloning this repository and running `go build main.go` and running the output binary file

## Quick Start

Configure your credentials file the same way you use for the AWS CLI or the SDK (command-line options are not supported). To re-iterate the [documentation](https://docs.aws.amazon.com/sdk-for-go/api/):

>  When using the SDK you'll generally need your AWS credentials to authenticate with AWS services. The SDK supports multiple methods of supporting these credentials. By default the SDK will source credentials automatically from its default credential chain. [..] The common items in the credential chain are the following:
>
>  * Environment Credentials - Set of environment variables (see [supported list](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html#envvars-list))
>  * Shared Credentials [file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) (~/.aws/credentials)
>  * EC2 Instance Role Credentials

Launch the binary and there you go!

## Navigation

I tried to follow a vim-like key configuration for most of the cases. There's a help page for every window, which you can display by hitting `?`. These are common keys found across all windows: 



| Key  | Function |
| :--: | :------: |
| TAB  | Move to neighboring windows         |
| ?     |View help messages (if available) |
|q|Move back one page (will exit this help message) |
|Space|Select Option in a radio box/tree view (except in a confirmation box) |
|hjkl|Movement keys|

## Known Issues

TODO (see [issues](https://github.com/rfc2119/aws-tui/issues) for now)

## Architecture

A quick overview of the architecture is drafted at [architecture.md](https://github.com/rfc2119/aws-tui/blob/master/architecture.md)

## Contributing

Obviously this would be a huge effort for anyone to do all the work alone. If you have any issues, feature requests or would like to maintain/contribute to the software, please do not hesitate to submit an [issue](https://github.com/rfc2119/aws-tui/issues).

TODO: Contributing guide

## Support

\#aws-tui on Freenode

## Acknowledgement

I am forever indebted to all open-source community and its projects. The interface uses the work of:

1.  [tcell](https://github.com/gdamore/tcell/): grid-based terminal view for Go
2. [tview](https://github.com/rivo/tview/): the wonderful and amazing framework for creating rich interactive terminal-based programs
3. The [AWS SDK V2](https://github.com/aws/aws-sdk-go-v2/) for Go
4. The [Go](https://github.com/golang/go) team

# Current Working Services

| Service Name | Implemented | Description |
| :----------: | :---------: | :---------: |
|Amazon EC2 |  ✓  | Virtual Servers in the Cloud|
|AWS Lambda |  | Run Code Without Thinking About Servers|
|Amazon Virtual Private Cloud (VPC) |  | Isolated Cloud Resources|
|AWS Elastic Beanstalk |  | AWS Application Container|
|Amazon EC2 Auto Scaling |  | Add or remove compute capacity to meet changes in demand|
|AWS Batch |  | Fully managed batch processing at any scale|
|AWS Serverless Application Repository |  | Discover, deploy, publish and share serverless applications|
|Amazon Elastic Container Registry |  | Store and Retrieve Docker Images|
|Amazon Elastic Container Service |  | Run containerized applications in production|
|AWS Fargate |  | Run containers without managing servers or clusters|
|Amazon Elastic Kubernetes Service (EKS) |  | Fully managed Kubernetes service|
|Amazon S3 |  | Scalable Storage in the Cloud|
|Amazon Elastic Block Store (EBS) |  | Scalable Storage in the Cloud|
|Amazon Glacier |  | Low-Cost Archive Storage in the Cloud|
|AWS Snowball |  | Move petabyte-scale data sets|
|AWS Storage Gateway |  | Integrates on-premises IT environments with Cloud storage|
|Amazon Elastic File System |  | Full managed file system for EC2|
|AWS Backup |  | Centralized backup across AWS services|
|Amazon Relational Database Service (RDS) |  | Managed Relational Database Service|
|Amazon DynamoDB |  | Dynamic Databases in the Cloud|
|Amazon Aurora |  | MySQL and PostgreSQL Compatible Relational Database Built for the Cloud|
|Amazon ElastiCache |  | In-Memory Caching Service|
|Amazon Neptune |  | Fast, reliable graph database built for the cloud|
|Amazon Keyspaces (for Apache Cassandra) |  | Managed Cassandra-compatible database|
|AWS Database Migration Service |  | Migrate your databases to AWS with minimal downtime|
|AWS Server Migration Service |  | Easy migration of on-premises workloads to AWS|
|Database Migration Service |  | Migrate your databases to AWS with minimal downtime|
|Amazon CloudFront |  | Global Content Delivery Network|
|AWS Direct Connect |  | Dedicated Network Connection to AWS|
|Amazon Route 53 |  | A reliable and cost-effective way to route end users to Internet applications|
|AWS Transit Gateway |  | Easily scale VPC and account connections|
|AWS PrivateLink |  | Access services hosted on AWS easily and securely by keeping your network traffic within the AWS network|
|Elastic Load Balancing |  | Distribute incoming traffic across multiple targets|
|Amazon API Gateway |  | Create, Publish, Maintain, Monitor, and Secure APIs at Any Scale|
|AWS AppSync |  | Build data-driven apps with real-time and offline capabilities|
|AWS CodeDeploy |  | Automate Code Deployments|
|AWS X-Ray |  | Analyze and debug your applications|
|AWS CodeBuild |  | Build and test code with continuous scaling.|
|AWS CodeCommit |  | Securely host highly scalable private Git repositories. Collaborate on code.|
|Amazon WorkSpaces |  | Virtual desktops in the cloud|
|Amazon CloudWatch |  | Resource and Application Monitoring|
|AWS Organizations |  | Central governance and management across AWS accounts|
|Amazon EC2 Systems Manager |  | Configure and manage Amazon EC2 and on-premises system|
|AWS CloudFormation |  | Templates for AWS Resource Creation|
|AWS CloudTrail |  | Track user activity and API usage|
|AWS Config |  | AWS resource inventory and configuration history|
|AWS Management Console |  | Web-Based User Interface|
|AWS License Manager |  | Set rules to manage, discover, and report software license usage|
|AWS Personal Health Dashboard |  | Personalized view of AWS service health|
|AWS Budgets |  | Set custom budgets that alert you when you exceed your budgeted thresholds|
|AWS Cost Explorer |  | Visualize, understand, and manage your AWS costs and usage over time|
|AWS Cost & Usage Report |  | Dive deeper into your costs and usage|
|Reserved Instance (RI) Reporting |  | Manage and monitor your instance reservations|
|AWS Identity and Access Management (IAM) |  ✓  | Configurable AWS Access Controls|
|Amazon Cognito |  | User Sign-up and Sign-in|
|AWS Directory Service |  | Host and Manage Active Dirctory|
|AWS Key Management Service |  | Easily create and control the keys used to encrypt your data|
|AWS Secrets Manager |  | Rotate, manage, and retrieve secrets|
|AWS Certificate Manager (ACM) |  | Provision, manage, and deploy SSL/TLS certificates|
|Amazon Redshift |  | Fast, Simple, Cost-effective Data Warehousing|
|Amazon Elasticsearch Service |  | Fully managed, reliable, and scalable Elasticsearch service|
|Amazon Elastic MapReduce |  | Hosted Hadoop Framework|
|Amazon Kinesis Data Streams |  | Amazon Kinesis Data Streams|
|Amazon Kinesis Data Firehose |  | Prepare and load real-time data streams into data stores and analytics tools|
|AWS Glue |  | Prepare and load data|
|Amazon Athena |  | Query data in S3 using SQL|
|Amazon Managed Streaming for Apache Kafka (MSK) |  | Fully managed, highly available, and secure Apache Kafka service|
|Amazon Simple Workflow |  | Workflow service for coordinating applications|
|AWS Step Functions |  | Build distributed applications using visual workflows|
|Amazon EventBridge |  | Serverless event bus that connects application data from your own apps, SaaS, and AWS services|
|Amazon Simple Queue Service (SQS) |  | Message Queue Service|
|Amazon Simple Notification Service (SNS) |  | Push Notification Service|
|FreeRTOS |  | IoT operating system for microcontrollers|
|AWS IoT Greengrass |  | Bring local compute, messaging, data management, sync, and ML inference capabilities to edge devices|
|AWS IoT Device Defender |  | Security management for IoT devices|
|AWS IoT Core |  | Easily and securely connect devices to the cloud. Reliably scale to billions of devices and trillions of messages.|
|AWS IoT Device Management |  | Onboard, organize, monitor, and remotely manage connected devices at scale|
|AWS Premium Support |  | One-on-one, Fast-response Support Channel|
|AWS Deep Learning AMIs |  | Quickly build deep learning applications|
|Amazon Polly |  | Turn text into lifelike speech using deep learning|
|Amazon Transcribe |  | Automatically convert speech to text|
|Amazon SageMaker |  | Machine learning for every developer and data scientist|
|Amazon GameLift |  | Simple, fast, cost-effective dedicated game server hosting.|
|AWS Elemental MediaConvert |  | Process video files and clips to prepare on-demand content for distribution or archiving|


```

```
