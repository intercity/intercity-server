#!/usr/bin/env ruby

require_relative "lib/intercity_server/installer"

class IntercityServerCli
  def self.start
    case ARGV[0]
    when "help"
      IntercityServerCli.new.usage
    when "install"
      IntercityServerCli.new.install
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
  end

  def install
    IntercityServer::Installer.execute
  end
end

IntercityServerCli.start
