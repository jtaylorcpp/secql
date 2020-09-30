package aws

import (
	"strings"

	awsTypes "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/sirupsen/logrus"
)

func GetAllEC2Instances(regionalSess *session.Session, region string) ([]*model.EC2Instance, error) {
	instanceModels := []*model.EC2Instance{}
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
						Region:           region,
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
					/*sshClient, sshError := aws.NewEC2SSHSession(regionalSess, *instanceModel)
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

					osqPackages, err := osquery.GetPackages(sshClient, osqOSInfo)
					if err != nil {
						logrus.Errorf("got error from osquery package discovery: %s", err.Error())
					}

					packages := make([]*model.OSPackage, len(osqPackages))
					for idx, pkg := range osqPackages {
						packages[idx] = &model.OSPackage{
							ID:         pkg.Name,
							Version:    pkg.Version,
							Source:     pkg.Source,
							Size:       pkg.Size,
							Arch:       pkg.Arch,
							Revision:   pkg.Revision,
							Status:     pkg.Status,
							Maintainer: pkg.Maintainer,
							Section:    pkg.Section,
							Priority:   pkg.Priority,
						}
					}

					instanceModel.OsPackages = packages

					listeningApps, err := osquery.GetListeningApplications(sshClient, osqOSInfo)
					if err != nil {
						logrus.Errorf("got error from osquery listener discovery: %s", err.Error())
					}
					listeningApplications := make([]*model.ListeningApplication, len(listeningApps))
					for idx, app := range listeningApps {
						listeningApplications[idx] = &model.ListeningApplication{
							ID:      app.Name,
							Address: app.Address,
							Port:    app.Port,
							Pid:     app.Pid,
						}
					}
					instanceModel.ListeningApplications = listeningApplications*/
					instanceModels = append(instanceModels, instanceModel)
				}
			}
			return !lastPage
		})
	if err != nil {
		logrus.Errorf("error describing ec2 instances: %v", err.Error())
	}

	return instanceModels, err
}
