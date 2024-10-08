package kube

const (
	KubectlCliNotFoundErrorMessage = `
	{Red}Kubectl CLI is not detected on your system.

	{Magenta}Well, this is a bit... unexpected. Unexpected because if you are trying to connect to Kubernetes,
	you are probably an expert, but that 'expert' {Cyan}COMPLETELY forgot to install kubectl!!{Magenta} 
	So you may be a {Cyan}noob{Magenta}, or just a crazy guy. 

	In both cases I can do something for you.

	How to pronounce 'kubectl': https://www.youtube.com/watch?v=2wgAIvXpJqU
	Some courses for you: https://www.youtube.com/achetronic

	{White}It is possible to remediate this issue installing it from:

	{White}Pretty way: {Reset}{Cyan}Use your OS package manager
	{White}Experts website: {Reset}{Cyan}https://kubernetes.io/docs/tasks/tools/#kubectl
	`

	// CommandArgsNoTargetErrorMessage is the message thrown when the user is trying to establish a connection
	// without defining a target
	CommandArgsNoTargetErrorMessage = `
	{Red}Impossible to get target ID from arguments.

	{White}To connect to a target, you need to connect using {Cyan}Target ID {White}as follows:
	{Bold}{White}console ~$ {Green}bbb connect kube {Cyan}{ttcp_example}`

	// AuthorizeSessionErrorMessage message is thrown when there is an error different from 4xx on authorize-session
	// command execution
	AuthorizeSessionErrorMessage = `
	{Red}Error executing 'authorize-session' command.

	{White}Command execution returned:
	{Italic}%s{Reset}

	{White}H.Boundary command under the hood returned:
	{Italic}%s`

	// AuthorizeSessionUserErrorMessage message is thrown when there is a 4xx error on authorize-session command execution
	AuthorizeSessionUserErrorMessage = `
	{Red}There is something wrong with your target when authorizing the session

	{Magenta}Review following points:
	* Your H.Boundary token is still valid
	* You have properly written the target id
	* You have enough permissions to authorize a session against to that target

	{White}Under the hood:
	{Italic}%s`

	// NotKubeTargetErrorMessage message is thrown when user is trying to use 'kube' connection type on non-kube targets
	NotKubeTargetErrorMessage = `
	{Red}Selected target is not configured as Kubernetes target.
	
	{Magenta}Contact your H.Boundary administrators and may the force be with you`

	// ConnectionSuccessfulMessage represents the message thrown when everything finished as expected
	// and shows how to use recently created (connection + kubeconfig) to the user
	ConnectionSuccessfulMessage = `
	{Magenta}Some reminders: {Reset}

	* Your H.Boundary session ID is: {Bold}{Cyan}%s {Reset}
	* Your session will expire in: {Bold}{Cyan}%s {Reset}


	Execute the following command to kill the connection once you don't need it: 	
	{Bold}{White}console ~$ {Green}%s {Reset}

	Execute the following command to kill ALL the connections at once:
	{Bold}{White}console ~$ {Green}%s {Reset}


	{Magenta}You are ready to query your Kubernetes cluster using command as follows: {Reset}
	{Bold}{White}console ~$ {Green}%s`
)
