require "fileutils"
require "highline"

module IntercityServer
  class Installer
    attr_reader :hostname, :use_ssl, :letsencrypt_email

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

      cli.choose do |menu|
        menu.prompt = "Do you want to use LetsEncrypt for SSL?\n" \
          "IMPORTANT: The hostname should be public and reachable for the LetsEncrypt servers.\n"\
          "The SSL certificate can't be generated if LetsEncrypt can't reach the domain!"
        menu.choice(:yes) { @use_ssl = true }
        menu.choices(:no) { @use_ssl = false }
      end

      if use_ssl
        @letsencrypt_email = cli.ask("What is the email address we can use for LetsEncrypt") do |q|
          q.validate  = email_regex
        end
      end

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

      cli.say "---- Done\n\n"
      if use_ssl
        cli.say ""
        cli.say "================="
        cli.say "== IMPORTANT:  =="
        cli.say "================="
        cli.say "Keep in mind that it can take up to 3 minutes until your Intercity instance is reachable over HTTPS."
        cli.say "This is due to the delay at Lets Encrypt with issueing the certificates"
      end
    end

    private

    def install_docker
      `wget -nv -O - https://get.docker.com/ | sh`
    end

    def clone_intercity
      FileUtils.mkdir_p "/var/intercity"
      `git clone https://github.com/intercity/intercity-docker.git -b 0-4-stable /var/intercity`
    end

    def copy_configuration
      `cp /var/intercity/samples/app.yml /var/intercity/containers/`
    end

    def hostname_regex
      /(?!.{256})(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9] )?\.)+(?:[a-z]{1,63}|xn--[a-z0-9]{1,59})/
    end

    def email_regex
      /([^@]+)@([^\.]+)/
    end

    def replace_values
      config_file = "/var/intercity/containers/app.yml"
      config_content = File.read config_file
      config_content = config_content.gsub(/intercity\.example\.com/, hostname)

      if use_ssl
        config_content = config_content.gsub(/#- "templates\/web\.ssl\.template.yml"/, '- "templates/web.ssl.template.yml"')
        config_content = config_content.gsub(/#- "templates\/web\.letsencrypt\.ssl\.template.yml"/,
                                             '- "templates/web.letsencrypt.ssl.template.yml"')
        config_content = config_content.gsub(/LETSENCRYPT_ACCOUNT_EMAIL: "example@example.com"/,
                                             "LETSENCRYPT_ACCOUNT_EMAIL: \"#{letsencrypt_email}\"")
      end

      File.open(config_file, "w") {|file| file.puts config_content }
    end

    def build_intercity
      `/var/intercity/launcher bootstrap app`
    end

    def start_intercity
      `/var/intercity/launcher start app`
      `/var/intercity/launcher restart app` if use_ssl
    end

    def check_existing_install
      return unless Dir.exist?("/var/intercity")
      HighLine.new.say "Looks like Intercity is already installed."
      exit 1
    end
  end
end
