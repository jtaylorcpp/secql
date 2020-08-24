package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jtaylorcpp/secql/graph/aws"
	"github.com/jtaylorcpp/secql/graph/generated"
	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func (r *queryResolver) Ec2Instances(ctx context.Context) ([]*model.EC2Instance, error) {
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

						sshError := aws.NewEC2SSHSession(regionalSess, *instanceModel)
						if sshError != nil {
							logrus.Errorf("ssh error: %s\n", sshError.Error())
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func init() {
	logrus.SetLevel(logrus.DebugLevel)
}
