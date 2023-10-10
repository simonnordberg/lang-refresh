# frozen_string_literal: true

module AlbumImport
  Album = Struct.new(:uid) do
    def self.from_dict(dict)
      new(dict["UID"])
    end
  end

  Photo = Struct.new(:uid) do
    def self.from_dict(dict)
      new(dict["PhotoUID"])
    end
  end
end
