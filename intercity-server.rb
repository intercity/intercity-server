#!/usr/bin/env ruby

require_relative "lib/intercity_server/installer"

class IntercityServerCli
  def self.start
    case ARGV[0]
    when "help"
      IntercityServerCli.new.usage
    when "install"
      IntercityServerCli.new.install
    when "restart"
      IntercityServerCli.new.restart
    when "update"
      IntercityServerCli.new.update
    else
      IntercityServerCli.new.usage
    end
  end

  def usage
    puts "Usage: \tintercity-server COMMAND"
    puts ""
    puts "Commands:"
    puts "    help - Show the commands available"
    puts "    install - Run the setup for installing intercity-server"
    puts "    restart - Restart your Intercity instance"
    puts "    update - Update your Intercity instance"
  end

  def install
    if installed?
      puts "Intercity is already installed."
      exit 1
    end
    IntercityServer::Installer.execute
  end

  def restart
    ensure_installed
    `/var/intercity/launcher restart app`
  end

  def update
    ensure_installed
    `/var/intercity/launcher rebuild app`
  end

  private

  def installed?
    Dir.exist?("/var/intercity")
  end

  def ensure_installed
    return if installed?
    puts "Intercity is not yet installed."
    puts "To install Intercity run:"
    puts "   intercity-server install"
    exit 1
  end
end

IntercityServerCli.start
