package main

import (
	"log"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("intercity-server", "Manage your Intercity instance with ease")

	version = app.Command("version", "Show the version of installed Intercity and Intercity server")

	install         = app.Command("install", "Install Intercity on this server.")
	installHostname = install.Arg("hostname", "The hostname you want to run Intercity on.").Required().String()

	update = app.Command("update", "Update your Intercity instance.")

	restart = app.Command("restart", "Restart your Intercity instance")
)

func main() {
	kingpin.Version("0.3.0")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case install.FullCommand():
		if _, err := os.Stat("/var/intercity"); os.IsNotExist(err) {
			installDocker()
			downloadIntercity()
			configureIntercity()
			buildIntercity()
			startIntercity()
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
		println("intercity-server:", "0.2.0")
		println("intercity-docker:", "0.4.1")
		println("intercity-web:", "0.2.0")
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
