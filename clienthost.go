package ubot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/1354092549/wsrpc"
	"github.com/gorilla/websocket"
)

type RPCConfigurator func(rpc *wsrpc.WebsocketRPC, rpcConn *wsrpc.WebsocketRPCConn) error
type ClientRegistrar func(managerUrl *url.URL, manager *Manager) (string, error)

func AssertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func DialRouter(operator string, urlStr string, registerClient ClientRegistrar) (*websocket.Conn, error) {
	switch strings.ToLower(operator) {
	case "applyto":
		managerUrl, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		if managerUrl.User != nil {
			user := managerUrl.User.Username()
			password, _ := managerUrl.User.Password()
			userData := make(url.Values)
			userData["user"] = []string{user}
			userData["password"] = []string{password}

			getTokenUrl := *managerUrl
			getTokenUrl.User = nil
			if getTokenUrl.Scheme == "wss" {
				getTokenUrl.Scheme = "https"
			} else {
				getTokenUrl.Scheme = "http"
			}
			getTokenUrl.Path = "/api/manager/get_token"
			getTokenUrlStr := getTokenUrl.String()

			fmt.Printf("Getting manager token from %s\n", getTokenUrlStr)
			res, err := http.PostForm(getTokenUrlStr, userData)
			if err != nil {
				return nil, errors.New("failed to get the manager token: " + err.Error())
			}
			if res.StatusCode != http.StatusOK {
				return nil, errors.New("failed to get the manager token: " + res.Status)
			}
			defer res.Body.Close()
			managerTokenBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, errors.New("failed to get the manager token: " + err.Error())
			}
			managerToken := string(managerTokenBytes)
			fmt.Println("Got manager token")

			managerUrl.User = nil
			managerUrl.RawQuery = "token=" + managerToken
		}
		managerUrlStr := managerUrl.String()
		fmt.Printf("Applying to %s\n", managerUrlStr)
		managerConn, _, err := websocket.DefaultDialer.Dial(managerUrlStr, nil)
		if err != nil {
			return nil, errors.New("failed to dial to manager: " + err.Error())
		}
		managerRPCConn := wsrpc.NewWebsocketRPC().Connect(managerConn)
		go func() {
			managerRPCConn.ServeConn()
			managerConn.Close()
		}()
		var manager Manager
		manager.Get(managerRPCConn)
		urlStr, err = registerClient(managerUrl, &manager)
		_ = managerConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			return nil, errors.New("failed to register the client: " + err.Error())
		}
		fmt.Println("Applied")
	case "connect":
		break
	default:
		return nil, errors.New("unknown operator: " + operator)
	}
	fmt.Printf("Connecting to %s\n", urlStr)
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		return nil, errors.New("failed to dial to router: " + err.Error())
	}
	return conn, nil
}

func HostClient(registerClient ClientRegistrar, configRPC RPCConfigurator) error {
	var err error
	fmt.Println("Press Ctrl+C to shut down.")
	if len(os.Args) < 3 {
		return errors.New("missing required arguments")
	}
	operator, urlStr := os.Args[1], os.Args[2]
	var conn *websocket.Conn
	for retryCount := 0; retryCount < 5; retryCount++ {
		conn, err = DialRouter(operator, urlStr, registerClient)
		if err == nil {
			break
		}
		fmt.Println("Failed to connect to UBot Router, it will try again in 5 seconds.")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("Failed to connect to UBot Router after 5 attempts: %v", err)
	}
	fmt.Println("Connected")
	rpc := wsrpc.NewWebsocketRPC()
	rpcConn := rpc.Connect(conn)
	err = configRPC(rpc, rpcConn)
	if err != nil {
		return err
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	fn := make(chan int, 1)
	go func() {
		rpcConn.ServeConn()
		conn.Close()
		fn <- 0
	}()
	select {
	case <-sc:
		fmt.Println("Shutting down...")
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		<-fn
	case <-fn:
		fmt.Println("Session ended")
	}
	return nil
}

func HostAccount(id string, creater func(*AccountEventEmitter) *Account) error {
	return HostClient(func(managerUrl *url.URL, manager *Manager) (string, error) {
		token, err := manager.RegisterAccount(id)
		if err != nil {
			return "", err
		}
		urlObj := *managerUrl
		urlObj.Path = "/api/account"
		query := url.Values{}
		query.Set("id", id)
		query.Set("token", token)
		urlObj.RawQuery = query.Encode()
		return urlObj.String(), nil
	}, func(rpc *wsrpc.WebsocketRPC, rpcConn *wsrpc.WebsocketRPCConn) error {
		remoteObj := new(AccountEventEmitter)
		remoteObj.Get(rpcConn)
		localObj := creater(remoteObj)
		localObj.Register(rpc)
		return nil
	})
}

func HostApp(id string, creater func(*AppApi) *App) error {
	return HostClient(func(managerUrl *url.URL, manager *Manager) (string, error) {
		token, err := manager.RegisterApp(id)
		if err != nil {
			return "", err
		}
		urlObj := *managerUrl
		urlObj.Path = "/api/app"
		query := url.Values{}
		query.Set("id", id)
		query.Set("token", token)
		urlObj.RawQuery = query.Encode()
		return urlObj.String(), nil
	}, func(rpc *wsrpc.WebsocketRPC, rpcConn *wsrpc.WebsocketRPCConn) error {
		remoteObj := new(AppApi)
		remoteObj.Get(rpcConn)
		localObj := creater(remoteObj)
		localObj.Register(rpc)
		return nil
	})
}
