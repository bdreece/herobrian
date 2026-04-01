package ec2

import (
	"context"
	"maps"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bdreece/herobrian/pkg/minecraft/host"
	"github.com/bdreece/herobrian/pkg/minecraft/instance"
	"github.com/spf13/viper"
)

type Provider struct {
	client *ec2.Client
	config map[string]host.Config
	instance.ProviderFactory
}

func NewProvider(client *ec2.Client, factory instance.ProviderFactory) (*Provider, error) {
	config := map[string]host.Config{}
	if err := viper.UnmarshalKey("minecraft:hosts", &config); err != nil {
		return nil, err
	}

	return &Provider{client, config, factory}, nil
}

var _ host.Provider = (*Provider)(nil)

// DescribeHosts implements [host.Describer].
func (provider *Provider) DescribeHosts(ctx context.Context, ids ...string) (map[string]*host.Info, error) {
	instanceOutput, err := provider.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: provider.resolveHosts(ids),
	})

	if err != nil {
		return nil, err
	}

	hostsByType := map[types.InstanceType][]*host.Info{}
	for _, r := range instanceOutput.Reservations {
		for _, i := range r.Instances {
			info := host.Info{
				ID:           *i.InstanceId,
				Architecture: string(i.Architecture),
				Processor: host.ProcessorInfo{
					CoreCount:      i.CpuOptions.CoreCount,
					ThreadsPerCore: i.CpuOptions.ThreadsPerCore,
				},
			}

			if i.PublicDnsName != nil {
				info.DNSName = *i.PublicDnsName
			}
			if i.PublicIpAddress != nil {
				info.IPAddress = *i.PublicIpAddress
			}
			if i.PlatformDetails != nil {
				info.Platform = *i.PlatformDetails
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

	hosts := maps.Collect(func(yield func(string, *host.Info) bool) {
		for _, typeInfo := range typeOutput.InstanceTypes {
			image := host.ImageInfo{
				Type:    string(typeInfo.InstanceType),
				Memory:  typeInfo.MemoryInfo.SizeInMiB,
				Network: *typeInfo.NetworkInfo.NetworkPerformance,
			}

			for _, host := range hostsByType[typeInfo.InstanceType] {
				host.Image = &image

				if !yield(host.ID, host) {
					return
				}
			}
		}
	})

	return hosts, nil
}

// StartHosts implements [host.Starter].
func (provider *Provider) StartHosts(ctx context.Context, ids ...string) error {
	_, err := provider.client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: ids,
	})

	return err
}

// StopHosts implements [host.Stopper].
func (provider *Provider) StopHosts(ctx context.Context, ids ...string) error {
	_, err := provider.client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: ids,
	})

	return err
}

// RestartHosts implements [host.Restarter].
func (provider *Provider) RestartHosts(ctx context.Context, ids ...string) error {
	_, err := provider.client.RebootInstances(ctx, &ec2.RebootInstancesInput{
		InstanceIds: ids,
	})

	return err
}

// CheckHosts implements [host.Checker].
func (provider *Provider) CheckHosts(ctx context.Context, ids ...string) (map[string]host.Status, error) {
	output, err := provider.client.DescribeInstanceStatus(ctx, &ec2.DescribeInstanceStatusInput{
		InstanceIds: ids,
	})

	if err != nil {
		return nil, err
	}

	statuses := maps.Collect(func(yield func(string, host.Status) bool) {
		for _, status := range output.InstanceStatuses {
			id := *status.InstanceId

			state := host.Status(*status.InstanceState.Code & 255)

			if !yield(id, state) {
				return
			}
		}
	})

	return statuses, nil
}

func (provider *Provider) resolveHosts(hosts []string) []string {
	return slices.Collect(func(yield func(string) bool) {
		for _, host := range hosts {
			if !yield(provider.config[host].ID) {
				return
			}
		}
	})
}
