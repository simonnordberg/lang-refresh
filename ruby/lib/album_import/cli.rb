# frozen_string_literal: true

require "optparse"
require "httparty"

module AlbumImport
  class CLI
    def run(argv)
      opts = parse_opts! argv
      run_import(opts[:host], opts[:username], opts[:password], opts[:directory])
    rescue OptionParser::MissingArgument => e
      puts e.message
      exit 1
    end

    def run_import(hostname, username, password, directory)
      client = Client.new(hostname)

      begin
        if client.authenticate?(username, password)
          importer = Importer.new(client)
          importer.import_albums_from_directory(directory)
        end
      rescue => e
        puts "Error: #{e.backtrace}"
        raise
      ensure
        client.logout
      end
    end

    def parse_opts!(argv)
      options = {}
      OptionParser.new do |opts|
        opts.banner = "Usage: album_import [options]"
        opts.on("--directory directory", "Directory containing images")
        opts.on("--host host", "Base URL")
        opts.on("--username username", "Username")
        opts.on("--password password", "Password")
      end.parse! argv, into: options

      %i[directory host username password].each do |opt|
        raise OptionParser::MissingArgument, "Option --#{opt} is required" unless options[opt]
      end
      options
    end
  end
end
