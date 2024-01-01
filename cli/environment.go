package cli

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
	"strings"
	"runtime"

	"github.com/joho/godotenv"
)

var deb_base = []string {"Ubuntu", "Pop", "Debian", "MX", "Raspbian", "Kali"}
var rh_base = []string {"Fedora", "Rocky", "AlmaLinux", "CentOS", "RedHatEnterpriseServer", "Oracle", "ClearOS", "AmazonAMI"}
var arch_base = []string {"Arch", "Garuda", "Manjaro", "Endeavour"}

// Load envfile and get environment variables
func GetEnv() {

	// Get home dir path
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Can't find home directory of system...", err)
		os.Exit(1)
	}

	env.HomeDir = homedir
	env.DBDir = homedir + "/.config/app"

	dir := homedir + "/.config/app"
	dBase := ""

	// Getting settings from settings.conf if it exists
	if _, conferr := os.Stat(dir + "/settings.conf"); conferr == nil {
		err := godotenv.Load(dir + "/settings.conf")
		if err != nil {
			fmt.Println(utils.ColorRed, "Error loading settings.conf", utils.ColorReset)
		}
		

		if yay := os.Getenv("YAY"); yay == "n" {
			env.Yay = false
		}
		if flatpak := os.Getenv("FLATPAK"); flatpak == "n" {
			env.Flatpak = false
		}
		if snap := os.Getenv("SNAP"); snap == "n" {
			env.Snap = false
		}
		if appimage := os.Getenv("APPIMAGE"); appimage == "n" {
			env.AppImage = false
		}
	} else {
		if _, err := os.Stat(dir); os.IsNotExist(err) {	
			fmt.Println(utils.ColorYellow, "\nFirst time running... Creating config dir...\n\n", utils.ColorReset)
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Exiting...\n")
				os.Exit(1)
			}
		} 

		if werr := utils.WriteToFile("YAY = n\nFLATPAK = n\nSNAP = n\nAPPIMAGE = n", dir + "/settings.conf"); werr != nil {
			utils.PrintErrorExit("Write settings.conf Error:", werr)
		}
	}

	// determine OS
	osType := runtime.GOOS
	switch osType {
	case "linux":
		env.OSType = "Linux"
	case "darwin":
		env.OSType = "Mac"
	case "windows":
		env.OSType = "Windows"
	default:
		fmt.Print(utils.ColorRed, "Unsupported Operating System... Exiting...\n\n", utils.ColorReset)
		os.Exit(1)
	}

	if osType == "linux" {
		d, err := exec.Command("/usr/bin/lsb_release", "-i", "-s").Output()
		if err != nil {
			fmt.Println("Error:", err)
		}
		distro := strings.TrimSuffix(string(d), "\n")

		//fmt.Println("Distro:", distro)
		env.Distro = distro
		
		// Check if it's a Debian based
		for i := 0; i < len(deb_base); i++ {
			if distro == deb_base[i] {
				dBase = "debian"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Check if it's a RedHat based
		for i := 0; i < len(rh_base); i++ {
			if distro == rh_base[i] {
				dBase = "redhat"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Check if it's a Arch based
		for i := 0; i < len(arch_base); i++ {
			if distro == arch_base[i] {
				dBase = "arch"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}
		
		// Temporarily disabling package manager that don't exist
		if env.Base == "arch" && env.Yay == true {
			yay, _ := utils.CheckIfExists(env.YayCmd)
			if yay == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Yay because it's not installed on your system. Suppress this message by disabling Yay on app by running \"app -m yay disable\"...\n", utils.ColorReset)
			}
		}

		if env.Flatpak == true {
			flatpak, _ := utils.CheckIfExists(env.FlatpakCmd)
			if flatpak == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Flatpak because it's not installed on your system. Suppress this message by disabling Flatpak on app by running \"app -m flatpak disable\"...\n", utils.ColorReset)
			}
		}

		if env.Snap == true {
			snap, _ := utils.CheckIfExists(env.SnapCmd)
			if snap == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Snap because it's not installed on your system. Suppress this message by disabling Snap on app by running \"app -m snap disable\"...\n", utils.ColorReset)
			}
		}
	}	

}

