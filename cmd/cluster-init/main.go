package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	onboardcmd "github.com/openshift/ci-tools/cmd/cluster-init/cmd/onboard"
	"github.com/openshift/ci-tools/cmd/cluster-init/cmd/provision"
)

func main() {
	log := logrus.NewEntry(logrus.StandardLogger())
	ctx := handleSignals(signals.SetupSignalHandler(), log)

	// TODO: onboard is treated like the root command. Create a real
	// root command and attach onboard to it.
	onboardCmd := onboardcmd.New()

	provisionCmd, err := provision.NewProvision(ctx, log)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	onboardCmd.AddCommand(provisionCmd)

	if err := onboardCmd.Execute(); err != nil {
		logrus.Fatalf("%s", err)
	}
}

func handleSignals(signalCtx context.Context, log *logrus.Entry) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-signalCtx.Done()
		log.Warn("Received interrupt signal")
		cancel()
	}()

	return ctx
}
