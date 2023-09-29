from urllib.parse import urljoin

import requests

from model import Album, Photo

SESSION_ENDPOINT = "/api/v1/session"
ALBUMS_ENDPOINT = "/api/v1/albums"
FILES_ENDPOINT = "/api/v1/files"


class Client:
    def __init__(self, base_uri):
        self.base_uri = base_uri
        self.session_id = None

    def authenticate(self, username, password):
        url = urljoin(self.base_uri, SESSION_ENDPOINT)
        payload = {
            "username": username,
            "password": password,
        }
        response = self._post(url, json=payload)
        if response.status_code == 200:
            self.session_id = response.json().get("id")
        return self.session_id is not None

    def get_albums(self):
        url = urljoin(self.base_uri, ALBUMS_ENDPOINT)
        params = {
            "count": 10,
        }
        response = self._get(url, params)
        if response.status_code == 200:
            return [*map(lambda x: Album.from_dict(x), response.json())]
        else:
            return []

    def create_album(self, title, description):
        url = urljoin(self.base_uri, ALBUMS_ENDPOINT)
        payload = {
            "Title": title,
            "Description": description
        }
        response = self._post(url, json=payload)
        if response.status_code == 200:
            return Album.from_dict(response.json())
        else:
            return False

    def get_photo(self, sha):
        url = urljoin(self.base_uri,
                      "{api}/{sha}".format(
                          api=FILES_ENDPOINT,
                          sha=sha
                      ))
        response = self._get(url)
        if response.status_code == 200:
            return Photo.from_dict(response.json())
        else:
            return False

    def add_album_photo(self, album, photo):
        url = urljoin(self.base_uri,
                      "{api}/{id}/photos".format(
                          api=ALBUMS_ENDPOINT,
                          id=album.uid)
                      )
        payload = {
            "photos": [photo.uid]
        }
        response = self._post(url, json=payload)
        if response.status_code == 200:
            return [*map(lambda x: Photo.from_dict(x), response.json()["added"])]
        else:
            return []

    def logout(self):
        if self.session_id:
            url = urljoin(urljoin(self.base_uri, SESSION_ENDPOINT), self.session_id)
            response = requests.delete(url)
            self.session_id = None
            return response.status_code == 200
        return False

    def _default_headers(self):
        return {
            "X-Session-ID": self.session_id
        }

    def _get(self, url, params=None) -> requests.Response:
        return requests.get(url, params=params, headers=self._default_headers())

    def _post(self, url, json=None) -> requests.Response:
        return requests.post(url, json=json, headers=self._default_headers())

    def _delete(self, url) -> requests.Response:
        return requests.delete(url, headers=self._default_headers())
