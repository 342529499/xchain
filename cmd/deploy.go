package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/1851616111/util/rand"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/1851616111/xchain/pkg/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strings"
	"reflect"
)

type XChainOption struct {
	name string
	path string
	lang string

	//TODO ：支持yaml，json启动, 不然复杂的参数和函数让用户输入比较麻烦
	initJson string

	targetAddress string
}

func newCommandDeploy(out io.Writer) (*cobra.Command, *XChainOption) {
	options := new(XChainOption)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the specified chaincode to the network.",
		Long:  "Deploy the specified chaincode to the network.e",
		Run: func(c *cobra.Command, args []string) {
			if err := options.validate(); err != nil {
				fmt.Printf("deploy xcode %s error: %v\n", err)
				return
			}
			options.Run(c, args)

		},
	}

	parseXChainCmd(cmd, options)
	return cmd, options
}

func parseXChainCmd(cmd *cobra.Command, option *XChainOption) {
	flags := cmd.Flags()

	flags.StringVarP(&option.path, "path", "p", "", "code path(local,http,https)")
	flags.StringVarP(&option.name, "name", "n", "", "default name is code repository name")
	flags.StringVarP(&option.lang, "lang", "l", "golang", "xcode deploy language")
	flags.StringVarP(&option.initJson, "init", "i", "", "xcode init json msg")
	flags.StringVarP(&option.targetAddress, "target", "t", "0.0.0.0:10960", "xcode deploy target address")

}

func (o *XChainOption) validate() error {
	if len(strings.TrimSpace(o.path)) == 0 {
		return errors.New("deploy param path must not be nil.")
	}

	if len(strings.TrimSpace(o.initJson)) == 0 {
		return errors.New("deploy param init must not be nil.")
	}

	if err := json.Unmarshal([]byte(o.initJson), &map[string]interface{}{}); err != nil {
		return fmt.Errorf("deploy param init parse json format err:%v", err)
	}

	return nil
}

func (o *XChainOption) Run(c *cobra.Command, args []string) {
	if err := o.validate(); err != nil {
		log.Printf("command run err %v\n", err.Error())
	}

	getName := func() (name string) {

		if o.name != "" {
			name = o.name
		} else {
			//没有指定xcode的名字
			if os.IsPathSeparator(uint8(rune(o.lang[len(o.lang)-1]))) {
				o.lang = o.lang[:len(o.lang)-1]
			}

			tmp := strings.TrimLeftFunc(o.lang, func(c rune) bool {
				return os.IsPathSeparator(uint8(c))
			})

			name = fmt.Sprintf("%s-%s", tmp, rand.String(4))
		}

		return
	}

	getInput := func() (*pb.XCodeInput, error) {
		in := &pb.XCodeInput{}
		if err := json.Unmarshal([]byte(o.initJson), &in); err != nil {
			return nil, fmt.Errorf("Chaincode argument error: %s", err)
		}

		return in, nil
	}

	input, err := getInput()
	if err != nil {
		fmt.Printf("deploy parse param init err %v\n", input)
	}

	deploySpec := &pb.XCodeSpec{
		Type: pb.XCodeSpec_GOLANG,
		XcodeID: &pb.XCodeID{
			Path: o.path,
			Name: getName(),
		},
		XcodeMsg: input,
	}

	conn, err := grpc.Dial(o.targetAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("deploy at target(%s), dial err: %v.", o.targetAddress, err)
		return
	}

	cli := pb.NewNetClient(conn)

	serverStream, err := cli.Connect(context.Background())
	if err != nil {
		log.Printf("deploy at target(%s), connect err: %v.", o.targetAddress, err)
		return
	}

	if err = serverStream.Send(server.MakeDeployMsg(deploySpec)); err != nil {
		fmt.Printf("deploy xcode %s fail , send errMsg:%v\n", err)
		return
	}

	var result *pb.Message
	for {
		result, err = serverStream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a deploy answer: %v", err)
		}
		break
	}

	if server.IsMessageFunc(result, func(a *pb.Message) bool{
		return reflect.DeepEqual(a, result)
	}) {
		fmt.Printf("deploy xcode %s fail, errMsg:%v\n", *result)
	} else {
		fmt.Printf("deploy xcode %s ok.", deploySpec.XcodeID.Name)
	}
}
