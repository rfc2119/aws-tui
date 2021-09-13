package common

// v2 of the aws go sdk is used

const (
	ACTION_INSTANCES_STATUS_UPDATE = iota
	// ACTION_INSTANCE_STATUS_UPDATE
	ACTION_VOLUME_MODIFIED
	ACTION_ERROR // TODO

	// Defining the services themselves as numeric constants
	// Used onwards to tweak service names and configs. This will probably be replaced
	ServiceEc2
	ServiceLambda
	ServiceVirtualPrivateCloud
	ServiceElasticBeanstalk
	ServiceEc2AutoScaling
	ServiceBatch
	ServiceServerlessApplicationRepository
	ServiceElasticContainerRegistry
	ServiceElasticContainerService
	ServiceFargate
	ServiceElasticKubernetesService
	ServiceS3
	ServiceElasticBlockStore
	ServiceGlacier
	ServiceSnowball
	ServiceStorageGateway
	ServiceElasticFileSystem
	ServiceBackup
	ServiceRelationalDatabaseService
	ServiceDynamodb
	ServiceAurora
	ServiceElasticache
	ServiceNeptune
	ServiceKeyspaces
	ServiceDatabaseMigrationService
	ServiceServerMigrationService
	databaseMigrationService
	ServiceCloudfront
	ServiceDirectConnect
	ServiceRoute53
	ServiceTransitGateway
	ServicePrivatelink
	elasticLoadBalancing
	ServiceApiGateway
	ServiceAppsync
	ServiceCodedeploy
	ServiceXRay
	ServiceCodebuild
	ServiceCodecommit
	ServiceWorkspaces
	ServiceCloudwatch
	ServiceOrganizations
	ServiceEc2SystemsManager
	ServiceCloudformation
	ServiceCloudtrail
	ServiceConfig
	ServiceManagementConsole
	ServiceLicenseManager
	ServicePersonalHealthDashboard
	ServiceBudgets
	ServiceCostExplorer
	ServiceCostUsageReport
	reservedInstanceRiReporting
	ServiceIdentityAndAccessManagement
	ServiceCognito
	ServiceDirectoryService
	ServiceKeyManagementService
	ServiceSecretsManager
	ServiceCertificateManager
	ServiceRedshift
	ServiceElasticsearchService
	ServiceElasticMapreduce
	ServiceKinesisDataStreams
	ServiceKinesisDataFirehose
	ServiceGlue
	ServiceAthena
	ServiceManagedStreamingForApacheKafka
	ServiceSimpleWorkflow
	ServiceStepFunctions
	ServiceEventbridge
	ServiceSimpleQueueService
	ServiceSimpleNotificationService
	freertos
	ServiceIotGreengrass
	ServiceIotDeviceDefender
	ServiceIotCore
	ServiceIotDeviceManagement
	ServicePremiumSupport
	ServiceDeepLearningAmis
	ServicePolly
	ServiceTranscribe
	ServiceSagemaker
	ServiceGamelift
	ServiceElementalMediaconvert

	// filters
	FILTER_AFFINITY
	FILTER_ARCHITECTURE
	FILTER_AVAILABILITY_ZONE
	FILTER_BLOCK_DEVICE_MAPPING_ATTACH_TIME
	FILTER_BLOCK_DEVICE_MAPPING_DELETE_ON_TERMINATION
	FILTER_BLOCK_DEVICE_MAPPING_DEVICE_NAME
	FILTER_BLOCK_DEVICE_MAPPING_STATUS
	FILTER_BLOCK_DEVICE_MAPPING_VOLUME_ID
	FILTER_BLOCK_DEVICE_MAPPING_VOLUME_TYPE
	FILTER_BLOCK_DEVICE_MAPPING_ENCRYPTED
	FILTER_CLIENT_TOKEN
	FILTER_DESCRIPTION
	FILTER_DNS_NAME
	FILTER_EGRESS_IP_PERMISSION_CIDR
	FILTER_EGRESS_IP_PERMISSION_FROM_PORT
	FILTER_EGRESS_IP_PERMISSION_GROUP_ID
	FILTER_EGRESS_IP_PERMISSION_GROUP_NAME
	FILTER_EGRESS_IP_PERMISSION_IPV6_CIDR
	FILTER_EGRESS_IP_PERMISSION_PREFIX_LIST_ID
	FILTER_EGRESS_IP_PERMISSION_PROTOCOL
	FILTER_EGRESS_IP_PERMISSION_TO_PORT
	FILTER_EGRESS_IP_PERMISSION_USER_ID
	FILTER_ENA_SUPPORT
	FILTER_GROUP_ID
	FILTER_GROUP_NAME
	FILTER_HIBERNATION_OPTIONS_CONFIGURED
	FILTER_HOST_ID
	FILTER_HYPERVISOR
	FILTER_IAM_INSTANCE_PROFILE_ARN
	FILTER_IMAGE_ID
	FILTER_IMAGE_TYPE
	FILTER_INSTANCE_ID
	FILTER_INSTANCE_LIFECYCLE
	FILTER_INSTANCE_STATE_CODE
	FILTER_INSTANCE_STATE_NAME
	FILTER_INSTANCE_TYPE
	FILTER_INSTANCE_GROUP_ID
	FILTER_INSTANCE_GROUP_NAME
	FILTER_IP_ADDRESS
	FILTER_IP_PERMISSION_CIDR
	FILTER_IP_PERMISSION_FROM_PORT
	FILTER_IP_PERMISSION_GROUP_ID
	FILTER_IP_PERMISSION_GROUP_NAME
	FILTER_IP_PERMISSION_IPV6_CIDR
	FILTER_IP_PERMISSION_PREFIX_LIST_ID
	FILTER_IP_PERMISSION_PROTOCOL
	FILTER_IP_PERMISSION_TO_PORT
	FILTER_IP_PERMISSION_USER_ID
	FILTER_IS_PUBLIC
	FILTER_KERNEL_ID
	FILTER_KEY_NAME
	FILTER_LAUNCH_INDEX
	FILTER_LAUNCH_TIME
	FILTER_METADATA_OPTIONS_HTTP_TOKENS
	FILTER_METADATA_OPTIONS_HTTP_PUT_RESPONSE_HOP_LIMIT
	FILTER_METADATA_OPTIONS_HTTP_ENDPOINT
	FILTER_MONITORING_STATE
	FILTER_NAME
	FILTER_NETWORK_INTERFACE_ADDRESSES_PRIVATE_IP_ADDRESS
	FILTER_NETWORK_INTERFACE_ADDRESSES_PRIMARY
	FILTER_NETWORK_INTERFACE_ADDRESSES_ASSOCIATION_PUBLIC_IP
	FILTER_NETWORK_INTERFACE_ADDRESSES_ASSOCIATION_IP_OWNER_ID
	FILTER_NETWORK_INTERFACE_ASSOCIATION_PUBLIC_IP
	FILTER_NETWORK_INTERFACE_ASSOCIATION_IP_OWNER_ID
	FILTER_NETWORK_INTERFACE_ASSOCIATION_ALLOCATION_ID
	FILTER_NETWORK_INTERFACE_ASSOCIATION_ASSOCIATION_ID
	FILTER_NETWORK_INTERFACE_ATTACHMENT_ATTACHMENT_ID
	FILTER_NETWORK_INTERFACE_ATTACHMENT_INSTANCE_ID
	FILTER_NETWORK_INTERFACE_ATTACHMENT_INSTANCE_OWNER_ID
	FILTER_NETWORK_INTERFACE_ATTACHMENT_DEVICE_INDEX
	FILTER_NETWORK_INTERFACE_ATTACHMENT_STATUS
	FILTER_NETWORK_INTERFACE_ATTACHMENT_ATTACH_TIME
	FILTER_NETWORK_INTERFACE_ATTACHMENT_DELETE_ON_TERMINATION
	FILTER_NETWORK_INTERFACE_AVAILABILITY_ZONE
	FILTER_NETWORK_INTERFACE_DESCRIPTION
	FILTER_NETWORK_INTERFACE_GROUP_ID
	FILTER_NETWORK_INTERFACE_GROUP_NAME
	FILTER_NETWORK_INTERFACE_IPV6_ADDRESSES_IPV6_ADDRESS
	FILTER_NETWORK_INTERFACE_MAC_ADDRESS
	FILTER_NETWORK_INTERFACE_NETWORK_INTERFACE_ID
	FILTER_NETWORK_INTERFACE_OWNER_ID
	FILTER_NETWORK_INTERFACE_PRIVATE_DNS_NAME
	FILTER_NETWORK_INTERFACE_REQUESTER_ID
	FILTER_NETWORK_INTERFACE_REQUESTER_MANAGED
	FILTER_NETWORK_INTERFACE_STATUS
	FILTER_NETWORK_INTERFACE_SOURCE_DEST_CHECK
	FILTER_NETWORK_INTERFACE_SUBNET_ID
	FILTER_NETWORK_INTERFACE_VPC_ID
	FILTER_OWNER_ALIAS
	FILTER_OWNER_ID
	FILTER_PLACEMENT_GROUP_NAME
	FILTER_PLACEMENT_PARTITION_NUMBER
	FILTER_PLATFORM
	FILTER_PRIVATE_DNS_NAME
	FILTER_PRIVATE_IP_ADDRESS
	FILTER_PRODUCT_CODE
	FILTER_PRODUCT_CODE_TYPE
	FILTER_RAMDISK_ID
	FILTER_REASON
	FILTER_REQUESTER_ID
	FILTER_RESERVATION_ID
	FILTER_ROOT_DEVICE_NAME
	FILTER_ROOT_DEVICE_TYPE
	FILTER_SOURCE_DEST_CHECK
	FILTER_SPOT_INSTANCE_REQUEST_ID
	FILTER_STATE
	FILTER_STATE_REASON_CODE
	FILTER_STATE_REASON_MESSAGE
	FILTER_SUBNET_ID
	FILTER_TAG_KEY
	FILTER_TENANCY
	FILTER_VIRTUALIZATION_TYPE
	FILTER_VPC_ID
)

type awsServiceDescription struct {
	Name        string
	Description string
	Available   bool
}

var AWServicesDescriptions = map[int]awsServiceDescription{
	ServiceEc2: {
		Name: "Amazon EC2", Description: "Virtual Servers in the Cloud", Available: true,
	},

	ServiceLambda: {
		Name: "AWS Lambda", Description: "Run Code Without Thinking About Servers", Available: false,
	},

	ServiceVirtualPrivateCloud: {
		Name: "Amazon Virtual Private Cloud (VPC)", Description: "Isolated Cloud Resources", Available: false,
	},

	ServiceElasticBeanstalk: {
		Name: "AWS Elastic Beanstalk", Description: "AWS Application Container", Available: false,
	},

	ServiceEc2AutoScaling: {
		Name: "Amazon EC2 Auto Scaling", Description: "Add or remove compute capacity to meet changes in demand", Available: false,
	},

	ServiceBatch: {
		Name: "AWS Batch", Description: "Fully managed batch processing at any scale", Available: false,
	},

	ServiceServerlessApplicationRepository: {
		Name: "AWS Serverless Application Repository", Description: "Discover, deploy, publish and share serverless applications", Available: false,
	},

	ServiceElasticContainerRegistry: {
		Name: "Amazon Elastic Container Registry", Description: "Store and Retrieve Docker Images", Available: false,
	},

	ServiceElasticContainerService: {
		Name: "Amazon Elastic Container Service", Description: "Run containerized applications in production", Available: false,
	},

	ServiceFargate: {
		Name: "AWS Fargate", Description: "Run containers without managing servers or clusters", Available: false,
	},

	ServiceElasticKubernetesService: {
		Name: "Amazon Elastic Kubernetes Service (EKS)", Description: "Fully managed Kubernetes service", Available: false,
	},

	ServiceS3: {
		Name: "Amazon S3", Description: "Scalable Storage in the Cloud", Available: false,
	},

	ServiceElasticBlockStore: {
		Name: "Amazon Elastic Block Store (EBS)", Description: "Scalable Storage in the Cloud", Available: false,
	}, // Merged in the EC2 service for now

	ServiceGlacier: {
		Name: "Amazon Glacier", Description: "Low-Cost Archive Storage in the Cloud", Available: false,
	},

	ServiceSnowball: {
		Name: "AWS Snowball", Description: "Move petabyte-scale data sets", Available: false,
	},

	ServiceStorageGateway: {
		Name: "AWS Storage Gateway", Description: "Integrates on-premises IT environments with Cloud storage", Available: false,
	},

	ServiceElasticFileSystem: {
		Name: "Amazon Elastic File System", Description: "Full managed file system for EC2", Available: false,
	},

	ServiceBackup: {
		Name: "AWS Backup", Description: "Centralized backup across AWS services", Available: false,
	},

	ServiceRelationalDatabaseService: {
		Name: "Amazon Relational Database Service (RDS)", Description: "Managed Relational Database Service", Available: false,
	},

	ServiceDynamodb: {
		Name: "Amazon DynamoDB", Description: "Dynamic Databases in the Cloud", Available: false,
	},

	ServiceAurora: {
		Name: "Amazon Aurora", Description: "MySQL and PostgreSQL Compatible Relational Database Built for the Cloud", Available: false,
	},

	ServiceElasticache: {
		Name: "Amazon ElastiCache", Description: "In-Memory Caching Service", Available: false,
	},

	ServiceNeptune: {
		Name: "Amazon Neptune", Description: "Fast, reliable graph database built for the cloud", Available: false,
	},

	ServiceKeyspaces: {
		Name: "Amazon Keyspaces (for Apache Cassandra)", Description: "Managed Cassandra-compatible database", Available: false,
	},

	ServiceDatabaseMigrationService: {
		Name: "AWS Database Migration Service", Description: "Migrate your databases to AWS with minimal downtime", Available: false,
	},

	ServiceServerMigrationService: {
		Name: "AWS Server Migration Service", Description: "Easy migration of on-premises workloads to AWS", Available: false,
	},

	databaseMigrationService: {
		Name: "Database Migration Service", Description: "Migrate your databases to AWS with minimal downtime", Available: false,
	},

	ServiceCloudfront: {
		Name: "Amazon CloudFront", Description: "Global Content Delivery Network", Available: false,
	},

	ServiceDirectConnect: {
		Name: "AWS Direct Connect", Description: "Dedicated Network Connection to AWS", Available: false,
	},

	ServiceRoute53: {
		Name: "Amazon Route 53", Description: "A reliable and cost-effective way to route end users to Internet applications", Available: false,
	},

	ServiceTransitGateway: {
		Name: "AWS Transit Gateway", Description: "Easily scale VPC and account connections", Available: false,
	},

	ServicePrivatelink: {
		Name: "AWS PrivateLink", Description: "Access services hosted on AWS easily and securely by keeping your network traffic within the AWS network", Available: false,
	},

	elasticLoadBalancing: {
		Name: "Elastic Load Balancing", Description: "Distribute incoming traffic across multiple targets", Available: false,
	},

	ServiceApiGateway: {
		Name: "Amazon API Gateway", Description: "Create, Publish, Maintain, Monitor, and Secure APIs at Any Scale", Available: false,
	},

	ServiceAppsync: {
		Name: "AWS AppSync", Description: "Build data-driven apps with real-time and offline capabilities", Available: false,
	},

	ServiceCodedeploy: {
		Name: "AWS CodeDeploy", Description: "Automate Code Deployments", Available: false,
	},

	ServiceXRay: {
		Name: "AWS X-Ray", Description: "Analyze and debug your applications", Available: false,
	},

	ServiceCodebuild: {
		Name: "AWS CodeBuild", Description: "Build and test code with continuous scaling.", Available: false,
	},

	ServiceCodecommit: {
		Name: "AWS CodeCommit", Description: "Securely host highly scalable private Git repositories. Collaborate on code.", Available: false,
	},

	ServiceWorkspaces: {
		Name: "Amazon WorkSpaces", Description: "Virtual desktops in the cloud", Available: false,
	},

	ServiceCloudwatch: {
		Name: "Amazon CloudWatch", Description: "Resource and Application Monitoring", Available: false,
	},

	ServiceOrganizations: {
		Name: "AWS Organizations", Description: "Central governance and management across AWS accounts", Available: false,
	},

	ServiceEc2SystemsManager: {
		Name: "Amazon EC2 Systems Manager", Description: "Configure and manage Amazon EC2 and on-premises system", Available: false,
	},

	ServiceCloudformation: {
		Name: "AWS CloudFormation", Description: "Templates for AWS Resource Creation", Available: false,
	},

	ServiceCloudtrail: {
		Name: "AWS CloudTrail", Description: "Track user activity and API usage", Available: false,
	},

	ServiceConfig: {
		Name: "AWS Config", Description: "AWS resource inventory and configuration history", Available: false,
	},

	ServiceManagementConsole: {
		Name: "AWS Management Console", Description: "Web-Based User Interface", Available: false,
	},

	ServiceLicenseManager: {
		Name: "AWS License Manager", Description: "Set rules to manage, discover, and report software license usage", Available: false,
	},

	ServicePersonalHealthDashboard: {
		Name: "AWS Personal Health Dashboard", Description: "Personalized view of AWS service health", Available: false,
	},

	ServiceBudgets: {
		Name: "AWS Budgets", Description: "Set custom budgets that alert you when you exceed your budgeted thresholds", Available: false,
	},

	ServiceCostExplorer: {
		Name: "AWS Cost Explorer", Description: "Visualize, understand, and manage your AWS costs and usage over time", Available: false,
	},

	ServiceCostUsageReport: {
		Name: "AWS Cost & Usage Report", Description: "Dive deeper into your costs and usage", Available: false,
	},

	reservedInstanceRiReporting: {
		Name: "Reserved Instance (RI) Reporting", Description: "Manage and monitor your instance reservations", Available: false,
	},

	ServiceIdentityAndAccessManagement: {
		Name: "AWS Identity and Access Management (IAM)", Description: "Configurable AWS Access Controls", Available: true,
	},

	ServiceCognito: {
		Name: "Amazon Cognito", Description: "User Sign-up and Sign-in", Available: false,
	},

	ServiceDirectoryService: {
		Name: "AWS Directory Service", Description: "Host and Manage Active Dirctory", Available: false,
	},

	ServiceKeyManagementService: {
		Name: "AWS Key Management Service", Description: "Easily create and control the keys used to encrypt your data", Available: false,
	},

	ServiceSecretsManager: {
		Name: "AWS Secrets Manager", Description: "Rotate, manage, and retrieve secrets", Available: false,
	},

	ServiceCertificateManager: {
		Name: "AWS Certificate Manager (ACM)", Description: "Provision, manage, and deploy SSL/TLS certificates", Available: false,
	},

	ServiceRedshift: {
		Name: "Amazon Redshift", Description: "Fast, Simple, Cost-effective Data Warehousing", Available: false,
	},

	ServiceElasticsearchService: {
		Name: "Amazon Elasticsearch Service", Description: "Fully managed, reliable, and scalable Elasticsearch service", Available: false,
	},

	ServiceElasticMapreduce: {
		Name: "Amazon Elastic MapReduce", Description: "Hosted Hadoop Framework", Available: false,
	},

	ServiceKinesisDataStreams: {
		Name: "Amazon Kinesis Data Streams", Description: "Amazon Kinesis Data Streams", Available: false,
	},

	ServiceKinesisDataFirehose: {
		Name: "Amazon Kinesis Data Firehose", Description: "Prepare and load real-time data streams into data stores and analytics tools", Available: false,
	},

	ServiceGlue: {
		Name: "AWS Glue", Description: "Prepare and load data", Available: false,
	},

	ServiceAthena: {
		Name: "Amazon Athena", Description: "Query data in S3 using SQL", Available: false,
	},

	ServiceManagedStreamingForApacheKafka: {
		Name: "Amazon Managed Streaming for Apache Kafka (MSK)", Description: "Fully managed, highly available, and secure Apache Kafka service", Available: false,
	},

	ServiceSimpleWorkflow: {
		Name: "Amazon Simple Workflow", Description: "Workflow service for coordinating applications", Available: false,
	},

	ServiceStepFunctions: {
		Name: "AWS Step Functions", Description: "Build distributed applications using visual workflows", Available: false,
	},

	ServiceEventbridge: {
		Name: "Amazon EventBridge", Description: "Serverless event bus that connects application data from your own apps, SaaS, and AWS services", Available: false,
	},

	ServiceSimpleQueueService: {
		Name: "Amazon Simple Queue Service (SQS)", Description: "Message Queue Service", Available: false,
	},

	ServiceSimpleNotificationService: {
		Name: "Amazon Simple Notification Service (SNS)", Description: "Push Notification Service", Available: false,
	},

	freertos: {
		Name: "FreeRTOS", Description: "IoT operating system for microcontrollers", Available: false,
	},

	ServiceIotGreengrass: {
		Name: "AWS IoT Greengrass", Description: "Bring local compute, messaging, data management, sync, and ML inference capabilities to edge devices", Available: false,
	},

	ServiceIotDeviceDefender: {
		Name: "AWS IoT Device Defender", Description: "Security management for IoT devices", Available: false,
	},

	ServiceIotCore: {
		Name: "AWS IoT Core", Description: "Easily and securely connect devices to the cloud. Reliably scale to billions of devices and trillions of messages.", Available: false,
	},

	ServiceIotDeviceManagement: {
		Name: "AWS IoT Device Management", Description: "Onboard, organize, monitor, and remotely manage connected devices at scale", Available: false,
	},

	ServicePremiumSupport: {
		Name: "AWS Premium Support", Description: "One-on-one, Fast-response Support Channel", Available: false,
	},

	ServiceDeepLearningAmis: {
		Name: "AWS Deep Learning AMIs", Description: "Quickly build deep learning applications", Available: false,
	},

	ServicePolly: {
		Name: "Amazon Polly", Description: "Turn text into lifelike speech using deep learning", Available: false,
	},

	ServiceTranscribe: {
		Name: "Amazon Transcribe", Description: "Automatically convert speech to text", Available: false,
	},

	ServiceSagemaker: {
		Name: "Amazon SageMaker", Description: "Machine learning for every developer and data scientist", Available: false,
	},

	ServiceGamelift: {
		Name: "Amazon GameLift", Description: "Simple, fast, cost-effective dedicated game server hosting.", Available: false,
	},

	ServiceElementalMediaconvert: {
		Name: "AWS Elemental MediaConvert", Description: "Process video files and clips to prepare on-demand content for distribution or archiving", Available: false},
}

// convenient maps *shrugs*

// map of sub items (tree children) names appearing at front page. this should be modeled as a tree object with children as tree nodes. this works for now
var ServiceChildrenNames = map[int][]string{
	ServiceEc2:                         []string{"Instances", "Volumes"},
	ServiceIdentityAndAccessManagement: []string{"TODO"},
}

// map of filter names and some of the default values
var FilterNames = map[int][]string{
	FILTER_AFFINITY:                                            []string{"affinity", "default", "host"},
	FILTER_ARCHITECTURE:                                        []string{"architecture", "i386", "x86_64", "arm64"},
	FILTER_AVAILABILITY_ZONE:                                   []string{"availability-zone"},
	FILTER_BLOCK_DEVICE_MAPPING_ATTACH_TIME:                    []string{"block-device-mapping.attach-time"},
	FILTER_BLOCK_DEVICE_MAPPING_DELETE_ON_TERMINATION:          []string{"block-device-mapping.delete-on-termination"},
	FILTER_BLOCK_DEVICE_MAPPING_DEVICE_NAME:                    []string{"block-device-mapping.device-name"},
	FILTER_BLOCK_DEVICE_MAPPING_STATUS:                         []string{"block-device-mapping.status", "attaching", "attached", "detaching", "detached"},
	FILTER_BLOCK_DEVICE_MAPPING_VOLUME_ID:                      []string{"block-device-mapping.volume-id"},
	FILTER_BLOCK_DEVICE_MAPPING_VOLUME_TYPE:                    []string{"block-device-mapping.volume-type", "gp2", "io1", "io2", "st1", "sc1", "standard"},
	FILTER_BLOCK_DEVICE_MAPPING_ENCRYPTED:                      []string{"block-device-mapping.encrypted"},
	FILTER_CLIENT_TOKEN:                                        []string{"client-token"},
	FILTER_DESCRIPTION:                                         []string{"description"},
	FILTER_DNS_NAME:                                            []string{"dns-name"},
	FILTER_EGRESS_IP_PERMISSION_CIDR:                           []string{"egress.ip-permission.cidr"},
	FILTER_EGRESS_IP_PERMISSION_FROM_PORT:                      []string{"egress.ip-permission.from-port"},
	FILTER_EGRESS_IP_PERMISSION_GROUP_ID:                       []string{"egress.ip-permission.group-id"},
	FILTER_EGRESS_IP_PERMISSION_GROUP_NAME:                     []string{"egress.ip-permission.group-name"},
	FILTER_EGRESS_IP_PERMISSION_IPV6_CIDR:                      []string{"egress.ip-permission.ipv6-cidr"},
	FILTER_EGRESS_IP_PERMISSION_PREFIX_LIST_ID:                 []string{"egress.ip-permission.prefix-list-id"},
	FILTER_EGRESS_IP_PERMISSION_PROTOCOL:                       []string{"egress.ip-permission.protocol", "tcp", "icmp", "udp"},
	FILTER_EGRESS_IP_PERMISSION_TO_PORT:                        []string{"egress.ip-permission.to-port"},
	FILTER_EGRESS_IP_PERMISSION_USER_ID:                        []string{"egress.ip-permission.user-id"},
	FILTER_ENA_SUPPORT:                                         []string{"ena-support"},
	FILTER_GROUP_ID:                                            []string{"group-id"},
	FILTER_GROUP_NAME:                                          []string{"group-name"},
	FILTER_HIBERNATION_OPTIONS_CONFIGURED:                      []string{"hibernation-options.configured"},
	FILTER_HOST_ID:                                             []string{"host-id"},
	FILTER_HYPERVISOR:                                          []string{"hypervisor", "ovm", "xen"},
	FILTER_IAM_INSTANCE_PROFILE_ARN:                            []string{"iam-instance-profile.arn"},
	FILTER_IMAGE_ID:                                            []string{"image-id"},
	FILTER_IMAGE_TYPE:                                          []string{"image-type", "machine", "kernel", "ramdisk"},
	FILTER_INSTANCE_ID:                                         []string{"instance-id"},
	FILTER_INSTANCE_LIFECYCLE:                                  []string{"instance-lifecycle", "spot", "scheduled"},
	FILTER_INSTANCE_STATE_CODE:                                 []string{"instance-state-code", "0", "16", "32", "48", "64", "80"},
	FILTER_INSTANCE_STATE_NAME:                                 []string{"instance-state-name", "pending", "running", "shutting-down", "terminated", "stopping", "stopped"},
	FILTER_INSTANCE_TYPE:                                       []string{"instance-type"},
	FILTER_INSTANCE_GROUP_ID:                                   []string{"instance.group-id"},
	FILTER_INSTANCE_GROUP_NAME:                                 []string{"instance.group-name"},
	FILTER_IP_ADDRESS:                                          []string{"ip-address"},
	FILTER_IP_PERMISSION_CIDR:                                  []string{"ip-permission.cidr"},
	FILTER_IP_PERMISSION_FROM_PORT:                             []string{"ip-permission.from-port"},
	FILTER_IP_PERMISSION_GROUP_ID:                              []string{"ip-permission.group-id"},
	FILTER_IP_PERMISSION_GROUP_NAME:                            []string{"ip-permission.group-name"},
	FILTER_IP_PERMISSION_IPV6_CIDR:                             []string{"ip-permission.ipv6-cidr"},
	FILTER_IP_PERMISSION_PREFIX_LIST_ID:                        []string{"ip-permission.prefix-list-id"},
	FILTER_IP_PERMISSION_PROTOCOL:                              []string{"ip-permission.protocol", "tcp", "icmp", "udp"},
	FILTER_IP_PERMISSION_TO_PORT:                               []string{"ip-permission.to-port"},
	FILTER_IP_PERMISSION_USER_ID:                               []string{"ip-permission.user-id"},
	FILTER_IS_PUBLIC:                                           []string{"is-public"},
	FILTER_KERNEL_ID:                                           []string{"kernel-id"},
	FILTER_KEY_NAME:                                            []string{"key-name"},
	FILTER_LAUNCH_INDEX:                                        []string{"launch-index"},
	FILTER_LAUNCH_TIME:                                         []string{"launch-time"},
	FILTER_METADATA_OPTIONS_HTTP_TOKENS:                        []string{"metadata-options.http-tokens", "optional", "required"},
	FILTER_METADATA_OPTIONS_HTTP_PUT_RESPONSE_HOP_LIMIT:        []string{"metadata-options.http-put-response-hop-limit"},
	FILTER_METADATA_OPTIONS_HTTP_ENDPOINT:                      []string{"metadata-options.http-endpoint", "enabled", "disabled"},
	FILTER_MONITORING_STATE:                                    []string{"monitoring-state", "disabled", "enabled"},
	FILTER_NAME:                                                []string{"name", "ubuntu/images/hvm-ssd/*"},
	FILTER_NETWORK_INTERFACE_ADDRESSES_PRIVATE_IP_ADDRESS:      []string{"network-interface.addresses.private-ip-address"},
	FILTER_NETWORK_INTERFACE_ADDRESSES_PRIMARY:                 []string{"network-interface.addresses.primary"},
	FILTER_NETWORK_INTERFACE_ADDRESSES_ASSOCIATION_PUBLIC_IP:   []string{"network-interface.addresses.association.public-ip"},
	FILTER_NETWORK_INTERFACE_ADDRESSES_ASSOCIATION_IP_OWNER_ID: []string{"network-interface.addresses.association.ip-owner-id"},
	FILTER_NETWORK_INTERFACE_ASSOCIATION_PUBLIC_IP:             []string{"network-interface.association.public-ip"},
	FILTER_NETWORK_INTERFACE_ASSOCIATION_IP_OWNER_ID:           []string{"network-interface.association.ip-owner-id"},
	FILTER_NETWORK_INTERFACE_ASSOCIATION_ALLOCATION_ID:         []string{"network-interface.association.allocation-id"},
	FILTER_NETWORK_INTERFACE_ASSOCIATION_ASSOCIATION_ID:        []string{"network-interface.association.association-id"},
	FILTER_NETWORK_INTERFACE_ATTACHMENT_ATTACHMENT_ID:          []string{"network-interface.attachment.attachment-id"},
	FILTER_NETWORK_INTERFACE_ATTACHMENT_INSTANCE_ID:            []string{"network-interface.attachment.instance-id"},
	FILTER_NETWORK_INTERFACE_ATTACHMENT_INSTANCE_OWNER_ID:      []string{"network-interface.attachment.instance-owner-id"},
	FILTER_NETWORK_INTERFACE_ATTACHMENT_DEVICE_INDEX:           []string{"network-interface.attachment.device-index"},
	FILTER_NETWORK_INTERFACE_ATTACHMENT_STATUS:                 []string{"network-interface.attachment.status", "attaching", "attached", "detaching", "detached"},
	FILTER_NETWORK_INTERFACE_ATTACHMENT_ATTACH_TIME:            []string{"network-interface.attachment.attach-time"},
	FILTER_NETWORK_INTERFACE_ATTACHMENT_DELETE_ON_TERMINATION:  []string{"network-interface.attachment.delete-on-termination"},
	FILTER_NETWORK_INTERFACE_AVAILABILITY_ZONE:                 []string{"network-interface.availability-zone"},
	FILTER_NETWORK_INTERFACE_DESCRIPTION:                       []string{"network-interface.description"},
	FILTER_NETWORK_INTERFACE_GROUP_ID:                          []string{"network-interface.group-id"},
	FILTER_NETWORK_INTERFACE_GROUP_NAME:                        []string{"network-interface.group-name"},
	FILTER_NETWORK_INTERFACE_IPV6_ADDRESSES_IPV6_ADDRESS:       []string{"network-interface.ipv6-addresses.ipv6-address"},
	FILTER_NETWORK_INTERFACE_MAC_ADDRESS:                       []string{"network-interface.mac-address"},
	FILTER_NETWORK_INTERFACE_NETWORK_INTERFACE_ID:              []string{"network-interface.network-interface-id"},
	FILTER_NETWORK_INTERFACE_OWNER_ID:                          []string{"network-interface.owner-id"},
	FILTER_NETWORK_INTERFACE_PRIVATE_DNS_NAME:                  []string{"network-interface.private-dns-name"},
	FILTER_NETWORK_INTERFACE_REQUESTER_ID:                      []string{"network-interface.requester-id"},
	FILTER_NETWORK_INTERFACE_REQUESTER_MANAGED:                 []string{"network-interface.requester-managed"},
	FILTER_NETWORK_INTERFACE_STATUS:                            []string{"network-interface.status", "available", "in-use"},
	FILTER_NETWORK_INTERFACE_SOURCE_DEST_CHECK:                 []string{"network-interface.source-dest-check"},
	FILTER_NETWORK_INTERFACE_SUBNET_ID:                         []string{"network-interface.subnet-id"},
	FILTER_NETWORK_INTERFACE_VPC_ID:                            []string{"network-interface.vpc-id"},
	FILTER_OWNER_ALIAS:                                         []string{"owner-alias", "amazon", "aws-marketplace", "self"},
	FILTER_OWNER_ID:                                            []string{"owner-id"},
	FILTER_PLACEMENT_GROUP_NAME:                                []string{"placement-group-name"},
	FILTER_PLACEMENT_PARTITION_NUMBER:                          []string{"placement-partition-number"},
	FILTER_PLATFORM:                                            []string{"platform", "windows"},
	FILTER_PRIVATE_DNS_NAME:                                    []string{"private-dns-name"},
	FILTER_PRIVATE_IP_ADDRESS:                                  []string{"private-ip-address"},
	FILTER_PRODUCT_CODE:                                        []string{"product-code"},
	FILTER_PRODUCT_CODE_TYPE:                                   []string{"product-code.type", "devpay", "marketplace"},
	FILTER_RAMDISK_ID:                                          []string{"ramdisk-id"},
	FILTER_REASON:                                              []string{"reason"},
	FILTER_REQUESTER_ID:                                        []string{"requester-id"},
	FILTER_RESERVATION_ID:                                      []string{"reservation-id"},
	FILTER_ROOT_DEVICE_NAME:                                    []string{"root-device-name"},
	FILTER_ROOT_DEVICE_TYPE:                                    []string{"root-device-type", "ebs", "instance-store"},
	FILTER_SOURCE_DEST_CHECK:                                   []string{"source-dest-check"},
	FILTER_SPOT_INSTANCE_REQUEST_ID:                            []string{"spot-instance-request-id"},
	FILTER_STATE:                                               []string{"state", "available", "pending", "failed"},
	FILTER_STATE_REASON_CODE:                                   []string{"state-reason-code"},
	FILTER_STATE_REASON_MESSAGE:                                []string{"state-reason-message"},
	FILTER_SUBNET_ID:                                           []string{"subnet-id"},
	FILTER_TAG_KEY:                                             []string{"tag-key"},
	FILTER_TENANCY:                                             []string{"tenancy", "dedicated", "default", "host"},
	FILTER_VIRTUALIZATION_TYPE:                                 []string{"virtualization-type", "paravirtual", "hvm"},
	FILTER_VPC_ID:                                              []string{"vpc-id"},
}

// actions are defined for each service. right now, this is a way to define the behavior that the view should follow. in some sense, there ought to be some generalization of this instead of defining everything manually, but we'll see how things goes

// this is the generic action structure. the data field is an interface. this permits all other structures (defined next) to be passed onto the channel. on the receiving side of the work channel, each receiver should assert the type of the data field and act accordingly
type Action struct {
	Type int
	Data interface{}
}

// UPDATE: will deprecate these for now
// these are the manually defined data structures that any first party
// should expect when receiving/sending an action. these structures are the "data" field in the action
// notice the similarity in using a name similar to the action, but in camel case
// type InstanceStatusUpdate ec2.InstanceStatus
//
// type InstanceStatusesUpdate []ec2.InstanceStatus
