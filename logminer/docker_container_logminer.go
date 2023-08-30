package logminer

import (
	"bufio"
	"github.com/fsouza/go-dockerclient"
	"io"
	"log"
	"nginx-proxy-metrics/config"
	"nginx-proxy-metrics/logline"
	"nginx-proxy-metrics/persistence"
	"time"
)

const (
	DockerDaemonSocket = "unix:///var/run/docker.sock"
	LogLineDelimiter   = '\n'
)

type DockerContainerLogMiner struct {
	HttpRequestPersister persistence.HttpRequestPersister
	LoglineParser        logline.Parser
}

func (logMiner *DockerContainerLogMiner) Mine() {
	containerId := findProxyContainerId()

	log.Println("Attaching container log listener.")

	client, err := docker.NewClient(DockerDaemonSocket)
	if err != nil {
		panic(err)
	}

	sinceTime := time.Now()

	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	logMiner.parseAndPersistStdPipesOutput(stdoutReader, stderrReader)

	log.Println("Starting to get logs from docker daemon.")
	for {
		dockerLogErr := client.Logs(docker.LogsOptions{
			Container:         containerId,
			OutputStream:      stdoutWriter,
			ErrorStream:       stderrWriter,
			Stdout:            true,
			Stderr:            true,
			Follow:            true,
			Tail:              "all",
			Since:             sinceTime.Unix(),
			InactivityTimeout: 0,
		})

		sinceTime = time.Now()

		if dockerLogErr != nil {
			panic(dockerLogErr)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (logMiner *DockerContainerLogMiner) parseAndPersistStdPipesOutput(stdout, stderr io.Reader) {
	listenToPipe := func(input io.Reader) {
		buf := bufio.NewReader(input)

		for {
			line, _ := buf.ReadString(LogLineDelimiter)

			httpRequest, err := logMiner.LoglineParser.Parse(line)
			if err != nil {
				log.Printf("Error while parsing log line, reason: '%s', log line: '%s'", err, line)
			} else {
				logMiner.HttpRequestPersister.Persist(httpRequest)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	log.Println("Listening to stdout and stderr pipes.")

	go listenToPipe(stdout)
	go listenToPipe(stderr)
}

func findProxyContainerId() string {
	log.Println("Finding proxy container.")

	client, err := docker.NewClient(DockerDaemonSocket)
	if err != nil {
		panic(err)
	}

	proxyContainerName := config.ProxyContainerName

	filters := make(map[string][]string)
	filters["name"] = append(filters["name"], "/"+proxyContainerName)

	log.Println("Getting proxy container with name: ", proxyContainerName)
	containers, err := client.ListContainers(docker.ListContainersOptions{Filters: filters})
	if err != nil {
		panic(err)
	}

	if len(containers) <= 0 {
		panic("No running container found with specified name.")
	}

	proxyContainer := containers[0]
	proxyContainerId := proxyContainer.ID

	log.Printf("Selected container id: %s", proxyContainerId)

	return proxyContainerId
}
