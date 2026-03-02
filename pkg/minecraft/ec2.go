package minecraft

import (
	"context"
	"maps"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Provider struct {
	client *ec2.Client
}

func NewEC2Provider(client *ec2.Client) *EC2Provider {
	return &EC2Provider{client}
}

var _ HostProvider = (*EC2Provider)(nil)

func (provider *EC2Provider) Hosts(ctx context.Context, ids ...string) ([]*HostInfo, error) {
	instanceOutput, err := provider.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: ids,
	})

	if err != nil {
		return nil, err
	}

	hostsByType := map[types.InstanceType][]*HostInfo{}
	for _, r := range instanceOutput.Reservations {
		for _, i := range r.Instances {
			info := HostInfo{
				ID:           *i.InstanceId,
				Architecture: string(i.Architecture),
				DNSName:      *i.PublicDnsName,
				IPAddress:    *i.PublicIpAddress,
				Platform:     *i.PlatformDetails,
				Processor: ProcessorInfo{
					CoreCount:      i.CpuOptions.CoreCount,
					ThreadsPerCore: i.CpuOptions.ThreadsPerCore,
				},
			}

			hostsByType[i.InstanceType] = append(hostsByType[i.InstanceType], &info)
		}
	}

	typeOutput, err := provider.client.DescribeInstanceTypes(ctx, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: slices.Collect(maps.Keys(hostsByType)),
	})

	if err != nil {
		return nil, err
	}

	hosts := slices.Collect(func(yield func(*HostInfo) bool) {
		for _, typeInfo := range typeOutput.InstanceTypes {
			image := ImageInfo{
				Type:    string(typeInfo.InstanceType),
				Memory:  typeInfo.MemoryInfo.SizeInMiB,
				Network: *typeInfo.NetworkInfo.NetworkPerformance,
			}

			for _, host := range hostsByType[typeInfo.InstanceType] {
				host.Image = &image

				if !yield(host) {
					return
				}
			}
		}
	})

	return hosts, nil
}

func (provider *EC2Provider) StartHosts(ctx context.Context, ids ...string) error {
	_, err := provider.client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: ids,
	})

	return err
}

func (provider *EC2Provider) StopHosts(ctx context.Context, ids ...string) error {
	_, err := provider.client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: ids,
	})

	return err
}

func (provider *EC2Provider) RestartHosts(ctx context.Context, ids ...string) error {
	_, err := provider.client.RebootInstances(ctx, &ec2.RebootInstancesInput{
		InstanceIds: ids,
	})

	return err
}
