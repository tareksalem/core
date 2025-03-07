package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	corecluster "github.com/projecteru2/core/cluster"
	"github.com/projecteru2/core/engine"
	enginetypes "github.com/projecteru2/core/engine/types"
	"github.com/projecteru2/core/log"
	"github.com/projecteru2/core/types"
	coretypes "github.com/projecteru2/core/types"
	"github.com/projecteru2/core/utils"

	"github.com/docker/distribution/reference"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/blkiodev"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerapi "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/registry"
	"github.com/docker/go-units"
)

type fuckDockerStream struct {
	conn net.Conn
	buf  io.Reader
}

func (f fuckDockerStream) Read(p []byte) (n int, err error) {
	return f.buf.Read(p)
}

func (f fuckDockerStream) Close() error {
	return f.conn.Close()
}

func mergeStream(stream io.ReadCloser) io.Reader {
	outr, outw := io.Pipe()

	go func() {
		defer stream.Close()
		_, err := stdcopy.StdCopy(outw, outw, stream)
		_ = outw.CloseWithError(err)
	}()

	return outr
}

// FuckDockerStream will copy docker stream to stdout and err
func FuckDockerStream(stream dockertypes.HijackedResponse) io.ReadCloser {
	outr := mergeStream(ioutil.NopCloser(stream.Reader))
	return fuckDockerStream{stream.Conn, outr}
}

// make mount paths
// 使用volumes, 参数格式跟docker一样
// volumes:
//     - "/foo-data:$SOMEENV/foodata:rw"
func makeMountPaths(opts *enginetypes.VirtualizationCreateOptions) ([]string, map[string]struct{}) {
	binds := []string{}
	volumes := make(map[string]struct{})

	var expandENV = func(env string) string {
		envMap := map[string]string{}
		for _, env := range opts.Env {
			parts := strings.Split(env, "=")
			envMap[parts[0]] = parts[1]
		}
		return envMap[env]
	}

	for _, path := range opts.Volumes {
		expanded := os.Expand(path, expandENV)
		parts := strings.Split(expanded, ":")
		if len(parts) == 2 {
			binds = append(binds, fmt.Sprintf("%s:%s:rw", parts[0], parts[1]))
			volumes[parts[1]] = struct{}{}
		} else if len(parts) >= 3 {
			binds = append(binds, fmt.Sprintf("%s:%s:%s", parts[0], parts[1], parts[2]))
			volumes[parts[1]] = struct{}{}
			if len(parts) == 4 {
				log.Warn("[makeMountPaths] docker engine not support volume with size limit")
			}
		}
	}

	return binds, volumes
}

func makeResourceSetting(cpu float64, memory int64, cpuMap map[string]int64, numaNode string, iopsOptions map[string]string, remap bool) dockercontainer.Resources {
	resource := dockercontainer.Resources{}

	resource.CPUQuota = 0
	resource.CPUShares = defaultCPUShare
	resource.CPUPeriod = corecluster.CPUPeriodBase
	if cpu > 0 {
		resource.CPUQuota = int64(cpu * float64(corecluster.CPUPeriodBase))
	} else if cpu == -1 {
		resource.CPUQuota = -1
	}

	if len(cpuMap) > 0 {
		cpuIDs := []string{}
		for cpuID := range cpuMap {
			cpuIDs = append(cpuIDs, cpuID)
		}
		resource.CpusetCpus = strings.Join(cpuIDs, ",")
		// numaNode will empty or numaNode
		resource.CpusetMems = numaNode

		if remap {
			resource.CPUShares = int64(1024)
		} else {
			// unrestrained cpu quota for binding
			resource.CPUQuota = -1
			// cpu share for fragile pieces
			if _, divpart := math.Modf(cpu); divpart > 0 {
				resource.CPUShares = int64(math.Round(float64(1024) * divpart))
			}
		}
	}
	resource.Memory = memory
	resource.MemorySwap = memory
	resource.MemoryReservation = memory / 2
	if memory != 0 && memory/2 < int64(units.MiB*4) {
		resource.MemoryReservation = int64(units.MiB * 4)
	}

	if len(iopsOptions) > 0 {
		var readIOPSDevices, writeIOPSDevices, readBPSDevices, writeBPSDevices []*blkiodev.ThrottleDevice
		for device, options := range iopsOptions {
			parts := strings.Split(options, ":")
			for len(parts) < 4 {
				parts = append(parts, "0")
			}
			var readIOPS, writeIOPS, readBPS, writeBPS int64
			readIOPS, _ = utils.ParseRAMInHuman(parts[0])
			writeIOPS, _ = utils.ParseRAMInHuman(parts[1])
			readBPS, _ = utils.ParseRAMInHuman(parts[2])
			writeBPS, _ = utils.ParseRAMInHuman(parts[3])

			readIOPSDevices = append(readIOPSDevices, &blkiodev.ThrottleDevice{
				Path: device,
				Rate: uint64(readIOPS),
			})
			writeIOPSDevices = append(writeIOPSDevices, &blkiodev.ThrottleDevice{
				Path: device,
				Rate: uint64(writeIOPS),
			})
			readBPSDevices = append(readBPSDevices, &blkiodev.ThrottleDevice{
				Path: device,
				Rate: uint64(readBPS),
			})
			writeBPSDevices = append(writeBPSDevices, &blkiodev.ThrottleDevice{
				Path: device,
				Rate: uint64(writeBPS),
			})
		}
		resource.BlkioDeviceReadIOps = readIOPSDevices
		resource.BlkioDeviceWriteIOps = writeIOPSDevices
		resource.BlkioDeviceReadBps = readBPSDevices
		resource.BlkioDeviceWriteBps = writeBPSDevices
	}

	return resource
}

// 只要一个image的前面, tag不要
func normalizeImage(image string) string {
	if strings.Contains(image, ":") {
		t := strings.Split(image, ":")
		return t[0]
	}
	return image
}

// image begin
// MakeAuthConfigFromRemote Calculate encoded AuthConfig from registry and eru-core config
// See https://github.com/docker/cli/blob/16cccc30f95c8163f0749eba5a2e80b807041342/cli/command/registry.go#L67
func makeEncodedAuthConfigFromRemote(authConfigs map[string]coretypes.AuthConfig, remote string) (string, error) {
	ref, err := reference.ParseNormalizedNamed(remote)
	if err != nil {
		return "", err
	}

	// Resolve the Repository name from fqn to RepositoryInfo
	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return "", err
	}

	serverAddress := repoInfo.Index.Name
	if authConfig, exists := authConfigs[serverAddress]; exists {
		if encodedAuth, err := encodeAuthToBase64(authConfig); err == nil {
			return encodedAuth, nil
		}
		return "", err
	}
	return "dummy", nil
}

// EncodeAuthToBase64 serializes the auth configuration as JSON base64 payload
// See https://github.com/docker/cli/blob/master/cli/command/registry.go#L41
func encodeAuthToBase64(authConfig coretypes.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

// Image tag
// 格式严格按照 Hub/HubPrefix/appname:tag 来
func createImageTag(config types.DockerConfig, appname, tag string) string {
	prefix := strings.Trim(config.Namespace, "/")
	if prefix == "" {
		return fmt.Sprintf("%s/%s:%s", config.Hub, appname, tag)
	}
	return fmt.Sprintf("%s/%s/%s:%s", config.Hub, prefix, appname, tag)
}

func makeCommonPart(build *enginetypes.Build) (string, error) {
	tmpl := template.Must(template.New("common").Parse(commonTmpl))
	out := bytes.Buffer{}
	if err := tmpl.Execute(&out, build); err != nil {
		return "", err
	}
	return out.String(), nil
}

func makeUserPart(opts *enginetypes.BuildContentOptions) (string, error) {
	tmpl := template.Must(template.New("user").Parse(userTmpl))
	out := bytes.Buffer{}
	if err := tmpl.Execute(&out, opts); err != nil {
		return "", err
	}
	return out.String(), nil
}

func makeMainPart(_ *enginetypes.BuildContentOptions, build *enginetypes.Build, from string, commands, copys []string) (string, error) {
	var buildTmpl []string
	common, err := makeCommonPart(build)
	if err != nil {
		return "", err
	}
	buildTmpl = append(buildTmpl, from, common)
	if len(copys) > 0 {
		buildTmpl = append(buildTmpl, copys...)
	}
	if len(commands) > 0 {
		buildTmpl = append(buildTmpl, commands...)
	}
	return strings.Join(buildTmpl, "\n"), nil
}

// Dockerfile
func createDockerfile(dockerfile, buildDir string) error {
	f, err := os.Create(filepath.Join(buildDir, "Dockerfile"))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(dockerfile)
	return err
}

// CreateTarStream create a tar stream
func CreateTarStream(path string) (io.ReadCloser, error) {
	tarOpts := &archive.TarOptions{
		ExcludePatterns: []string{},
		IncludeFiles:    []string{"."},
		Compression:     archive.Uncompressed,
		NoLchown:        true,
	}
	return archive.TarWithOptions(path, tarOpts)
}

// GetIP Get hostIP
func GetIP(ctx context.Context, daemonHost string) string {
	u, err := url.Parse(daemonHost)
	if err != nil {
		log.Errorf(ctx, "[GetIP] GetIP %s failed %v", daemonHost, err)
		return ""
	}
	return u.Hostname()
}

func makeDockerClient(_ context.Context, config coretypes.Config, client *http.Client, endpoint string) (engine.API, error) {
	cli, err := dockerapi.NewClientWithOpts(
		dockerapi.WithHost(endpoint),
		dockerapi.WithVersion(config.Docker.APIVersion),
		dockerapi.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return &Engine{cli, config}, nil
}

func useCNI(labels map[string]string) bool {
	for k, v := range labels {
		if k == "cni" && v == "1" {
			return true
		}
	}
	return false
}
