package common

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// v2 of the aws go sdk is used

const (
	ACTION_INSTANCE_STATUS_UPDATE = iota
	ACTION_INSTANCES_STATUS_UPDATE
	ACTION_INSTANCE_HALP_ME

	// Defining the services themselves as numeric constants
	// Used onwards to tweak service names and configs. This will probably be replaced
	SERVICE_EC2
    SERVICE_IAM

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

// convenient maps *shrugs*
// map to unify service names
var ServiceNames = map[int]string{
	SERVICE_EC2: "ec2",             // plus EBS as well
    SERVICE_IAM: "iam",
}

// map of subitems (tree children) names appearing at front page. this should be modeled as a tree object with children as tree nodes. this works for now
var ServiceChildrenNames = map[int][]string{
    SERVICE_EC2: []string{"Instances", "Volumes"},
    SERVICE_IAM: []string{"TODO"},
}

var AvailableServices = map[int]bool{
    SERVICE_EC2: true,
    SERVICE_IAM: true,
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
    FILTER_EGRESS_IP_PERMISSION_CIDR: []string{"egress.ip-permission.cidr"},
    FILTER_EGRESS_IP_PERMISSION_FROM_PORT: []string{"egress.ip-permission.from-port"},
    FILTER_EGRESS_IP_PERMISSION_GROUP_ID: []string{"egress.ip-permission.group-id"},
    FILTER_EGRESS_IP_PERMISSION_GROUP_NAME: []string{"egress.ip-permission.group-name"},
    FILTER_EGRESS_IP_PERMISSION_IPV6_CIDR: []string{"egress.ip-permission.ipv6-cidr"},
    FILTER_EGRESS_IP_PERMISSION_PREFIX_LIST_ID: []string{"egress.ip-permission.prefix-list-id"},
    FILTER_EGRESS_IP_PERMISSION_PROTOCOL: []string{"egress.ip-permission.protocol","tcp","icmp","udp"},
    FILTER_EGRESS_IP_PERMISSION_TO_PORT: []string{"egress.ip-permission.to-port"},
    FILTER_EGRESS_IP_PERMISSION_USER_ID: []string{"egress.ip-permission.user-id"},
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
    FILTER_IP_PERMISSION_CIDR: []string{"ip-permission.cidr"},
    FILTER_IP_PERMISSION_FROM_PORT: []string{"ip-permission.from-port"},
    FILTER_IP_PERMISSION_GROUP_ID: []string{"ip-permission.group-id"},
    FILTER_IP_PERMISSION_GROUP_NAME: []string{"ip-permission.group-name"},
    FILTER_IP_PERMISSION_IPV6_CIDR: []string{"ip-permission.ipv6-cidr"},
    FILTER_IP_PERMISSION_PREFIX_LIST_ID: []string{"ip-permission.prefix-list-id"},
    FILTER_IP_PERMISSION_PROTOCOL: []string{"ip-permission.protocol","tcp","icmp","udp"},
    FILTER_IP_PERMISSION_TO_PORT: []string{"ip-permission.to-port"},
    FILTER_IP_PERMISSION_USER_ID: []string{"ip-permission.user-id"},
	FILTER_IS_PUBLIC:                                           []string{"is-public"},
	FILTER_KERNEL_ID:                                           []string{"kernel-id"},
	FILTER_KEY_NAME:                                            []string{"key-name"},
	FILTER_LAUNCH_INDEX:                                        []string{"launch-index"},
	FILTER_LAUNCH_TIME:                                         []string{"launch-time"},
	FILTER_METADATA_OPTIONS_HTTP_TOKENS:                        []string{"metadata-options.http-tokens", "optional", "required"},
	FILTER_METADATA_OPTIONS_HTTP_PUT_RESPONSE_HOP_LIMIT:        []string{"metadata-options.http-put-response-hop-limit"},
	FILTER_METADATA_OPTIONS_HTTP_ENDPOINT:                      []string{"metadata-options.http-endpoint", "enabled", "disabled"},
	FILTER_MONITORING_STATE:                                    []string{"monitoring-state", "disabled", "enabled"},
    FILTER_NAME:                                                []string{"name","ubuntu/images/hvm-ssd/*"},
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

// these are the manually defined data structures that any first party
// should expect when receiving/sending an action. these structures are the "data" field in the action
// notice the similarity in using a name similar to the action, but in camel case
type InstanceStatusUpdate ec2.InstanceStatus

// TODO: is type(InstanceStatusesUpdate) different from type([]ec2.InstanceStatus)
type InstanceStatusesUpdate []ec2.InstanceStatus
