# frozen_string_literal: true

module AlbumImport
  class Client
    SESSION_ENDPOINT = "/api/v1/session"
    ALBUMS_ENDPOINT = "/api/v1/albums"
    FILES_ENDPOINT = "/api/v1/files"

    include HTTParty

    def initialize(hostname)
      self.class.base_uri hostname
      @session_id = nil
    end

    def authenticate?(username, password)
      response = self.class.post(
        SESSION_ENDPOINT,
        body: JSON.generate({username: username, password: password}),
        headers: default_headers
      )
      @session_id = response["id"]
      @session_id.present?
    end

    def get_photo(sha)
      response = self.class.get(
        "#{FILES_ENDPOINT}/#{sha}",
        headers: default_headers
      )
      Photo.from_dict(response)
    end

    def get_albums
      self.class.get(
        ALBUMS_ENDPOINT,
        headers: default_headers,
        query: {count: 10}
      ).map { |e| Album.from_dict(e) }
    end

    def create_album(title, description)
      response = self.class.post(
        ALBUMS_ENDPOINT,
        headers: default_headers,
        body: JSON.generate(
          {
            Title: title,
            Description: description
          }
        )
      )
      Album.from_dict(response)
    end

    def add_album_photo(album, photo)
      self.class.post(
        "#{ALBUMS_ENDPOINT}/#{album.uid}/photos",
        headers: default_headers,
        body: JSON.generate(
          {
            photos: [photo.uid]
          }
        )
      ).success?
    end

    def logout
      self.class.delete("#{SESSION_ENDPOINT}/#{@session_id}") if @session_id.present?
      @session_id = nil
      @session_id.blank?
    end

    private

    def default_headers
      {
        "X-Session-ID": @session_id,
        "Content-Type": "application/json"
      }
    end
  end
end
