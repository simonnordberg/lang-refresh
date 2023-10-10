# frozen_string_literal: true

module AlbumImport
  class Importer
    def initialize(client)
      @client = client
    end

    def import_albums_from_directory(dir)
      return false unless Dir.exist? dir
      Dir.foreach(dir) { |album_dir| import_album_from_directory(File.join(dir, album_dir)) }
    end

    def import_album_from_directory(album_path, metadata_filename = "metadata.json")
      return false unless Dir.exist? album_path

      metadata_file = File.join(album_path, metadata_filename)
      return false unless File.exist? metadata_file

      metadata = read_metadata(metadata_file)
      return false unless metadata.present?

      album = create_album(metadata)
      import_album_images(album_path, album) if album.present?
    end

    def import_album_images(album_path, album, suffixes = %w[.jpg .jpeg .png])
      Dir.foreach(album_path) do |filename|
        image = File.join(album_path, filename)
        next unless File.exist? image
        next unless suffixes.include? File.extname(image)
        sha = sha1sum(image)
        photo = get_photo(sha)
        add_album_photo(album, photo) if photo.present?
      end
    end

    def sha1sum(filename)
      Digest::SHA1.file(filename).hexdigest
    end

    def add_album_photo(album, photo)
      @client.add_album_photo(album, photo)
    end

    def read_metadata(metadata)
      JSON.parse(File.read(metadata))
    end

    def get_photo(sha)
      @client.get_photo(sha)
    end

    def create_album(metadata)
      @client.create_album(metadata["title"],
                           metadata["description"])
    end
  end
end
