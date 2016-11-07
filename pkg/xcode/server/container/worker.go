package container

import (
	"bytes"
	"github.com/1851616111/go-dockerclient"
	"sync"

	pb "github.com/1851616111/xchain/pkg/protos"
)

var (
	clientOnce *sync.Once = new(sync.Once)

	//TODO：每次有操作时都创建一个client？？
	client     *docker.Client
	refreshErr error

	default_StopContainer_Timeout_Sec = 60
)

func getDefaultClient() (*docker.Client, error) {
	if client == nil {
		clientOnce.Do(func() {
			client, refreshErr = docker.NewClient("unix:///var/run/docker.sock")
		})
	}

	return client, refreshErr
}

type container struct {
	client *docker.Client
}

//ping docker 失败后刷新client
func (c *container) RefreshClient() error {
	client, refreshErr = getDefaultClient()
	if refreshErr != nil {
		logger.Printf("refresh container(docker) client failed. err: %v\n", refreshErr)
		return refreshErr
	}

	return nil
}

func (w *Worker) listImage() ([]docker.APIImages, error) {
	workOpts, ok := w.opts.(*docker.ListImagesOptions)
	if !ok {
		return nil, InterfaceAssertError("*go-dockerclient.ListImagesOptions")
	}

	_, err := client.ListImages(*workOpts)
	if err != nil {
		logger.Printf("Error list images err: %s", err)
		return nil, err
	}

	return client.ListImages(*workOpts)
}

func (w *Worker) buildImage() error {

	var workSepc *pb.XCodeSpec
	var workOpts *docker.BuildImageOptions
	var ok bool

	if workSepc, ok = w.metadata.(*pb.XCodeSpec); !ok {
		return InterfaceAssertError("*xcode.WorkSpec")
	}

	if workOpts, ok = w.opts.(*docker.BuildImageOptions); !ok {
		return InterfaceAssertError("*go-dockerclient.BuildImageOptions")
	}

	reader, err := GetXCodePackageBytes(workSepc)
	if err != nil {
		return err
	}

	output := bytes.NewBuffer(nil)

	workOpts.OutputStream = output
	workOpts.InputStream = reader

	if err := client.BuildImage(*workOpts); err != nil {
		logger.Printf("Error building images: %s", err)
		logger.Printf("Image Output:\n********************\n%s\n********************", output.String())
		return err
	}
	return nil
}

func (w *Worker) createContainer() error {
	workOpts, ok := w.opts.(*docker.CreateContainerOptions)
	if !ok {
		return InterfaceAssertError("*go-dockerclient.CreateContainerOptions")
	}

	_, err := client.CreateContainer(*workOpts)
	if err != nil {
		logger.Printf("Error create container: %s", err)
		return err
	}

	logger.Printf("Created container: %s", workOpts.Config.Image)

	client.StartContainer(workOpts.Name, nil)
	return nil
}

func (w *Worker) removeContainer() error {
	id, ok := w.metadata.(string)
	if ok {
		return InterfaceAssertError("string")
	}

	err := client.StopContainer(id, uint(default_StopContainer_Timeout_Sec))
	if err != nil {
		logger.Printf("Stop container %s(%s)", id, err)
	} else {
		logger.Printf("Stopped container %s", id)
	}
	err = client.KillContainer(docker.KillContainerOptions{ID: id})
	if err != nil {
		logger.Printf("Kill container %s (%s)", id, err)
	} else {
		logger.Printf("Killed container %s", id)
	}
	err = client.RemoveContainer(docker.RemoveContainerOptions{ID: id, Force: true})
	if err != nil {
		logger.Printf("Remove container %s (%s)", id, err)
	} else {
		logger.Printf("Removed container %s", id)
	}

	return err
}

func (w *Worker) removeImage() error {
	var id string
	var workOpts *docker.RemoveImageOptions
	var ok bool

	if id, ok = w.metadata.(string); !ok {
		return InterfaceAssertError("string")
	}

	if workOpts, ok = w.opts.(*docker.RemoveImageOptions); !ok {
		return InterfaceAssertError("*go-dockerclient.RemoveImageOptions")
	}

	err := client.RemoveImageExtended(id, *workOpts)
	if err != nil {
		logger.Printf("error while destroying image: %s", err)
	} else {
		logger.Printf("Destroyed image %s", id)
	}

	return nil
}
