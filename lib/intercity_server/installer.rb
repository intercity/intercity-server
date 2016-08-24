require "fileutils"
require "highline"

module IntercityServer
  class Installer
    attr_reader :hostname

    def self.execute
      Installer.new.execute
    end

    def execute
      check_existing_install

      cli = HighLine.new

      @hostname = cli.ask("What is the hostname? (e.g.: intercity.example.com)") do |q|
        q.validate = hostname_regex
      end
      cli.say "Hostname is set to #{@hostname}"

      cli.say "---- Installing docker"
      install_docker

      cli.say "---- Downloading Intercity"
      clone_intercity

      cli.say "---- Configuring Intercity"
      copy_configuration
      replace_values

      cli.say "---- Building Intercity"
      build_intercity

      cli.say "---- Starting Intercity"
      start_intercity

      cli.say "---- Done"
    end

    private

    def install_docker
      `wget -nv -O - https://get.docker.com/ | sh`
    end

    def clone_intercity
      FileUtils.mkdir_p "/var/intercity"
      `git clone https://github.com/intercity/intercity-docker.git -b master /var/intercity`
    end

    def copy_configuration
      `cp /var/intercity/samples/app.yml /var/intercity/containers/`
    end

    def hostname_regex
      /(?!.{256})(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9] )?\.)+(?:[a-z]{1,63}|xn--[a-z0-9]{1,59})/
    end

    def replace_values
      config_file = "/var/intercity/containers/app.yml"
      config_content = File.read config_file
      config_content = config_content.gsub(/intercity\.example\.com/, hostname)

      File.open(config_file, "w") {|file| file.puts config_content }
    end

    def build_intercity
      `/var/intercity/launcher bootstrap app`
    end

    def start_intercity
      `/var/intercity/launcher start app`
    end

    def check_existing_install
      return unless Dir.exist?("/var/intercity")
      HighLine.new.say "Looks like Intercity is already installed."
      exit 1
    end
  end
end
