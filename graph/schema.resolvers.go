package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jtaylorcpp/secql/aws"
	"github.com/jtaylorcpp/secql/graph/generated"
	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/jtaylorcpp/secql/osquery"
	"github.com/sirupsen/logrus"
)

func (r *eC2InstanceResolver) OsInfo(ctx context.Context, obj *model.EC2Instance) (*model.OSInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *eC2InstanceResolver) OsPackages(ctx context.Context, obj *model.EC2Instance) ([]*model.OSPackage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *eC2InstanceResolver) ListeningApplications(ctx context.Context, obj *model.EC2Instance) ([]*model.ListeningApplication, error) {
	panic(fmt.Errorf("not implemented"))
}

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
		instances, err := aws.GetAllEC2Instances(regionalSess, region)
		if err != nil {
			logrus.Errorf("error in region %s: %s", region, err.Error())
		}
		for _, instance := range instances {
			if !r.Resolver.Cache.Exists(instance.ID) {
				var osqueryHost string
				if instance.Public {
					osqueryHost = instance.PublicIP
				} else {
					osqueryHost = instance.PrivateIP
				}

				osqueryConfig := &osquery.ClientOpts{
					Host: osqueryHost,
					EC2SSHConfig: &osquery.OSQueryEC2SSHConfig{
						ID:        instance.ID,
						Region:    region,
						AZ:        instance.AvailabilityZone,
						IsPublic:  instance.Public,
						PublicIP:  instance.PublicIP,
						PrivateIP: instance.PrivateIP,
					},
				}

				client, err := osquery.NewClient(osqueryConfig)
				if err != nil {
					logrus.Errorf("unable to create osquery client for ec2 instance %v: %v", instance.ID, err.Error())
				} else {
					r.Resolver.Cache.Put(instance.ID, client)
				}
			}
		}
		instanceModels = append(instanceModels, instances...)
	}
	return instanceModels, nil
}

// EC2Instance returns generated.EC2InstanceResolver implementation.
func (r *Resolver) EC2Instance() generated.EC2InstanceResolver { return &eC2InstanceResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type eC2InstanceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
