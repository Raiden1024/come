package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

// Config json
type Config struct {
	FireToken string `json:"fireToken"`
	Host      string
}

// OpenfireSession ... Retrieve Session datas
type OpenfireSession struct {
	Session OpenfireDatas `json:"session"`
}

// OpenfireSession ... Retrieve All Sessions Datas
type OpenfireSessions struct {
	Sessions []OpenfireDatas `json:"session"`
}

// OpenfireDatas ... Retrieve Sessions Fields
type OpenfireDatas struct {
	SessionID      string `json:"sessionId"`
	Username       string `json:"username"`
	Ressource      string `json:"ressource"`
	Node           string `json:"node"`
	SessionStatus  string `json:"sessionStatus"`
	PresenceStatus string `json:"presenceStatus"`
	Priority       string `json:"priority"`
	HostAddress    string `json:"hostAddress"`
	HostName       string `json:"hostName"`
	CreationDate   string `json:"creationDate"`
	LastActionDate string `json:"lastActionDate"`
	Secure         bool   `json:"secure"`
}

var (
	httpClient = &http.Client{}
	arguments  = os.Args
	config     = Config{}
	usageTxt   = "Usage: come [ARGUMENT] [USER]\n" +
		"-c ou c    SSH connection to a user machine, ex: come -c <user>\n" +
		"-i ou i    Display User IP Address, ex: come -i <user>\n" +
		"-w ou w    Wait for user Online status, ex: come -w <user>\n" +
		"-l ou l    Display active sessions list\n" +
		"-v ou v    Print Version\n" +
		"-h ou h    This help"
)

func init() {
	if _, err := os.Stat(os.ExpandEnv("$HOME") + "/.config/comecfg.json"); err != nil {
		fmt.Fprint(os.Stdout, "No configuration found."+
			"\nTo use 'come' you must have a functional OpenFire Server\n"+
			"With REST API plugin installed.\n")
		fmt.Print("openfire server address (format needed: 'http(s)://host(:port)': ")
		entry := bufio.NewScanner(os.Stdin)
		entry.Scan()
		hostname := entry.Text()
		fmt.Print("openfire server API Secret Key: ")
		entry.Scan()
		token := entry.Text()
		f, err := os.OpenFile(os.ExpandEnv("$HOME")+"/.config/comecfg.json", os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not write configuration: %v\n", err)
		}
		config = Config{
			FireToken: token,
			Host:      hostname,
		}
		newConfig, err := json.Marshal(&config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not write configuration: %v\n", err)
		}
		f.Write(newConfig)
		f.Close()
		fmt.Println("Configuration created in " + os.ExpandEnv("$HOME") + "/.config/comecfg.json")
	}
	f, err := os.OpenFile(os.ExpandEnv("$HOME")+"/.config/comecfg.json", os.O_RDONLY, 0664)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read configuration: %v", err)
	}
	reader, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read configuration: %v", err)
	}
	f.Close()
	if err := json.Unmarshal(reader, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Could not read json file: %v", err)
	}
}

func listSessions() error {
	openfireSessions := OpenfireSessions{}
	req, err := http.NewRequest("GET", config.Host+"/plugins/restapi/v1/sessions/", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", config.FireToken)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &openfireSessions)
	for _, v := range openfireSessions.Sessions {
		fmt.Println(v.Username, v.HostAddress)
	}
	return nil
}

func sshConnect(userid string) (string, error) {
	openfireSession := OpenfireSession{}
	req, err := http.NewRequest("GET", config.Host+"/plugins/restapi/v1/sessions/"+userid, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", config.FireToken)
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &openfireSession)
	return openfireSession.Session.HostAddress, nil
}

func main() {
	if len(arguments) < 2 {
		fmt.Println(usageTxt)
	} else if len(arguments) > 3 {
		fmt.Println("Too much parameters")
	} else if len(arguments) == 2 {
		if arguments[1] == "-h" || arguments[1] == "h" {
			fmt.Println(usageTxt)
		} else if arguments[1] == "-v" || arguments[1] == "v" {
			fmt.Println("COME (COnnect ME) - version: 1.1")
		} else if arguments[1] == "-l" || arguments[1] == "l" {
			err := listSessions()
			if err != nil {
				fmt.Printf("%v", err)
			}
		} else {
			fmt.Println("Username missing")
		}
	} else if len(arguments) == 3 {
		if arguments[1] == "-c" || arguments[1] == "c" {
			ip, err := sshConnect(arguments[2])
			if err != nil {
				fmt.Printf("Unable to establish connection: %v", err)
			}
			fmt.Println("Waiting for connection...")
			cmd := exec.Command("ssh", "root@"+ip)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.Run()

		} else if arguments[1] == "-w" || arguments[1] == "w" {
			fmt.Println("Waiting for User Online Status...")
			for {
				ip, err := sshConnect(arguments[2])
				if err != nil {
					fmt.Printf("%v", err)
				}
				if ip != "" {
					break
				}
			}
			fmt.Println(arguments[2] + ": ONLINE")
		} else if arguments[1] == "-i" || arguments[1] == "i" {
			ip, err := sshConnect(arguments[2])
			if err != nil {
				fmt.Printf("Unable to obtain IP address: %v", err)
			}
			fmt.Printf("%s\n", ip)
		} else {
			fmt.Println("Unknown parameter, -h for help...")
		}
	}
}
