package cli

import (
	"cli/internal/models"
	service "cli/internal/services"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/skip2/go-qrcode"
)

type CLI struct {
	authService  *service.AuthService
	currentToken string
	currentUser  *models.User
}

func NewCLI(auth *service.AuthService) *CLI {
	return &CLI{
		authService: auth,
	}
}

func (c *CLI) updatePrompt(rl *readline.Instance) {
	if c.currentUser != nil {
		rl.SetPrompt(fmt.Sprintf("auth(%s)> ", c.currentUser.Username))
	} else {
		rl.SetPrompt("auth> ")
	}
}

func (c *CLI) Run() {

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "auth> ",
		HistoryFile: ".history",
		AutoComplete: readline.NewPrefixCompleter(
			readline.PcItem("register"),
			readline.PcItem("login"),
			readline.PcItem("logout"),
			readline.PcItem("whoami"),
			readline.PcItem("enable-2fa"),
			readline.PcItem("disable-2fa"),
			readline.PcItem("help"),
			readline.PcItem("exit"),
		),
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {

		c.updatePrompt(rl)

		command, err := rl.Readline()
		if err != nil {
			break
		}

		command = strings.TrimSpace(command)

		switch command {

		case "register":
			c.Register(rl)

		case "login":
			c.Login(rl)

		case "logout":
			c.Logout(rl)

		case "whoami":
			c.WhoAmI()

		case "help":
			c.Help()

		case "enable-2fa":
			c.Enable2FA(rl)

		case "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}

func (c *CLI) Register(rl *readline.Instance) {

	rl.SetPrompt("Username: ")
	username, _ := rl.Readline()

	rl.SetPrompt("Password: ")
	password, _ := rl.Readline()

	_, err := c.authService.Register(
		strings.TrimSpace(username),
		strings.TrimSpace(password),
	)

	if err != nil {
		fmt.Println(err)
		c.updatePrompt(rl)
		return
	}

	fmt.Println("Registration Successful")
	c.updatePrompt(rl)
}

func (c *CLI) Login(rl *readline.Instance) {

	rl.SetPrompt("Username: ")
	username, _ := rl.Readline()

	rl.SetPrompt("Password: ")
	password, _ := rl.Readline()

	otp := ""

	enabled, err := c.authService.IsMFAEnabled(strings.TrimSpace(username))
	if err != nil {
		fmt.Println(err)
		return
	}

	if enabled {
		rl.SetPrompt("OTP: ")
		otp, _ = rl.Readline()
	}

	user, session, err := c.authService.Login(
		strings.TrimSpace(username),
		strings.TrimSpace(password),
		strings.TrimSpace(otp),
	)

	if err != nil {
		fmt.Println(err)
		c.updatePrompt(rl)
		return
	}

	c.currentUser = user
	c.currentToken = session.Token

	fmt.Println("Login Successful")
	c.updatePrompt(rl)
}

func (c *CLI) Logout(rl *readline.Instance) {

	if c.currentToken == "" {
		fmt.Println("Please login first")
		return
	}

	err := c.authService.Logout(c.currentToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.currentToken = ""
	c.currentUser = nil

	fmt.Println("Logout Successful")
	c.updatePrompt(rl)
}

func (c *CLI) WhoAmI() {

	if c.currentToken == "" {
		fmt.Println("Please login first")
		return
	}

	user, session, err := c.authService.WhoAmI(c.currentToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nCurrent User")
	fmt.Println("------------")
	fmt.Println("Username :", user.Username)
	fmt.Println("MFA status :", user.MFAEnabled)
	fmt.Println("Registation Date :", user.CreatedAt.Format("02 Jan 2006 15:04:05"))
	fmt.Println("Session Token expiry :", session.ExpiresAt.Format("02 Jan 2006 15:04:05"))
	// fmt.Println("Session Expires :", user.SessionExpires.Format("02 Jan 2006 15:04:05"))

	if user.LastLogin != nil {
		fmt.Println("Last Login :", user.LastLogin.Format("02 Jan 2006 15:04:05"))
	}
}

func (c *CLI) Help() {
	if c.currentUser == nil {
		fmt.Println("\nAvailable Commands")
		fmt.Println("------------------")
		fmt.Println("register")
		fmt.Println("login")
		fmt.Println("help")
		fmt.Println("exit")
	} else {
		fmt.Println("\nAvailable Commands")
		fmt.Println("------------------")
		fmt.Println("whoami")
		fmt.Println("logout")
		fmt.Println("enable-2fa")
		fmt.Println("disable-2fa")
		fmt.Println("help")
		fmt.Println("exit")
	}
}
func (c *CLI) Enable2FA(rl *readline.Instance) {

	if c.currentToken == "" {
		fmt.Println("Please login first")
		return
	}

	secret, url, err := c.authService.Enable2FA(c.currentToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = qrcode.WriteFile(url, qrcode.Medium, 256, "qrcode.png")
	if err != nil {
		fmt.Println("Failed to generate QR:", err)
		return
	}

	fmt.Println("\n========== MFA SETUP ==========")
	fmt.Println("Add this secret to Google Authenticator")
	fmt.Println()
	fmt.Println(secret)
	fmt.Println()

	rl.SetPrompt("Enter OTP: ")
	otp, _ := rl.Readline()

	err = c.authService.Confirm2FA(
		c.currentToken,
		secret,
		strings.TrimSpace(otp),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("✅ MFA Enabled Successfully")

	c.updatePrompt(rl)
}
func (c *CLI) Disable2FA() {

	if c.currentToken == "" {
		fmt.Println("Please login first")
		return
	}

	err := c.authService.Disable2FA(c.currentToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("✅ MFA Disabled")
}
