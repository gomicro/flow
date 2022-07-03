package rds

import (
	gofmt "fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"

	"github.com/gomicro/flow/fmt"
)

func init() {
}

// ListCmd represents the list rds instances
var instancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List db instances",
	Long:  `List db instances from RDS`,
	Run:   instancesListFunc,
}

var padding = &paddings{}

func instancesListFunc(cmd *cobra.Command, args []string) {
	input := &rds.DescribeDBInstancesInput{}

	table := []*tableRow{}

	err := rdsSvc.DescribeDBInstancesPages(input,
		func(page *rds.DescribeDBInstancesOutput, lastPage bool) bool {
			for _, in := range page.DBInstances {
				name := *in.DBName
				if len(name)+1 > padding.name {
					padding.name = len(name) + 1
				}

				status := *in.DBInstanceStatus
				if len(status)+1 > padding.status {
					padding.status = len(status) + 1
				}

				region := *in.AvailabilityZone
				if len(region)+1 > padding.region {
					padding.region = len(region) + 1
				}

				size := *in.DBInstanceClass
				if len(size)+1 > padding.size {
					padding.size = len(size) + 1
				}

				endpoint := *in.Endpoint.Address
				if len(endpoint)+1 > padding.endpoint {
					padding.endpoint = len(endpoint) + 1
				}

				row := &tableRow{
					name:     name,
					status:   status,
					region:   region,
					size:     size,
					endpoint: endpoint,
				}

				table = append(table, row)
			}

			return !lastPage
		})
	if err != nil {
		fmt.Printf("Error listing instances: %s", err)
		os.Exit(1)
	}

	gofmt.Printf("%-*s | %-*s | %-*s | %-*s | %-*s\n",
		padding.name, "name",
		padding.status, "status",
		padding.region, "region",
		padding.size, "size",
		padding.endpoint, "endpont",
	)
	gofmt.Printf("%s | %s | %s | %s | %s\n",
		strings.Repeat("-", padding.name),
		strings.Repeat("-", padding.status),
		strings.Repeat("-", padding.region),
		strings.Repeat("-", padding.size),
		strings.Repeat("-", padding.endpoint),
	)

	for _, r := range table {
		gofmt.Printf("%s\n", r)
	}
}

type paddings struct {
	name     int
	status   int
	region   int
	size     int
	endpoint int
}

type tableRow struct {
	name     string
	status   string
	region   string
	size     string
	endpoint string
}

func (r *tableRow) String() string {
	return gofmt.Sprintf("%-*s | %-*s | %-*s | %-*s | %-*s",
		padding.name, r.name,
		padding.status, r.status,
		padding.region, r.region,
		padding.size, r.size,
		padding.endpoint, r.endpoint,
	)
}
