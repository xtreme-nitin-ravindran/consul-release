package helpers

import (
	"github.com/pivotal-cf-experimental/bosh-test/bosh"
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/ops"
)

func DeploymentVMs(boshClient bosh.Client, deploymentName string) ([]bosh.VM, error) {
	vms, err := boshClient.DeploymentVMs(deploymentName)
	if err != nil {
		return nil, err
	}

	for index := range vms {
		vms[index].IPs = nil
	}

	return vms, nil
}

func GetVMsFromManifestV1(manifest consul.Manifest) []bosh.VM {
	var vms []bosh.VM

	for _, job := range manifest.Jobs {
		for i := 0; i < job.Instances; i++ {
			vms = append(vms, bosh.VM{JobName: job.Name, Index: i, State: "running"})
		}
	}

	return vms
}

func GetVMsFromManifest(manifest consul.ManifestV2) []bosh.VM {
	var vms []bosh.VM

	for _, ig := range manifest.InstanceGroups {
		for i := 0; i < ig.Instances; i++ {
			vms = append(vms, bosh.VM{JobName: ig.Name, Index: i, State: "running"})
		}
	}

	return vms
}

func GetVMsFromManifestV2(manifest string) []bosh.VM {
	var vms []bosh.VM

	instanceGroups, err := ops.InstanceGroups(manifest)
	if err != nil {
		panic(err)
	}

	for _, ig := range instanceGroups {
		for i := 0; i < ig.Instances; i++ {
			vms = append(vms, bosh.VM{JobName: ig.Name, Index: i, State: "running"})
		}
	}

	return vms
}

func GetVMIPs(boshClient bosh.Client, deploymentName, jobName string) ([]string, error) {
	vms, err := boshClient.DeploymentVMs(deploymentName)
	if err != nil {
		return []string{}, err
	}

	for _, vm := range vms {
		if vm.JobName == jobName {
			return vm.IPs, nil
		}
	}

	return []string{}, nil
}
