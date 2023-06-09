package eraser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"aws-eraser/pkg/aws/auth"
	"aws-eraser/pkg/log"
	"aws-eraser/pkg/resources"
	"github.com/hashicorp/go-multierror"
)

func Erase(ctx context.Context, autoApprove bool, duration time.Duration, resourceStr, fileName, fileFormat string) error {
	logger := log.FromContext(ctx)
	logger.Info("parsing resources")
	accountResources, err := resources.Read(resourceStr, fileName, fileFormat)
	if err != nil {
		return err
	}
	if len(accountResources) == 0 {
		logger.Info("no resources found, exiting...")
		return nil
	}

	logger.Infof("the below resources will be deleted:\n\n%s", accountResources)
	if !isApproved(ctx, autoApprove) {
		return nil
	}

	if err := auth.InitClientsMap(ctx, accountResources); err != nil {
		return err
	}

	wg := multierror.Group{}
	for account, res := range accountResources {
		for _, r := range res {
			acc := account
			region := r.Region
			rType := r.Type
			rID := r.ID
			wg.Go(func() error {
				logger.Infof("deleting %s:%s %s %s", acc, region, rType, rID)
				return fmt.Errorf("failed: %s", rID)
			})
		}
	}
	return wg.Wait().ErrorOrNil()
}

func isApproved(ctx context.Context, autoApprove bool) bool {
	if autoApprove {
		return true
	}
	log.FromContext(ctx).Info("do you want to continue? [y/n]")
	var confirm string
	if _, err := fmt.Scanln(&confirm); err != nil {
		log.FromContext(ctx).Errorf("failed to read user input: %s", err.Error())
		return false
	}
	return strings.ToLower(confirm) == "y"
}
