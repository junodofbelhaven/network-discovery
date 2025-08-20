package ports

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"network-discovery/internal/models"

	"github.com/sirupsen/logrus"
)

// Scanner wraps Nmap CLI execution to detect open ports for given hosts
type Scanner struct {
	MaxWorkers     int
	TimeoutPerHost time.Duration
	logger         *logrus.Logger
}

func NewScanner(maxWorkers int) *Scanner {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	return &Scanner{
		MaxWorkers:     maxWorkers,
		TimeoutPerHost: 30 * time.Second,
		logger:         logger,
	}
}

func NewScannerWithLogger(maxWorkers int, logger *logrus.Logger) *Scanner {
	return &Scanner{
		MaxWorkers:     maxWorkers,
		TimeoutPerHost: 30 * time.Second,
		logger:         logger,
	}
}

// ScanHost runs nmap for a single IP and returns open ports
func (s *Scanner) ScanHost(ip string) ([]models.PortInfo, error) {
	if ip == "" {
		return nil, fmt.Errorf("empty ip")
	}

	// Build command: fast, no DNS, open ports only, XML to stdout
	args := []string{"-Pn", "-T4", "-n", "--open", "-oX", "-", ip}
	ctx, cancel := context.WithTimeout(context.Background(), s.TimeoutPerHost)
	defer cancel()

	cmd := exec.CommandContext(ctx, "nmap", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// If nmap is missing or times out, return error to let caller decide
		return nil, fmt.Errorf("nmap run failed: %v, stderr: %s", err, stderr.String())
	}

	ports, err := parseNmapXML(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("parse nmap xml failed: %v", err)
	}
	return ports, nil
}

// Minimal XML structures to extract open ports
type nmapRun struct {
	XMLName xml.Name   `xml:"nmaprun"`
	Hosts   []nmapHost `xml:"host"`
}

type nmapHost struct {
	Ports nmapPorts `xml:"ports"`
}

type nmapPorts struct {
	Ports []nmapPort `xml:"port"`
}

type nmapPort struct {
	Protocol string      `xml:"protocol,attr"`
	PortID   string      `xml:"portid,attr"`
	State    nmapState   `xml:"state"`
	Service  nmapService `xml:"service"`
}

type nmapState struct {
	State string `xml:"state,attr"`
}

type nmapService struct {
	Name string `xml:"name,attr"`
}

func parseNmapXML(data []byte) ([]models.PortInfo, error) {
	var run nmapRun
	if err := xml.Unmarshal(data, &run); err != nil {
		return nil, err
	}
	var result []models.PortInfo
	for _, h := range run.Hosts {
		for _, p := range h.Ports.Ports {
			if p.State.State != "open" {
				continue
			}
			portNum, _ := strconv.Atoi(p.PortID)
			result = append(result, models.PortInfo{
				Port:     portNum,
				Protocol: p.Protocol,
				Service:  p.Service.Name,
				State:    p.State.State,
			})
		}
	}
	return result, nil
}
