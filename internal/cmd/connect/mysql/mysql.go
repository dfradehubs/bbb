package mysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/boundary"
	"bbb/internal/fancy"
	"bbb/internal/globals"

	"net/url"
)

const (
	descriptionShort = `Create a connection to a mysql target`
	descriptionLong  = `
	Create a connection to a mysql target.
	It authorizes a session if needed, and open a mysql CLI using it`
)

var (
	database  string = ""
	extraArgs string = ""
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "mysql",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	cmd.Flags().StringVar(&database, "database", database, "Database to connect for")
	cmd.Flags().StringVar(&extraArgs, "extraArgs", extraArgs, "Extra arguments for mysql connection")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	var err error
	var consoleStderr bytes.Buffer
	var consoleStdout bytes.Buffer

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		fancy.Fatalf(globals.TokenRetrievalErrorMessage)
	}

	// We need a target to connect to
	if len(args) != 1 {
		fancy.Fatalf(CommandArgsNoTargetErrorMessage)
	}

	// mysql package used to open the connection
	var mysqlCli string = "mariadb"

	// Check if mysql is present in the system, if not, just return the connection to the user
	var mysqlCliPresent bool = true

	_, err = exec.LookPath(mysqlCli)
	if err != nil {
		fancy.Printf(MySQLCliNotFoundErrorMessage)
		mysqlCliPresent = false
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 1. Ask H.Boundary for an authorized session
	// This request will provide a session ID and brokered credentials associated to the target
	// (AuthorizeSession & Connect) are performed in separated steps to check type of target before connecting
	_, err = boundary.GetTargetAuthorizedSession(storedTokenReference, args[0], &consoleStdout, &consoleStderr)
	if err != nil {
		// Brutally fail when there is no output or error to handle anything
		if len(consoleStderr.Bytes()) == 0 && len(consoleStdout.Bytes()) == 0 {
			fancy.Fatalf(AuthorizeSessionErrorMessage, err.Error(), consoleStderr.String())
		}

		// Forward stderr to stdout for later processing
		consoleStdout = consoleStderr
	}

	//
	var response boundary.AuthorizeSessionResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	// On user failures, just inform the user
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		fancy.Fatalf(AuthorizeSessionUserErrorMessage, consoleStdout.String())
	}

	if len(response.Item.Credentials) == 0 {
		fancy.Printf(TargetWithNoCredentials)
	}

	//
	targetSessionToken := response.Item.AuthorizationToken

	// Extract host of the target in Boundary for later usage.
	// Remember some proxies use this to route
	targetSessionUrl, err := url.Parse(response.Item.Endpoint)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing session URL. You have to configure a valid URL in Boundary: "+err.Error())
	}

	//
	targetSessionHost := strings.Split(targetSessionUrl.Host, ":")
	if len(targetSessionHost) != 2 {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing session Host. Session URL must have <address>:<port> format: "+err.Error())
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 2. Create a TCP connection to the target with authorized session previously created
	// User commands will be performed over this connection
	sessionFileName := targetSessionToken[:10]
	connectCommand, err := boundary.GetSessionConnection(storedTokenReference, targetSessionToken)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing 'boundary connect' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	//
	stdoutFile := globals.BbbTemporaryDir + "/" + sessionFileName + ".out"
	connectSessionStdoutRaw, err := globals.GetFileContents(stdoutFile, true)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, err.Error())
	}

	//
	var connectSessionStdout boundary.ConnectSessionStdoutT
	err = json.Unmarshal(connectSessionStdoutRaw, &connectSessionStdout)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 3. Open mysql to the socket opened by the connection
	// We use mysql or just return the connection to the user

	// mysql arguments for authentication if needed
	mysqlCliArgs := []string{"--host=" + connectSessionStdout.Address, "--port=" + strconv.Itoa(connectSessionStdout.Port), "--skip-ssl-verify-server-cert"}

	// Use password to the cli command
	if response.Item.Credentials[0].Credential.Password != "" {
		mysqlCliArgs = append(mysqlCliArgs,
			"--password="+response.Item.Credentials[0].Credential.Password)
	}

	// Use user to the cli command
	if response.Item.Credentials[0].Credential.Username != "" {
		mysqlCliArgs = append(mysqlCliArgs,
			"--user="+response.Item.Credentials[0].Credential.Username)
	}

	// Extra arguments if defined
	if extraArgs != "" {
		mysqlCliArgs = append(mysqlCliArgs,
			extraArgs)
	}

	// Database to connect for if is defined
	if database != "" {
		mysqlCliArgs = append(mysqlCliArgs,
			"-D", database)
	}

	// Set mysql command
	mysqlCommand := exec.Command(mysqlCli, mysqlCliArgs...)

	mysqlCommand.Stdin = os.Stdin
	mysqlCommand.Stdout = os.Stdout
	mysqlCommand.Stderr = &consoleStderr

	durationStringFromNow, err := globals.GetDurationStringFromNow(connectSessionStdout.Expiration)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Error getting session duration: "+err.Error())
	}

	mysqlUrl := fmt.Sprintf("%s:%d", connectSessionStdout.Address, connectSessionStdout.Port)
	fancy.Printf(ConnectionSuccessfulMessage,
		connectSessionStdout.SessionId,
		durationStringFromNow,
		mysqlUrl)

	// If mysql is not present, just return the connection to the user
	if mysqlCliPresent {
		err = mysqlCommand.Run()

		if err != nil {
			fancy.Fatalf(globals.UnexpectedErrorMessage,
				"Failed executing '"+mysqlCli+"' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
		}
	} else {
		// If mysql is not present, just return the connection to the user and wait
		// for the user to close the connection manually with Cntrl+C
		WaitSignal()
	}

	// Clean up the connection
	err = connectCommand.Process.Kill()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"\nFailed killing background connection to H.Boundary: %v\n", err)
	}
}
