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
							PublicIP:  *instance.PublicIpAddress,
							PrivateIP: *instance.PrivateIpAddress,
						}

						if instanceModel.PublicIP != "" {
							instanceModel.Public = true
						}

						for _, tag := range instance.Tags {
							if strings.ToLower(*tag.Key) == "name" {
								instanceModel.Name = *tag.Value
							}
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
