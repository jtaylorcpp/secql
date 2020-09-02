package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	awsTypes "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jtaylorcpp/secql/graph/aws"
	"github.com/jtaylorcpp/secql/graph/generated"
	"github.com/jtaylorcpp/secql/graph/model"
	osquery "github.com/jtaylorcpp/secql/graph/osquery/interactive"
	"github.com/sirupsen/logrus"
)

func (r *queryResolver) Ec2Instances(ctx context.Context) ([]*model.EC2Instance, error) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugln("getting all ec2 instances")
	regions, err := aws.GetAllRegions(r.Session)
	if err != nil {
		logrus.Errorf("error getting regions: %s", err.Error())
		return []*model.EC2Instance{}, err
	}

	logrus.Debugf("available regions: %#v", regions)

	// for each region, get ec2 instances
	instanceModels := []*model.EC2Instance{}

	for _, region := range regions {
		logrus.Debugf("running in region: %s", region)
		regionalSess := aws.GetRegionalSession(r.Session, region)
		svc := ec2.New(regionalSess)
		input := &ec2.DescribeInstancesInput{
			//InstanceIds: []*string{},
		}
		err := svc.DescribeInstancesPages(input,
			func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
				for _, reservation := range page.Reservations {
					for _, instance := range reservation.Instances {
						instanceModel := &model.EC2Instance{
							ID: *instance.InstanceId,
							//Name      string  `json:"name"`
							PublicIP:         *instance.PublicIpAddress,
							PrivateIP:        *instance.PrivateIpAddress,
							AvailabilityZone: *instance.Placement.AvailabilityZone,
						}

						if instanceModel.PublicIP != "" {
							instanceModel.Public = true
						}

						for _, tag := range instance.Tags {
							if strings.ToLower(*tag.Key) == "name" {
								instanceModel.Name = *tag.Value
							}
						}

						descirbeImageInput := &ec2.DescribeImagesInput{
							Filters: []*ec2.Filter{
								&ec2.Filter{
									Name: awsTypes.String("image-id"),
									Values: []*string{
										instance.ImageId,
									},
								},
							},
						}

						output, imageError := svc.DescribeImages(descirbeImageInput)
						if imageError != nil {
							logrus.Errorf("image describe error: %s\n", imageError.Error())
						}

						if len(output.Images) == 0 {
							logrus.Error("no images available for image")
						} else {
							instanceModel.Ami = &model.Ami{
								ID: *output.Images[0].ImageId,
							}
						}

						logrus.Debugf("got instance: %#v", instanceModel)
						sshClient, sshError := aws.NewEC2SSHSession(regionalSess, *instanceModel)
						if sshError != nil {
							logrus.Errorf("ssh error: %s\n", sshError.Error())
						}

						logrus.Debugf("got client: %#v", *sshClient)

						osqOSInfo, err := osquery.GetOS(sshClient)
						if err != nil {
							logrus.Errorf("got error from osquery OS discovery: %s", err.Error())
						}
						logrus.Debugf("osquery OS info: %#v", osqOSInfo)
						instanceModel.OsInfo = &model.OSInfo{
							ID:             osqOSInfo.Name,
							Version:        osqOSInfo.Version,
							BuildVersion:   fmt.Sprintf("%v.%v.%v", osqOSInfo.Major, osqOSInfo.Minor, osqOSInfo.Patch),
							Arch:           osqOSInfo.Arch,
							PlatformDistro: osqOSInfo.Platform,
							PlatformBase:   osqOSInfo.PlatformLike,
						}
						instanceModels = append(instanceModels, instanceModel)
					}
				}
				return !lastPage
			})
		if err != nil {
			logrus.Errorf("error in region %s: %s", region, err.Error())
		}
	}
	return instanceModels, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
