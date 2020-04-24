package network

type Network struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
}

type Edge struct {
	NodeA string `json:"node_a"`
	NodeB string `json:"node_b"`
}
