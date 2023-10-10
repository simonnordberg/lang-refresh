# frozen_string_literal: true

require "test_helper"

module AlbumImport
  class TestCLI < Minitest::Test
    def test_parse_opts_missing
      assert_raises(OptionParser::MissingArgument, "Should raise missing argument") do
        CLI.new.send(:parse_opts!, %w[])
      end
    end

    def test_parse_opts_successful
      opts = CLI.new.send(:parse_opts!, %w[
        --host http://localhost:1111
        --directory /foo/bar
        --username user1name
        --password pass2word
      ])
      assert_equal({
                     host: "http://localhost:1111",
                     directory: "/foo/bar",
                     username: "user1name",
                     password: "pass2word",
                   }, opts)

    end
  end
end
