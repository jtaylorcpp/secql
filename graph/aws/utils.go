package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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

func NewEC2SSHKey() {

}
