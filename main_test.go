package main

import (
	"reflect"
	"testing"
)

func TestParseDockerImagesOutput_Realistic(t *testing.T) {
	output := `REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
repo1               latest              111111111111        2 weeks ago         1.5GB
repo2               v1                  222222222222        3 days ago          200MB
repo3               <none>              333333333333        5 hours ago         10KB
repo4               v2                  444444444444        10 minutes ago      0.5GB
`
	images := parseDockerImagesOutput(output)
	expected := []DockerImage{
		{
			Repository: "repo1",
			Tag:        "latest",
			ImageID:    "111111111111",
			Created:    "2 weeks ago",
			Size:       "1.5GB",
			SizeBytes:  1500000000,
		},
		{
			Repository: "repo4",
			Tag:        "v2",
			ImageID:    "444444444444",
			Created:    "10 minutes ago",
			Size:       "0.5GB",
			SizeBytes:  500000000,
		},
		{
			Repository: "repo2",
			Tag:        "v1",
			ImageID:    "222222222222",
			Created:    "3 days ago",
			Size:       "200MB",
			SizeBytes:  200000000,
		},
		{
			Repository: "repo3",
			Tag:        "<none>",
			ImageID:    "333333333333",
			Created:    "5 hours ago",
			Size:       "10KB",
			SizeBytes:  10000,
		},
	}
	if !reflect.DeepEqual(images, expected) {
		t.Errorf("parseDockerImagesOutput() = %#v; want %#v", images, expected)
	}
}

func TestParseSize_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1.5GB", 1500000000},
		{"200MB", 200000000},
		{"10KB", 10000},
		{"0.5GB", 500000000},
		{"0GB", 0},
		{"0MB", 0},
		{"0KB", 0},
		{"bad", 0},
		{"123", 0},
		{"123B", 0},
		{"", 0},
	}

	for _, tt := range tests {
		got := parseSize(tt.input)
		if got != tt.expected {
			t.Errorf("parseSize(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}

func TestFilterImages(t *testing.T) {
	images := []DockerImage{
		{Repository: "repo1", Tag: "latest", ImageID: "111111111111"},
		{Repository: "repo2", Tag: "v1", ImageID: "222222222222"},
		{Repository: "repo3", Tag: "<none>", ImageID: "333333333333"},
	}
	tests := []struct {
		filter   string
		expected []DockerImage
	}{
		{"repo1", []DockerImage{images[0]}},
		{"2222", []DockerImage{images[1]}},
		{"repo", images},
		{"", images},
		{"notfound", nil},
	}

	for _, tt := range tests {
		filtered := filterImages(images, tt.filter)
		if !reflect.DeepEqual(filtered, tt.expected) {
			t.Errorf("filter=%q: got %#v; want %#v", tt.filter, filtered, tt.expected)
		}
	}
}
