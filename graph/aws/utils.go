package aws

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

// GetAllRegions returns all of the AWS EC2 available regions
func GetAllRegions(sess *session.Session) ([]string, error) {
	svc := ec2.New(sess)
	input := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false),
	}

	result, err := svc.DescribeRegions(input)
	if err != nil {
		return []string{}, err
	}

	regions := make([]string, len(result.Regions))
	for rIdx, region := range result.Regions {
		regions[rIdx] = *region.RegionName
	}

	return regions, nil
}

func GetRegionalSession(sess *session.Session, region string) *session.Session {
	// Create a Session with a custom region
	regionalSession := sess.Copy(&aws.Config{
		Region: aws.String(region),
	})

	return regionalSession
}

func NewEC2SSHSession(sess *session.Session, instance model.EC2Instance) (*ssh.Client, error) {
	possibleUserNames := []string{"ec2-user", "ubuntu"}
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	logrus.Debugf("priv: %#v\n", key)
	logrus.Debugf("pub: %#v\n", key.Public())

	sshPubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)

	logrus.Debugf("ssh: %#v\n", string(sshPubKey.Marshal()))
	logrus.Debugf("ssh: %#v\n", string(sshPubKeyBytes))

	svc := ec2instanceconnect.New(sess)
	var sshClient *ssh.Client
	for _, user := range possibleUserNames {
		input := &ec2instanceconnect.SendSSHPublicKeyInput{
			AvailabilityZone: aws.String(instance.AvailabilityZone),
			InstanceId:       aws.String(instance.ID),
			InstanceOSUser:   aws.String(user),
			SSHPublicKey:     aws.String(string(sshPubKeyBytes)),
		}

		result, err := svc.SendSSHPublicKey(input)
		if err != nil {
			logrus.Error(err.Error())
			break
		}
		logrus.Debugf("ec2 instance connect for instance %s has result: %#v\n", instance.ID, result)

		// start ssh session
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			logrus.Error(err.Error())
			break
		}

		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		if instance.Public {
			sshClient, err = ssh.Dial("tcp", instance.PublicIP+":22", config)
			if err != nil {
				logrus.Error(err.Error())
				break
			}
		} else {
			sshClient, err = ssh.Dial("tcp", instance.PrivateIP+":22", config)
			if err != nil {
				logrus.Error(err.Error())
				break
			}
		}
		logrus.Debugf("ssh session: %#v", *sshClient)
		return sshClient, nil
	}

	return nil, errors.New("no ssh session was completed")
}

func parseTime(layout, value string) *time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return &t
}
