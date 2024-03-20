package spec

import (
	"context"
	"testing"

	"github.com/ermes-labs/api-go/api"
	"github.com/ermes-labs/api-go/infrastructure"
)

var infra = infrastructure.Infrastructure{
	Areas: []infrastructure.Area{
		{
			Node: infrastructure.Node{
				AreaName: "area",
				Host:     "host",
				GeoCoordinates: infrastructure.GeoCoordinates{
					Latitude:  0,
					Longitude: 0,
				},
			},
			Areas: []infrastructure.Area{
				{
					Node: infrastructure.Node{
						AreaName: "area1",
						Host:     "host1",
						GeoCoordinates: infrastructure.GeoCoordinates{
							Latitude:  50,
							Longitude: 20,
						},
					},
				},
				{
					Node: infrastructure.Node{
						AreaName: "area2",
						Host:     "host2",
						GeoCoordinates: infrastructure.GeoCoordinates{
							Latitude:  -40,
							Longitude: -30,
						},
					},
					Areas: []infrastructure.Area{
						{
							Node: infrastructure.Node{
								AreaName: "area3",
								Host:     "host1",
								GeoCoordinates: infrastructure.GeoCoordinates{
									Latitude:  40,
									Longitude: 100,
								},
							},
						},
					},
				},
			},
		},
	},
}

func TestInfrastructure[T api.Commands](t *testing.T, env Env[T]) {
	// Set up the environment.
	cmd, free := env.New("area2")
	defer free()

	err := cmd.LoadInfrastructure(context.Background(), infra)
	if err != nil {
		t.Errorf("failed to load the infrastructure: %v", err)
	}

	// Get the infrastructure.
	nodes, err := cmd.GetChildrenNodesOf(context.Background(), "area")

	if err != nil {
		t.Errorf("failed to get the children nodes of area: %v", err)
	}

	if len(nodes) != 2 {
		t.Errorf("invalid children nodes of area, expected 2, found %v", nodes)
	}

	node, err := cmd.GetParentNodeOf(context.Background(), "area1")

	if err != nil {
		t.Errorf("failed to get the parent node of area1: %v", err)
	}

	if node.AreaName != "area" {
		t.Errorf("invalid parent node of area1, expected %v, found %v", "area", node.AreaName)
	}

	node, err = cmd.GetParentNodeOf(context.Background(), "area")

	if err != nil {
		t.Errorf("failed to get the parent node of area: %v", err)
	}

	if node != nil {
		t.Errorf("invalid parent node of area, expected nil, found %v", node)
	}

	sessionId, err := cmd.CreateSession(context.Background(), api.NewCreateSessionOptionsBuilder().ClientGeoCoordinates(infrastructure.GeoCoordinates{
		Latitude:  0,
		Longitude: 0,
	}).Build())

	if err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	lookupNode, err := cmd.FindLookupNode(context.Background(), sessionId)

	if err != nil {
		t.Errorf("failed to find the lookup node: %v", err)
	}

	if lookupNode.AreaName != "area" {
		t.Errorf("invalid lookup node, expected %v, found %v", infra.Areas[0].Node, node)
	}

	sessionId, err = cmd.CreateSession(context.Background(), api.NewCreateSessionOptionsBuilder().Build())

	if err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	lookupNode, err = cmd.FindLookupNode(context.Background(), sessionId)

	if err != nil {
		t.Errorf("failed to find the lookup node: %v", err)
	}

	if lookupNode.AreaName != "area2" {
		t.Errorf("invalid lookup node, expected %v, found %v", "area1", node)
	}
}
