package installconfig

type BootstrapInPlaceInstallConfigSpec struct {
	BaseDomain        string
	ClusterName       string
	ClusterNetwork    string
	MachineNetwork    string
	ClusterSvcNetwork string
	InstallDisk       string
	PullSecret        string
	SshPubKey         string
}
