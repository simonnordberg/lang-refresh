use serde_json::Value;
use std::env;
use std::error::Error;
use std::fs;
use std::path::Path;
use tokio::fs::File;
use tokio::io::AsyncReadExt;
use url::Url;
mod api;
use api::Client;
use std::collections::HashSet;

async fn sha1sum(filename: &Path) -> Result<String, Box<dyn Error>> {
    let mut file = File::open(filename).await?;
    let mut hasher = openssl::sha::Sha1::new();

    let mut buffer = Vec::new();
    while let Ok(n) = file.read_to_end(&mut buffer).await {
        if n == 0 {
            break;
        }
        hasher.update(&buffer);
    }

    Ok(hex::encode(hasher.finish()))
}

struct Importer {
    client: Client,
}

impl Importer {
    async fn import_albums_from_directory(&self, path: &Path) -> Result<(), Box<dyn Error>> {
        let albums_path = Path::new(path);
        if !albums_path.exists() {
            eprintln!("does not exist: {:?}", path);
            return Ok(());
        }

        for entry in fs::read_dir(albums_path)? {
            let entry = entry?;
            let directory = entry.path();
            println!("Importing directory: {:?}", directory.display());
            self.import_album_from_directory(&directory, "metadata.json")
                .await?;
        }

        Ok(())
    }

    async fn import_album_from_directory(
        &self,
        path: &Path,
        metadata_filename: &str,
    ) -> Result<(), Box<dyn Error>> {
        let album_path = Path::new(path);
        let metadata_file = album_path.join(metadata_filename);

        if !album_path.exists() {
            eprintln!("does not exist: {:?}", path);
            return Ok(());
        }

        if !album_path.is_dir() {
            eprintln!("not a directory: {:?}", path);
            return Ok(());
        }

        if !metadata_file.exists() {
            eprintln!("does not exist: {:?}", metadata_file);
            return Ok(());
        }

        let metadata = self.read_metadata(&metadata_file).await?;
        let album = self.create_album(&metadata).await?;

        if let Some(album) = album {
            self.import_album_images(album_path, &album).await?;
        }

        Ok(())
    }

    async fn import_album_images(
        &self,
        album_path: &Path,
        album: &api::Album,
    ) -> Result<(), Box<dyn Error>> {
        let suffixes: HashSet<&str> = ["jpg", "jpeg", "png"].iter().cloned().collect();

        for entry in fs::read_dir(album_path)? {
            let entry = entry?;
            let file_path = entry.path();

            if !file_path.is_file() {
                println!("ignoring: {:?}", file_path.display());
                continue;
            }

            if !suffixes.contains(&file_path.extension().and_then(|e| e.to_str()).unwrap_or("")) {
                println!("ignoring: {:?}", file_path.display());
                continue;
            }

            let sha = sha1sum(&file_path).await?;
            let photo = self.client.get_photo(&sha).await?;
            self.client.add_album_photo(&album, &photo).await?;
            println!("added 1 photo to album {:?}", &album);
        }

        Ok(())
    }

    async fn create_album(&self, metadata: &Value) -> Result<Option<api::Album>, Box<dyn Error>> {
        let title = metadata["title"].as_str().unwrap_or_default();
        let description = metadata["description"].as_str().unwrap_or_default();

        let album = self.client.create_album(&title, &description).await?;
        Ok(Some(album))
    }

    async fn read_metadata(&self, metadata_file: &Path) -> Result<Value, Box<dyn Error>> {
        let contents = fs::read_to_string(metadata_file)?;
        let metadata: Value = serde_json::from_str(&contents)?;
        Ok(metadata)
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let args: Vec<String> = env::args().collect();
    if args.len() != 5 {
        eprintln!(
            "Usage: {} <host> <username> <password> <directory>",
            args[0]
        );
        std::process::exit(1);
    }

    let host = Url::parse(&args[1]).expect("invalid url");
    let username = &args[2];
    let password = &args[3];
    let directory = &args[4];

    println!("{} {} {} {}", host, username, password, directory);

    let mut importer = Importer {
        client: Client::new(host),
    };

    if !importer.client.authenticate(username, password).await? {
        eprintln!("Authentication failed");
        std::process::exit(1);
    }

    importer
        .import_albums_from_directory(Path::new(directory))
        .await?;
    importer.client.logout().await?;

    Ok(())
}
