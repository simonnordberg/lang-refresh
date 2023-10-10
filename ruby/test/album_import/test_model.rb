# frozen_string_literal: true

require "test_helper"

module AlbumImport
  class ModelTest < Minitest::Test
    def test_album
      album = Album.from_dict({ "UID" => "abc" })
      assert_equal("abc", album.uid)
    end

    def test_photo
      photo = Photo.from_dict({ "PhotoUID" => "abc" })
      assert_equal("abc", photo.uid)
    end
  end
end
