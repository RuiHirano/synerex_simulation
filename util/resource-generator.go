package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/go-yaml/yaml"
)

type Resource struct {
	ApiVersion string   `yaml:"apiVersion,omitempty"`
	Kind       string   `yaml:"kind,omitempty"`
	Metadata   Metadata `yaml:"metadata,omitempty"`
	Spec       Spec     `yaml:"spec,omitempty"`
}

type Spec struct {
	Replicas int      `yaml:"replicas,omitempty"`
	Templete Templete `yaml:"templete,omitempty"`
	Selector Selector `yaml:"selector,omitempty"`
	Ports    []Port   `yaml:"ports,omitempty"`
	Type     string   `yaml:"type,omitempty"`
}

type Templete struct {
	Metadata Metadata     `yaml:"metadata,omitempty"`
	Spec     TempleteSpec `yaml:"spec,omitempty"`
}

type TempleteSpec struct {
	Containers []Container `yaml:"containers,omitempty"`
}

type Container struct {
	Name            string `yaml:"name,omitempty"`
	Image           string `yaml:"image,omitempty"`
	ImagePullPolicy string `yaml:"imagePullPolicy,omitempty"`
	Stdin           bool   `yaml:"stdin,omitempty"`
	Tty             bool   `yaml:"tty,omitempty"`
	Env             []Env  `yaml:"env,omitempty"`
	Ports           []Port `yaml:"ports,omitempty"`
}

type Env struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Selector struct {
	App         string  `yaml:"app,omitempty"`
	MatchLabels []Label `yaml:"matchLabels,omitempty"`
}

type Port struct {
	Name          string `yaml:"name,omitempty"`
	Port          int    `yaml:"port,omitempty"`
	TargetPort    int    `yaml:"targetPort,omitempty"`
	ContainerPort int    `yaml:"containerPort,omitempty"`
}

type Metadata struct {
	Name   string  `yaml:"name,omitempty"`
	Labels []Label `yaml:"labels,omitempty"`
}

type Label struct {
	App string `yaml:"app,omitempty"`
}

type Area struct {
	Control   []Coord
	Duplicate []Coord
}

type Coord struct {
	Latitude  float64
	Longitude float64
}

type Option struct {
	FileName        string
	AreaCoords      []Coord
	DevideSquareNum int
	DuplicateRate   float64
}

// vis-monitor
func NewVisMonitorService() Resource {
	service := Resource{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   Metadata{Name: "vis-monitor"},
		Spec: Spec{
			Selector: Selector{App: "worker"},
			Ports: []Port{
				{
					Name:       "http",
					Port:       80,
					TargetPort: 9500,
				},
			},
			Type: "NodePort",
		},
	}
	return service
}

// worker
func NewWorkerService() Resource {
	service := Resource{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   Metadata{Name: "worker"},
		Spec: Spec{
			Selector: Selector{App: "worker"},
			Ports: []Port{
				{
					Name:       "synerex",
					Port:       700,
					TargetPort: 10000,
				},
				{
					Name:       "nodeid",
					Port:       600,
					TargetPort: 9000,
				},
			},
		},
	}
	return service
}

func NewWorker() Resource {
	worker := Resource{
		ApiVersion: "apps/v1",
		Kind:       "ReplicaSet",
		Metadata: Metadata{
			Name:   "worker",
			Labels: []Label{{App: "worker"}},
		},
		Spec: Spec{
			Replicas: 1,
			Selector: Selector{MatchLabels: []Label{{App: "worker"}}},
			Templete: Templete{
				Metadata: Metadata{Labels: []Label{{App: "worker"}}},
				Spec: TempleteSpec{
					Containers: []Container{
						{
							Name:            "nodeid-server",
							Image:           "synerex-simulation/nodeid-server:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
							},
							Ports: []Port{{ContainerPort: 9000}},
						},
						{
							Name:            "synerex-server",
							Image:           "synerex-simulation/synerex-server:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
								{
									Name:  "SYNEREX_SERVER",
									Value: ":10000",
								},
							},
							Ports: []Port{{ContainerPort: 10000}},
						},
						{
							Name:            "worker-provider",
							Image:           "synerex-simulation/worker-provider:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
								{
									Name:  "SYNEREX_SERVER",
									Value: ":10000",
								},
								{
									Name:  "MASTER_SYNEREX_SERVER",
									Value: "master:700",
								},
								{
									Name:  "MASTER_NODEID_SERVER",
									Value: "master:600",
								},
								{
									Name:  "PORT",
									Value: "9980",
								},
							},
							Ports: []Port{{ContainerPort: 9980}},
						},
						{
							Name:            "agent-provider",
							Image:           "synerex-simulation/agent-provider:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
								{
									Name:  "SYNEREX_SERVER",
									Value: ":10000",
								},
								{
									Name:  "AREA",
									Value: "{\"id\":0, \"name\":\"Nagoya\", \"duplicate_area\": [{\"latitude\":35.15515326666666, \"longitude\":136.97152533333332},{\"latitude\":35.15097806666667, \"longitude\":136.97152533333332},{\"latitude\":35.15097806666667, \"longitude\":136.9788893333333},{\"latitude\":35.15515326666666, \"longitude\":136.9788893333333}], \"control_area\": [{\"latitude\":35.15480533333333, \"longitude\":136.972139},{\"latitude\":35.151326, \"longitude\":136.972139},{\"latitude\":35.151326, \"longitude\":136.97827566666666},{\"latitude\":35.15480533333333, \"longitude\":136.97827566666666}]}",
								},
							},
						},
						{
							Name:            "visualization-provider",
							Image:           "synerex-simulation/visualization-provider:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
								{
									Name:  "SYNEREX_SERVER",
									Value: ":10000",
								},
								{
									Name:  "VIS_ADDRESS",
									Value: ":9500",
								},
							},
						},
					},
				},
			},
		},
	}
	return worker
}

// master
func NewMasterService() Resource {
	service := Resource{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   Metadata{Name: "master"},
		Spec: Spec{
			Selector: Selector{App: "master"},
			Ports: []Port{
				{
					Name:       "synerex",
					Port:       700,
					TargetPort: 10000,
				},
				{
					Name:       "nodeid",
					Port:       600,
					TargetPort: 9000,
				},
				{
					Name:       "master-provider",
					Port:       800,
					TargetPort: 9990,
				},
			},
		},
	}
	return service
}

func NewMaster() Resource {
	master := Resource{
		ApiVersion: "apps/v1",
		Kind:       "ReplicaSet",
		Metadata: Metadata{
			Name:   "master",
			Labels: []Label{{App: "master"}},
		},
		Spec: Spec{
			Replicas: 1,
			Selector: Selector{MatchLabels: []Label{{App: "master"}}},
			Templete: Templete{
				Metadata: Metadata{Labels: []Label{{App: "master"}}},
				Spec: TempleteSpec{
					Containers: []Container{
						{
							Name:            "nodeid-server",
							Image:           "synerex-simulation/nodeid-server:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
							},
							Ports: []Port{{ContainerPort: 9000}},
						},
						{
							Name:            "synerex-server",
							Image:           "synerex-simulation/synerex-server:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
								{
									Name:  "SYNEREX_SERVER",
									Value: ":10000",
								},
							},
							Ports: []Port{{ContainerPort: 10000}},
						},
						{
							Name:            "master-provider",
							Image:           "synerex-simulation/master-provider:latest",
							ImagePullPolicy: "Never",
							Env: []Env{
								{
									Name:  "NODEID_SERVER",
									Value: ":9000",
								},
								{
									Name:  "SYNEREX_SERVER",
									Value: ":10000",
								},
								{
									Name:  "PORT",
									Value: "9990",
								},
							},
							Ports: []Port{{ContainerPort: 9990}},
						},
					},
				},
			},
		},
	}
	return master
}

// simulator
func NewSimulatorService() Resource {
	service := Resource{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   Metadata{Name: "simulator"},
		Spec: Spec{
			Selector: Selector{App: "simulator"},
			Ports: []Port{
				{
					Name: "http",
					Port: 8000,
				},
			},
			Type: "NodePort",
		},
	}
	return service
}

func NewSimulator() Resource {
	simulator := Resource{
		ApiVersion: "apps/v1",
		Kind:       "ReplicaSet",
		Metadata: Metadata{
			Name:   "simulator",
			Labels: []Label{{App: "simulator"}},
		},
		Spec: Spec{
			Replicas: 1,
			Selector: Selector{MatchLabels: []Label{{App: "simulator"}}},
			Templete: Templete{
				Metadata: Metadata{Labels: []Label{{App: "simulator"}}},
				Spec: TempleteSpec{
					Containers: []Container{
						{
							Name:            "simulator",
							Image:           "synerex-simulation/simulator:latest",
							ImagePullPolicy: "Never",
							Stdin:           true,
							Tty:             true,
							Env: []Env{
								{
									Name:  "MASTER_ADDRESS",
									Value: "http://master:800",
								},
							},
							Ports: []Port{{ContainerPort: 8000}},
						},
					},
				},
			},
		},
	}
	return simulator
}

// gateway
func NewGateway() Resource {
	master := Resource{}
	return master
}

func main() {

	option := Option{
		FileName: "test.yaml",
		AreaCoords: []Coord{
			{Longitude: 136.972139, Latitude: 35.161764},
			{Longitude: 136.972139, Latitude: 35.151326},
			{Longitude: 136.990549, Latitude: 35.151326},
			{Longitude: 136.990549, Latitude: 35.161764},
		},
		DevideSquareNum: 3,
		DuplicateRate:   0.1,
	}

	rsrcs := createData(option)
	fmt.Printf("rsrcs: %v\n", rsrcs)

	// write yaml
	fileName := option.FileName
	for _, rsrc := range rsrcs {
		err := WriteOnFile(fileName, rsrc)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func createData(option Option) []Resource {
	rsrcs := []Resource{
		NewSimulator(),
		NewSimulatorService(),
		NewMaster(),
		NewMasterService(),
	}
	areas := AreaDivider(option.AreaCoords, option.DevideSquareNum, option.DuplicateRate)
	fmt.Printf("areas: %v\n", areas)
	return rsrcs
}

// ファイル名とデータをを渡すとyamlファイルに保存してくれる関数です。
func WriteOnFile(fileName string, data interface{}) error {
	// ここでデータを []byte に変換しています。
	buf, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, string(buf))   //書き込み
	fmt.Fprintln(file, string("---")) //書き込み
	return nil
}

func AreaDivider(areaCoords []Coord, divideSquareNum int, duplicateRate float64) []Area {
	var areas []Area
	maxLat, maxLon, minLat, minLon := GetCoordRange(areaCoords)
	for i := 0; i < divideSquareNum; i++ { // 横方向
		leftlon := (maxLon-minLon)*float64(i)/float64(divideSquareNum) + minLon
		rightlon := (maxLon-minLon)*(float64(i)+1)/float64(divideSquareNum) + minLon

		for k := 0; k < divideSquareNum; k++ { // 縦方向
			bottomlat := (maxLat-minLat)*float64(i)/float64(divideSquareNum) + minLat
			toplat := (maxLat-minLat)*(float64(i)+1)/float64(divideSquareNum) + minLat
			area := Area{
				Control: []Coord{
					{Longitude: leftlon, Latitude: toplat},
					{Longitude: leftlon, Latitude: bottomlat},
					{Longitude: rightlon, Latitude: bottomlat},
					{Longitude: rightlon, Latitude: toplat},
				},
				Duplicate: []Coord{
					{Longitude: leftlon - (rightlon-leftlon)*duplicateRate, Latitude: toplat + (toplat-bottomlat)*duplicateRate},
					{Longitude: leftlon - (rightlon-leftlon)*duplicateRate, Latitude: bottomlat - (toplat-bottomlat)*duplicateRate},
					{Longitude: rightlon + (rightlon-leftlon)*duplicateRate, Latitude: bottomlat - (toplat-bottomlat)*duplicateRate},
					{Longitude: rightlon + (rightlon-leftlon)*duplicateRate, Latitude: toplat + (toplat-bottomlat)*duplicateRate},
				},
			}
			areas = append(areas, area)
		}
	}

	return areas

}

func GetCoordRange(coords []Coord) (float64, float64, float64, float64) {
	maxLon, maxLat := math.Inf(-1), math.Inf(-1)
	minLon, minLat := math.Inf(0), math.Inf(0)
	for _, coord := range coords {
		if coord.Latitude > maxLat {
			maxLat = coord.Latitude
		}
		if coord.Longitude > maxLon {
			maxLon = coord.Longitude
		}
		if coord.Latitude < minLat {
			minLat = coord.Latitude
		}
		if coord.Longitude < minLon {
			minLon = coord.Longitude
		}
	}
	return maxLat, maxLon, minLat, minLon
}
