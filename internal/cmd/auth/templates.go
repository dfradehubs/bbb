package auth

const (

	// AuthErrorMessage message is thrown when there is an error
	// different from 4xx on 'auth' command execution
	AuthErrorMessage = `
	{Red}Error executing 'auth' command.

	{White}Command execution returned:
	{Italic}%s{Reset}

	{White}H.Boundary command under the hood returned:
	{Italic}%s`

	// AuthUserErrorMessage message is thrown when there is a 4xx error on auth command execution
	AuthUserErrorMessage = `
	{Red}There is something wrong when authenticating the user

	{Magenta}Review following points:
	* You have enough permissions to connect to H.Boundary

	{White}Under the hood:
	{Italic}%s`

	// AuthSuccessfulMessage represents the message thrown when everything finished as expected on 'auth' command
	AuthSuccessfulMessage = `
	{White}You are {Bold}{Cyan}successfully{Reset} authenticated in H.Boundary {Reset}


	{White}You are ready to list projects or their targets as follows: {Reset}
	{Bold}{White}console ~$ {Green}bbb list {Reset}
	{Bold}{White}console ~$ {Green}bbb list {Cyan}{abbreviation}{Reset}`
)
