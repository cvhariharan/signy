package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/engineerd/signy/pkg/intoto"
	"github.com/engineerd/signy/pkg/trust"
)

type intotoVerifyCmd struct {
	verificationDir string
	ref             string
}

func newIntotoVerifyCmd() *cobra.Command {
	i := intotoVerifyCmd{}
	cmd := &cobra.Command{
		Use:   "intoto-verify",
		Short: "execute the in-toto verification",
		RunE: func(cmd *cobra.Command, args []string) error {
			i.ref = args[0]
			return i.run()
		},
	}
	cmd.Flags().StringVarP(&i.verificationDir, "verifications", "", "", "directory to run verifications")

	return cmd
}

func (i *intotoVerifyCmd) run() error {
	target, err := trust.VerifyCNABTrust(i.ref, "", trustServer, tlscacert, trustDir)
	if err != nil {
		return err
	}

	m := &intoto.Metadata{}
	err = json.Unmarshal(*target.Custom, m)
	if err != nil {
		return err
	}

	fmt.Printf("\nWriting In-Toto metadata files into %v", i.verificationDir)
	err = intoto.WriteMetadataFiles(m, i.verificationDir)
	if err != nil {
		return err
	}

	return intoto.Verify(filepath.Join(i.verificationDir, intoto.LayoutDefaultName), i.verificationDir, filepath.Join(i.verificationDir, intoto.KeyDefaultName))
}
