package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("intercity-server", "Manage your Intercity instance with ease")

	version = app.Command("version", "Show the version of installed Intercity and Intercity server")

	install           = app.Command("install", "Install Intercity on this server.")
	installHostname   = install.Flag("hostname", "The hostname you want to run Intercity on.").Required().String()
	installCustomPort = install.Flag("custom-port", "Use custom port for Intercity").Bool()
	installUseSSL     = install.Flag("use-ssl", "Enable SSL using Let's Encrypt").Bool()
	installSSLEmail   = install.Flag("ssl-email", "Email address to use for Let's Encrypt").String()

	update = app.Command("update", "Update your Intercity instance.")

	restart = app.Command("restart", "Restart your Intercity instance")

	current_cli_version = "0.3.0"
)

func main() {
	kingpin.Version(current_cli_version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case install.FullCommand():
		if _, err := os.Stat("/var/intercity"); os.IsNotExist(err) {
			if !checkValidDomain(*installHostname) {
				println("Hostname is not valid.")
				println("Installation cannot continue.")
				os.Exit(1)
			}

			if *installUseSSL && !checkValidEmail(*installSSLEmail) {
				println("In order to enable SSL you need to provide a valid email address")
				println("that we can register with Let's Encrypt.")
				println("You can do so with the '--ssl-email=' flag")
				os.Exit(1)
			}

			installDocker()
			downloadIntercity()
			configureIntercity()
			buildIntercity()
			startIntercity()

			println("")
			println("Congratulations! Intercity is now installed.")
			println("You can reach your brand new Intercity installation on:")

			if *installCustomPort {
				if *installUseSSL {
					println(fmt.Sprintf("    %v%v:8443", determineProtocol(), *installHostname))
				} else {
					println(fmt.Sprintf("    %v%v:880", determineProtocol(), *installHostname))
				}
			} else {
				println(fmt.Sprintf("    %v%v", determineProtocol(), *installHostname))
			}

			if *installUseSSL {
				println("====================")
				println("=     WARNING      =")
				println("====================")
				println("Please keep in mind that Let's encrypt can take up to 5 minutes")
				println("to issue the certificate. Untill this is done")
				println("your Intercity server will be unreachable")
			}

		} else {
			println("Intercity is already installed.")
			println("If you want to update your Intercity instance, run:")
			println("\t intercity-server update")
		}

	case update.FullCommand():
		if _, err := os.Stat("/var/intercity"); os.IsNotExist(err) {
			println("Intercity is not installed.")
			println("To install Intercity, run:")
			println("\t intercity-server install")
		} else {
			updateIntercity()
		}

	case restart.FullCommand():
		if _, err := os.Stat("/var/intercity"); os.IsNotExist(err) {
			println("Intercity is not installed.")
			println("To install Intercity, run:")
			println("\t intercity-server install")
		} else {
			restartIntercity()
		}

	case version.FullCommand():
		println("intercity-server:", current_cli_version)
	}
}

func installDocker() {
	println("---- Installing Docker")

	if _, err := runCommand("which docker"); err != nil {
		if _, err := runCommand("wget -nv -O - https://get.docker.com | sh"); err != nil {
			println("Could not install Docker")
			log.Fatal(err)
		} else {
			println("     Done")
		}
	} else {
		println("     Docker is already installed. Let's continue")
	}
}

func downloadIntercity() {
	println("---- Downloading Intercity")

	path := filepath.Join("/var", "intercity")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		println("Could not create Intercity installation directory.")
		println(err.Error())
		log.Fatal(err)
	}

	cmd := "git clone https://github.com/intercity/intercity-docker.git -b 0-4-stable /var/intercity"
	if _, err := runCommand(cmd); err != nil {
		println("Could not download Intercity")
		log.Fatal(err)
	}
	println("     Done")
}

func configureIntercity() {
	println("---- Configuring Intercity")

	if _, err := runCommand("cp /var/intercity/samples/app.yml /var/intercity/containers/"); err != nil {
		log.Fatal(err)
	}

	configFile := "/var/intercity/containers/app.yml"
	if err := replaceData(configFile, "intercity.example.com", *installHostname); err != nil {
		log.Fatal(err)
	}

	if *installCustomPort {
		configFile := "/var/intercity/containers/app.yml"
		if err := replaceData(configFile, "80:80", "8880:80"); err != nil {
			log.Fatal(err)
		}

		if err := replaceData(configFile, "443:443", "8443:443"); err != nil {
			log.Fatal(err)
		}
	}

	if *installUseSSL {
		configFile := "/var/intercity/containers/app.yml"
		if err := replaceData(configFile, "#- \"templates/web.ssl.template.yml\"",
			"- \"templates/web.ssl.template.yml\""); err != nil {
			log.Fatal(err)
		}

		if err := replaceData(configFile, "#- \"templates/web.letsencrypt.ssl.template.yml\"",
			"- \"templates/web.letsencrypt.ssl.template.yml\""); err != nil {
			log.Fatal(err)
		}

		if err := replaceData(configFile, "LETSENCRYPT_ACCOUNT_EMAIL: \"example@example.com\"",
			"LETSENCRYPT_ACCOUNT_EMAIL: \""+*installSSLEmail+"\""); err != nil {
			log.Fatal(err)
		}
	}

	println("     Done")
}

func buildIntercity() {
	println("---- Building Intercity")
	if _, err := runCommand("/var/intercity/launcher bootstrap app"); err != nil {
		log.Fatal(err)
	}

	println("     Done")
}

func startIntercity() {
	println("---- Starting Intercity")

	if _, err := runCommand("/var/intercity/launcher start app"); err != nil {
		log.Fatal(err)
	}

	if _, err := runCommand("/var/intercity/launcher restart app"); err != nil {
		log.Fatal(err)
	}

	println("     Done")
}

func updateIntercity() {
	println("---- Updating Intercity")
	if _, err := runCommand("/var/intercity/launcher rebuild app"); err != nil {
		log.Fatal(err)
	}
	println("     Done")
}

func restartIntercity() {
	println("---- Restarting Intercity")
	if _, err := runCommand("/var/intercity/launcher restart app"); err != nil {
		log.Fatal(err)
	}
	println("     Done")
}

func determineProtocol() string {
	if *installUseSSL {
		return "https://"
	} else {
		return "http://"
	}
}
