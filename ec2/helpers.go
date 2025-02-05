package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

func GetEC2InstanceInfo(client ec2iface.EC2API, instances []*string) (output []*ec2.Instance, err error) {
	// Set up our DI input object
	diInput := &ec2.DescribeInstancesInput{
		InstanceIds: instances,
	}

	describeInstancesPager := func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range page.Reservations {
			output = append(output, reservation.Instances...)
		}

		// If it's not the last page, continue
		return !lastPage
	}

	// Fetch all the instances described
	if err = client.DescribeInstancesPages(diInput, describeInstancesPager); err != nil {
		return nil, fmt.Errorf("Could not describe EC2 instances\n%v", err)
	}

	return output, nil
}

// GetEC2InstanceTags accepts any number of instance strings and returns a populated InstanceTags{} object for each instance
func GetEC2InstanceTags(client ec2iface.EC2API, instances []*string) (ec2Tags map[string]Tags, err error) {
	instanceInfo, err := GetEC2InstanceInfo(client, instances)
	if err != nil {
		return nil, fmt.Errorf("Error when trying to retrieve EC2 instance tags\n%v", err)
	}

	ec2Tags = make(map[string]Tags)
	for _, i := range instanceInfo {
		tagMap := make(map[string]string)

		for _, tag := range i.Tags {
			tagMap[*tag.Key] = *tag.Value
		}

		ec2Tags[*i.InstanceId] = tagMap
	}

	return ec2Tags, nil
}
