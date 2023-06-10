package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"aws-eraser/pkg/eraser"
	"aws-eraser/pkg/log"
	"github.com/spf13/cobra"
)

const (
	defaultDuration = time.Minute * 20
	minDuration     = time.Minute * 5
)

var (
	autoApprove          *bool
	duration             *time.Duration
	resources            *string
	fileName             *string
	fileFormat           *string
	supportedFileFormats = map[string]struct{}{
		"json": {},
		"yaml": {},
	}
	cmd = &cobra.Command{
		Use:   "aws-eraser",
		Short: "The tool that helps clean AWS resources",
		Args:  cobra.ExactArgs(0),
		Run:   run,
	}
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(_ *cobra.Command, _ []string) {
	ctx := context.Background()
	*fileFormat = strings.ToLower(*fileFormat)
	if err := checkParams(); err != nil {
		log.FromContext(ctx).Error(err)
		return
	}
	if err := eraser.Erase(ctx, *autoApprove, *duration, *resources, *fileName, *fileFormat); err != nil {
		log.FromContext(ctx).Error(err)
	}
}

func checkParams() error {
	if *duration < minDuration {
		return fmt.Errorf("duration must be greater then: %s", minDuration)
	}
	if *resources == "" && *fileName == "" {
		return errors.New("resources or file must be provided")
	}
	if *fileName != "" && *fileFormat == "" {
		return errors.New("file format must be provided")
	}
	if _, ok := supportedFileFormats[*fileFormat]; !ok {
		return fmt.Errorf("unsupported file format: %s", *fileFormat)
	}
	return nil
}

func init() {
	autoApprove = cmd.Flags().BoolP("auto-approve", "a", false, "Auto approve erase")
	duration = cmd.Flags().DurationP("duration", "d", defaultDuration, "Erase maximum duration in minutes")
	resources = cmd.Flags().StringP("resources", "r", "", "Resource list [e.g.: vpc:account:region:vpc-id,eks:account:region:eks-name]")
	fileName = cmd.Flags().StringP("file", "i", "", "Resource file")
	_ = cmd.MarkFlagFilename("file")
	fileFormat = cmd.Flags().StringP("file-format", "f", "yaml", "Resource file format [ yaml | json ]")
}
