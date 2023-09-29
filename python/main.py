import argparse
import hashlib
import json
import logging
import os
import sys
from pathlib import Path

from api import Client


def sha1sum(filename):
    with open(filename, 'rb', buffering=0) as f:
        return hashlib.file_digest(f, 'sha1').hexdigest()


class Importer:
    def __init__(self, client):
        self.client = client

    def import_albums_from_directory(self, path):
        albums_path = Path(path)
        if not albums_path.exists():
            logging.getLogger().error("does not exist: %s", path)
            return False

        for directory in os.scandir(albums_path):
            logging.getLogger().info("Importing directory: %s", directory)
            self.import_album_from_directory(directory)

    def import_album_from_directory(self, path, metadata_filename="metadata.json"):
        album_path = Path(path)
        metadata_file = album_path.joinpath(metadata_filename)

        if not album_path.exists():
            logging.getLogger().error("does not exist: %s", path)
            return False

        if not album_path.is_dir():
            logging.getLogger().error("not a directory: %s", path)
            return False

        if not metadata_file.exists():
            logging.getLogger().error("does not exist: %s", metadata_file)
            return False

        metadata = self._read_metadata(metadata_file)
        if not metadata:
            logging.getLogger().error("unable to parse metadata %s", metadata_file)
            return False

        album = self.create_album(metadata)
        if album:
            self.import_album_images(album_path, album)

    def import_album_images(self, album_path, album, suffixes=None):
        if suffixes is None:
            suffixes = [".jpg", ".jpeg", ".png"]

        for file in album_path.glob("*"):
            if not file.is_file():
                logging.getLogger().debug("ignoring: %s", file)
                continue

            if file.suffix not in suffixes:
                logging.getLogger().debug("ignoring: %s", file)
                continue

            sha = sha1sum(file)
            photo = self.client.get_photo(sha)
            if photo:
                photos = self.client.add_album_photo(album, photo)
                logging.getLogger().debug("added %d photos to album %s", len(photos), album.title)

    def create_album(self, metadata):
        return self.client.create_album(metadata.get("title"),
                                        metadata.get("description"))

    @staticmethod
    def _read_metadata(metadata_file):
        with open(metadata_file) as f:
            metadata = json.loads(f.read())
            return metadata


def main(args):
    client = Client(args.get("host"))
    try:
        if not client.authenticate(args.get("username"), args.get("password")):
            logging.getLogger().error("authentication failed")
            sys.exit(-1)
        importer = Importer(client)
        importer.import_albums_from_directory(args.get("directory"))
    finally:
        client.logout()


if __name__ == '__main__':
    logging.basicConfig()
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser(description="Import Google Takeout export to Photoprism",
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument("-H", "--host", help="Base URL", default="http://localhost:2342")
    parser.add_argument("-u", "--username", help="Username", default="username")
    parser.add_argument("-p", "--password", help="Password", default="password")
    parser.add_argument("directory", help="Google Photos directory (containing albums)", type=Path)

    args = vars(parser.parse_args())
    main(args)
