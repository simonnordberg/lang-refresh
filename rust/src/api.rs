use reqwest::Client as HttpClient;
use reqwest::Error as HttpError;
use reqwest::Response as HttpResponse;
use serde::{Deserialize, Serialize};
use url::Url;

const SESSION_ENDPOINT: &str = "/api/v1/session";
const ALBUMS_ENDPOINT: &str = "/api/v1/albums";
const FILES_ENDPOINT: &str = "/api/v1/files";

#[derive(Debug, Serialize, Deserialize)]
pub struct Album {
    #[serde(alias = "UID")]
    pub uid: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Photo {
    #[serde(alias = "PhotoUID")]
    pub uid: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct AddPhotosResponse {
    added: Vec<Photo>,
}

pub struct Client {
    base_uri: Url,
    session_id: Option<String>,
    http_client: HttpClient,
}

impl Client {
    pub fn new(base_uri: Url) -> Self {
        Client {
            base_uri,
            session_id: None,
            http_client: HttpClient::new(),
        }
    }
}

#[derive(Debug, Deserialize)]
struct AuthenticationResponse {
    id: String,
}

impl Client {
    pub async fn authenticate(
        &mut self,
        username: &str,
        password: &str,
    ) -> Result<bool, HttpError> {
        let url = self.base_uri.join(SESSION_ENDPOINT).expect("invalid url");
        let payload = serde_json::json!({
            "username": username,
            "password": password,
        });

        let resp: AuthenticationResponse = self
            .http_client
            .post(url)
            .json(&payload)
            .send()
            .await?
            .json()
            .await?;

        self.session_id = Some(resp.id);
        Ok(self.session_id.is_some())
    }

    pub async fn create_album(&self, title: &str, description: &str) -> Result<Album, HttpError> {
        let url = self.base_uri.join(ALBUMS_ENDPOINT).expect("invalid url");
        let payload = serde_json::json!({
            "Title": title,
            "Description": description,
        });

        let response = self._post(&url, Some(payload)).await?;
        match response.error_for_status() {
            Ok(r) => Ok(r.json().await?),
            Err(err) => Err(err),
        }
    }

    pub async fn get_photo(&self, sha: &str) -> Result<Photo, HttpError> {
        let url = self
            .base_uri
            .join(&format!("{}/{}", FILES_ENDPOINT, sha))
            .expect("invalid url");

        let response = self._get(&url, None).await?;
        match response.error_for_status() {
            Ok(r) => Ok(r.json().await?),
            Err(err) => Err(err),
        }
    }

    pub async fn add_album_photo(
        &self,
        album: &Album,
        photo: &Photo,
    ) -> Result<Vec<Photo>, HttpError> {
        let url = self
            .base_uri
            .join(&format!("{}/{}/photos", ALBUMS_ENDPOINT, album.uid))
            .expect("invalid url");

        let payload = serde_json::json!({
            "photos": [photo.uid],
        });

        let response = self._post(&url, Some(payload)).await?;
        match response.error_for_status() {
            Ok(r) => Ok(r.json::<AddPhotosResponse>().await?.added),
            Err(err) => Err(err),
        }
    }

    pub async fn logout(&mut self) -> Result<bool, HttpError> {
        if let Some(session_id) = &self.session_id {
            let url = self
                .base_uri
                .join(&format!("{}/{}", SESSION_ENDPOINT, session_id))
                .expect("invalid url");
            let response = self._delete(&url).await?;
            match response.error_for_status() {
                Ok(_) => {
                    self.session_id = None;
                    return Ok(true);
                }
                Err(_) => Ok(false),
            }
        } else {
            Ok(false)
        }
    }

    fn _default_headers(&self) -> reqwest::header::HeaderMap {
        let mut headers = reqwest::header::HeaderMap::new();
        if let Some(session_id) = &self.session_id {
            headers.insert("X-Session-ID", session_id.parse().unwrap());
        }
        return headers;
    }

    async fn _get(
        &self,
        url: &Url,
        params: Option<Vec<(&str, &str)>>,
    ) -> Result<HttpResponse, HttpError> {
        let response = self
            .http_client
            .get(url.clone())
            .query(&params)
            .headers(self._default_headers())
            .send()
            .await?;
        Ok(response)
    }

    async fn _post(
        &self,
        url: &Url,
        json: Option<serde_json::Value>,
    ) -> Result<HttpResponse, HttpError> {
        let response = self
            .http_client
            .post(url.clone())
            .json(&json)
            .headers(self._default_headers())
            .send()
            .await?;
        Ok(response)
    }

    async fn _delete(&self, url: &Url) -> Result<HttpResponse, HttpError> {
        let response = self
            .http_client
            .delete(url.clone())
            .headers(self._default_headers())
            .send()
            .await?;
        Ok(response)
    }
}
