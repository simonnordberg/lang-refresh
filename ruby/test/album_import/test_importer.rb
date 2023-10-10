# frozen_string_literal: true

require "test_helper"

module AlbumImport
  class TestImporter < Minitest::Test
    def test_import_photo_album
      album_dir = File.join(File.dirname(__FILE__), "../fixtures/albums/album1")
      photos = %w[photo1.jpg photo2.jpg].sort

      mock_client = Minitest::Mock.new
      importer = Importer.new(mock_client)
      album = Album.new("album_uid")

      # Mock expectations
      mock_client.expect :create_album, album, ["Album 1 title", "Album 1 description"]
      photos.each do |photo|
        sha = checksum(File.join(album_dir, photo))
        photo = Photo.new(sha)

        mock_client.expect :get_photo, photo, [sha]
        mock_client.expect :add_album_photo, true, [album, photo]
      end

      importer.import_album_from_directory(album_dir)

      mock_client.verify
    end

    def test_import_multiple_albums
      albums_dir = File.join(File.dirname(__FILE__), "../fixtures/albums")

      mock_client = Minitest::Mock.new
      importer = Importer.new(mock_client)

      2.times { mock_client.expect :create_album, Object, [String, String] }
      3.times do
        mock_client.expect :get_photo, Object, [String]
        mock_client.expect :add_album_photo, true, [Object, Object]
      end

      importer.import_albums_from_directory(albums_dir)

      mock_client.verify
    end

    private

    def checksum(file)
      Digest::SHA1.file(file).hexdigest
    end
  end
end