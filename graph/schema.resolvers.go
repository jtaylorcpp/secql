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
	logrus.Debugf("getting os info for instance %v", obj.ID)
	var client osquery.Client = nil
	client = r.Resolver.Cache.Get(obj.ID)
	if client == nil {
		logrus.Debugf("no client for instance %v", obj.ID)
		var osqueryHost string
		if obj.Public {
			osqueryHost = "http://" + obj.PublicIP
		} else {
			osqueryHost = "http://" + obj.PrivateIP
		}

		osqueryConfig := &osquery.ClientOpts{
			Host: osqueryHost,
			EC2SSHConfig: &osquery.OSQueryEC2SSHConfig{
				ID:        obj.ID,
				Region:    obj.Region,
				AZ:        obj.AvailabilityZone,
				IsPublic:  obj.Public,
				PublicIP:  obj.PublicIP,
				PrivateIP: obj.PrivateIP,
			},
		}
		var err error = nil
		client, err = osquery.NewClient(osqueryConfig)
		if err != nil {
			logrus.Errorf("unable to create osquery client for ec2 instance %v: %v", obj.ID, err.Error())
		} else {
			r.Resolver.Cache.Put(obj.ID, client)
		}
	}

	if client == nil {
		return nil, fmt.Errorf("no client could be found for instance %v", obj.ID)
	}

	info, infoError := client.GetOSInfo()
	if infoError != nil {
		logrus.Errorf("error getting os info for instance %v: %v", obj.ID, infoError.Error())
	}

	return &info, infoError
}

func (r *eC2InstanceResolver) OsPackages(ctx context.Context, obj *model.EC2Instance) ([]*model.OSPackage, error) {
	logrus.Debugf("getting os packages for instance %v", obj.ID)
	var client osquery.Client = nil
	client = r.Resolver.Cache.Get(obj.ID)
	if client == nil {
		logrus.Debugf("no client for instance %v", obj.ID)
		var osqueryHost string
		if obj.Public {
			osqueryHost = "http://" + obj.PublicIP
		} else {
			osqueryHost = "http://" + obj.PrivateIP
		}

		osqueryConfig := &osquery.ClientOpts{
			Host: osqueryHost,
			EC2SSHConfig: &osquery.OSQueryEC2SSHConfig{
				ID:        obj.ID,
				Region:    obj.Region,
				AZ:        obj.AvailabilityZone,
				IsPublic:  obj.Public,
				PublicIP:  obj.PublicIP,
				PrivateIP: obj.PrivateIP,
			},
		}
		var err error = nil
		client, err = osquery.NewClient(osqueryConfig)
		if err != nil {
			logrus.Errorf("unable to create osquery client for ec2 instance %v: %v", obj.ID, err.Error())
		} else {
			r.Resolver.Cache.Put(obj.ID, client)
		}
	}

	if client == nil {
		return nil, fmt.Errorf("no client could be found for instance %v", obj.ID)
	}

	packages, pkgError := client.GetOSPackages()
	if pkgError != nil {
		logrus.Errorf("error getting os packages for instance %v: %v", obj.ID, pkgError.Error())
	}

	returnPkgs := make([]*model.OSPackage, len(packages))
	for idx, pkg := range packages {
		returnPkgs[idx] = &pkg
	}

	return returnPkgs, pkgError
}

func (r *eC2InstanceResolver) ListeningApplications(ctx context.Context, obj *model.EC2Instance) ([]*model.ListeningApplication, error) {
	logrus.Debugf("getting listening applications for instance %v", obj.ID)
	var client osquery.Client = nil
	client = r.Resolver.Cache.Get(obj.ID)
	if client == nil {
		logrus.Debugf("no client for instance %v", obj.ID)
		var osqueryHost string
		if obj.Public {
			osqueryHost = "http://" + obj.PublicIP
		} else {
			osqueryHost = "http://" + obj.PrivateIP
		}

		osqueryConfig := &osquery.ClientOpts{
			Host: osqueryHost,
			EC2SSHConfig: &osquery.OSQueryEC2SSHConfig{
				ID:        obj.ID,
				Region:    obj.Region,
				AZ:        obj.AvailabilityZone,
				IsPublic:  obj.Public,
				PublicIP:  obj.PublicIP,
				PrivateIP: obj.PrivateIP,
			},
		}
		var err error = nil
		client, err = osquery.NewClient(osqueryConfig)
		if err != nil {
			logrus.Errorf("unable to create osquery client for ec2 instance %v: %v", obj.ID, err.Error())
		} else {
			r.Resolver.Cache.Put(obj.ID, client)
		}
	}

	if client == nil {
		return nil, fmt.Errorf("no client could be found for instance %v", obj.ID)
	}

	apps, appError := client.GetListeningApplications()
	if appError != nil {
		logrus.Errorf("error getting listening applications for instance %v: %v", obj.ID, appError.Error())
	}

	returnApps := make([]*model.ListeningApplication, len(apps))
	for idx, app := range apps {
		returnApps[idx] = &app
	}

	return returnApps, appError
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
					osqueryHost = "http://" + instance.PublicIP
				} else {
					osqueryHost = "http://" + instance.PrivateIP
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
