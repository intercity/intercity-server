require "fileutils"
require "highline"

module IntercityServer
  class Installer
    attr_reader :hostname, :from_email

    def self.execute
      Installer.new.execute
    end

    def execute
      cli = HighLine.new

      @hostname = cli.ask("What is the hostname? (e.g.: intercity.example.com)") do |q|
        q.validate = hostname_regex
      end
      cli.say "Hostname is set to #{@hostname}"

      @from_email = cli.ask("What email should we use as sender (For example for activtion emails)") do |q|
        q.validate = /[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}/
      end
      cli.say "From email is set to #{@from_email}"

      cli.say "---- Installing docker"
      install_docker

      cli.say "---- Downloading Intercity"
      clone_intercity

      cli.say "---- Configuring Intercity"
      copy_configuration
      replace_values

      cli.say "---- Done"
    end

    private

    def install_docker
      `wget -nv -O - https://get.docker.com/ | sh`
    end

    def clone_intercity
      FileUtils.mkdir_p "/var/intercity"
      `git clone https://github.com/intercity/intercity-docker.git -b 0-3-stable /var/intercity`
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
      config_content = config_content.gsub(/example\@example.com/, from_email)

      File.open(config_file, "w") {|file| file.puts config_content }
    end
  end
end
