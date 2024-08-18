package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

// Define struct containing context and client
type config struct {
	ctx       context.Context
	client    *admin.APIClient
	projectID string
	gin       *gin.Engine
}

func newClient() (*config, error) {
	ctx := context.Background()
	_, b := os.LookupEnv("GIN_MODE")
	if !b {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	apiKey := os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
	apiSecret := os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("MONGODB_ATLAS_PUBLIC_KEY and MONGODB_ATLAS_PRIVATE_KEY must be set")
	}
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("MONGODB_ATLAS_PROJECT_ID must be set")
	}
	sdk, err := admin.NewClient(admin.UseDigestAuth(apiKey, apiSecret))
	if err != nil {
		return nil, err
	}
	r.SetTrustedProxies(nil)
	return &config{ctx: ctx, client: sdk, projectID: projectID, gin: r}, nil
}

func main() {
	client, err := newClient()
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not fetch IP access list: %v", err)
	}
	client.gin.POST("/api/v1/updateIP", func(c *gin.Context) {
		// Parse JSON input
		var input struct {
			IP      string `json:"ip"`
			Comment string `json:"comment"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// Check if IP already exists in access list
		if client.checkExistingIP(input.IP) {
			c.JSON(200, gin.H{"message": "IP already exists in access list"})
			return
		}
		// Create IP access list
		_, err := client.CreateIPAccessList(input.IP, input.Comment)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "IP added to access list"})
		// Remove old entries
		client.RemoveOldEntries(input.IP)
	})
	client.gin.Run(":8080")
}

// Check if IP already exists in the access list
func (c *config) checkExistingIP(ip string) bool {
	ipAccessList, err := c.ListIPAccessList()
	if err != nil {
		log.Fatalf("Could not fetch IP access list: %v", err)
	}
	for _, entry := range *ipAccessList.Results {
		if *entry.CidrBlock == ip {
			fmt.Printf("IP: %v already exists in accesslist, doing nothing.\n", *entry.CidrBlock)
			return true
		}
	}
	return false
}

func (c *config) ListIPAccessList() (*admin.PaginatedNetworkAccess, error) {
	ipAccessList, _, err := c.client.ProjectIPAccessListApi.ListProjectIpAccessLists(c.ctx, c.projectID).Execute()
	if err != nil {
		log.Fatalf("Could not fetch IP access list: %v", err)
	}
	return ipAccessList, nil
}

func (c *config) CreateIPAccessList(cidr string, comment string) (*admin.PaginatedNetworkAccess, error) {
	entry := &[]admin.NetworkPermissionEntry{
		{CidrBlock: &cidr,
			Comment: &comment},
	}
	admin, resp, err := c.client.ProjectIPAccessListApi.CreateProjectIpAccessList(c.ctx, c.projectID, entry).Execute()
	if err != nil {
		log.Fatalf("Could not create IP access list: %v (%v)", err, resp.Request.Response)
		return nil, err
	}
	return admin, nil
}

func (c *config) RemoveOldEntries(ip string) error {
	ipAccessList, err := c.ListIPAccessList()
	if err != nil {
		log.Fatalf("Could not fetch IP access list: %v", err)
	}
	for _, entry := range *ipAccessList.Results {
		if *entry.CidrBlock != ip {
			_, resp, err := c.client.ProjectIPAccessListApi.DeleteProjectIpAccessList(c.ctx, c.projectID, *entry.CidrBlock).Execute()
			if err != nil {
				log.Fatalf("Could not delete IP access list: %v (%v)", err, resp.Request.Response)
				return err
			}
		}
	}
	return nil
}
