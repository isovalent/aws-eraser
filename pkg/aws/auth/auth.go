package auth

import (
	"context"
	"fmt"

	"aws-eraser/pkg/log"
	"aws-eraser/pkg/resources"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
)

type clients struct {
	ASG *autoscaling.Client
	EC2 *ec2.Client
	ELB *elasticloadbalancing.Client
	EKS *eks.Client
	CF  *cloudformation.Client
}

var clientsMap = map[string]*clients{}

func InitClientsMap(ctx context.Context, accountResources resources.AccountResources) error {
	for account, res := range accountResources {
		for _, r := range res {
			key := clientKey(account, r.Region)
			if _, ok := clientsMap[key]; !ok {
				log.FromContext(ctx).Infof("initializing aws clients for: %s", key)
				cfg, err := config.LoadDefaultConfig(ctx, clientOpts(account, r.Region)...)
				if err != nil {
					return err
				}
				clientsMap[key] = &clients{
					ASG: autoscaling.NewFromConfig(cfg),
					CF:  cloudformation.NewFromConfig(cfg),
					EC2: ec2.NewFromConfig(cfg),
					ELB: elasticloadbalancing.NewFromConfig(cfg),
					EKS: eks.NewFromConfig(cfg),
				}
			}
		}
	}
	return nil
}

func clientKey(account, region string) string {
	return fmt.Sprintf("%s:%s", account, region)
}

func clientOpts(account, region string) []func(*config.LoadOptions) error {
	var opts []func(*config.LoadOptions) error
	if account != "default" {
		opts = append(opts, config.WithSharedConfigProfile(account))
	}
	if region != "default" {
		opts = append(opts, config.WithSharedConfigProfile(region))
	}
	return opts
}
