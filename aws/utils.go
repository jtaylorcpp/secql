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
	hellossh "github.com/helloyi/go-sshclient"
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

func NewEC2SSHSession(sess *session.Session, isPublic bool, id, az, publicIP, privateIP string) (*hellossh.Client, error) {
	possibleUserNames := []string{"ec2-user", "ubuntu"}
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	sshPubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)

	logrus.Debugf("ssh: %#v\n", string(sshPubKeyBytes))

	svc := ec2instanceconnect.New(sess)
	var sshClient *hellossh.Client
	for _, user := range possibleUserNames {
		input := &ec2instanceconnect.SendSSHPublicKeyInput{
			AvailabilityZone: aws.String(az),
			InstanceId:       aws.String(id),
			InstanceOSUser:   aws.String(user),
			SSHPublicKey:     aws.String(string(sshPubKeyBytes)),
		}

		result, err := svc.SendSSHPublicKey(input)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		logrus.Debugf("ec2 instance connect for user %s instance %s has result: %#v\n", user, id, result)

		// start ssh session
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		if isPublic {
			sshClient, err = hellossh.Dial("tcp", publicIP+":22", config)
			if err != nil {
				logrus.Error(err.Error())
				continue
			}
		} else {
			sshClient, err = hellossh.Dial("tcp", privateIP+":22", config)
			if err != nil {
				logrus.Error(err.Error())
				continue
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
