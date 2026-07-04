package rancher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const requestTimeout = 10 * time.Second

type DownstreamCluster struct {
	ID   string
	Name string
}

type rancherClustersResponse struct {
	Data []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func FetchDownstreamClusters(server, token string) ([]DownstreamCluster, error) {
	client := &http.Client{Timeout: requestTimeout}

	req, err := http.NewRequest("GET", strings.TrimRight(server, "/")+"/v3/clusters", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result rancherClustersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var clusters []DownstreamCluster
	for _, c := range result.Data {
		clusters = append(clusters, DownstreamCluster{ID: c.ID, Name: c.Name})
	}
	return clusters, nil
}

func FetchClusterKubeConfig(server, token, clusterID string) ([]byte, error) {
	client := &http.Client{Timeout: requestTimeout}

	req, err := http.NewRequest("POST", strings.TrimRight(server, "/")+"/v3/clusters/"+clusterID+"?action=generateKubeconfig", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Config string `json:"config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return []byte(result.Config), nil
}
