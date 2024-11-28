package types

// Config represents the structure of the configuration file (config.json).
// It includes the settings required to configure the application, such as
// a proxy for HTTP requests and an API key for game authentication.
//
// # Fields:
//   - Proxy: The residential proxy settings, which include IP, port, and authentication details.
//   - APIKey: The API key used to refresh game authentication.
//
// # Example config.json:
//
//	{
//		"proxy": {
//			"ip": "192.168.1.100",
//			"port": 8080,
//			"username": "proxyuser",
//			"password": "proxypass",
//			"socksType": 5,
//			"timeout": 30
//		},
//		"api_key": "your-api-key-here"
//	}
//
// # Example Usage:
//
//	config := Config{
//		Proxy: Proxy{
//			Ip:        "192.168.1.100",
//			Port:      8080,
//			Username:  "proxyuser",
//			Password:  "proxypass",
//			SocksType: 5,
//			Timeout:   30,
//		},
//		APIKey: "your-api-key-here",
//	}
//	fmt.Println(config.Proxy.Ip) // Output: 192.168.1.100
type Config struct {
	Proxy  Proxy  `json:"proxy"`   // Proxy contains the details of the HTTP/SOCKS proxy configuration.
	APIKey string `json:"api_key"` // APIKey is the key for authenticating API requests.
}

// Proxy represents the settings for configuring an SOCKS proxy server.
// It includes the proxy server's IP address, port, and optional authentication credentials.
//
// # Fields:
//   - Ip: The IP address of the proxy server (e.g., "192.168.1.100").
//   - Port: The port number for connecting to the proxy server (e.g., 8080).
//   - Username: The username for proxy authentication (if required).
//   - Password: The password for proxy authentication (if required).
//   - SocksType: The SOCKS protocol type (e.g., 4 or 5).
//   - Timeout: The timeout in seconds for proxy connections.
//
// # Example Usage:
//
//	proxy := Proxy{
//		Ip:        "192.168.1.100",
//		Port:      8080,
//		Username:  "proxyuser",
//		Password:  "proxypass",
//		SocksType: 5,
//		Timeout:   30,
//	}
//	fmt.Printf("Proxy: %s:%d\n", proxy.Ip, proxy.Port) // Output: Proxy: 192.168.1.100:8080
type Proxy struct {
	Ip        string `json:"ip"`        // Ip is the proxy server's IP address.
	Port      int    `json:"port"`      // Port is the proxy server's port.
	Username  string `json:"username"`  // Username is the username for proxy authentication.
	Password  string `json:"password"`  // Password is the password for proxy authentication.
	SocksType int    `json:"socksType"` // SocksType specifies the SOCKS protocol type (4 or 5).
	Timeout   int    `json:"timeout"`   // Timeout specifies the timeout in seconds for proxy connections.
}

// Account represents an entry in the accounts file (accounts.json).
// Each account includes associated game-specific data and Telegram settings for authentication refreshing.
//
// # Fields:
//   - GameData: The game-specific data associated with this account.
//   - TelegramData: The Telegram session information, including credentials and IDs.
//
// # Example accounts.json:
//
//	[
//		{
//			"game-data": "user=%7B%22id%22%3A78894796...",
//			"telegram": {
//				"tdataStringSession": "session-data-string",
//				"appId": "123456",
//				"appHash": "abcdef123456",
//				"telegramId": "987654321"
//			}
//		},
//		{
//			"game-data": "user=%7B%22id%22%3A78894796...",
//			"telegram": {
//				"tdataStringSession": "another-session-data",
//				"appId": "654321",
//				"appHash": "fedcba654321",
//				"telegramId": "123456789"
//			}
//		}
//	]
//
// # Example Usage:
//
//	account := Account{
//		GameData: "user=%7B%22id%22%3A78894796...",
//		TelegramData: TelegramData{
//			TdataStringSession: "session-data-string",
//			AppId:              "123456",
//			AppHash:            "abcdef123456",
//			TelegramId:         "987654321",
//		},
//	}
//	fmt.Println(account.GameData)             // Output: user=%7B%22id%22%3A78894796...
//	fmt.Println(account.TelegramId)           // Output: 987654321
type Account struct {
	GameData     string            `json:"game-data"` // Game-specific data associated with this account.
	TelegramData `json:"telegram"` // Telegram session information.
}

// TelegramData represents the Telegram session information, including credentials and IDs.
//
// # Fields:
//   - TdataStringSession: The session data string for Telegram.
//   - AppId: The application ID for Telegram.
//   - AppHash: The application hash for Telegram.
//   - TelegramId: The Telegram ID.
//
// # Example Usage:
//
//	telegramData := TelegramData{
//		TdataStringSession: "session-data-string",
//		AppId:              "123456",
//		AppHash:            "abcdef123456",
//		TelegramId:         "987654321",
//	}
//	fmt.Println(telegramData.TelegramId) // Output: 987654321
type TelegramData struct {
	TdataStringSession string `json:"tdataStringSession"`
	AppId              string `json:"appId"`
	AppHash            string `json:"appHash"`
	TelegramId         string `json:"telegramId"`
}
