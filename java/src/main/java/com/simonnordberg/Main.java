package com.simonnordberg;

import com.simonnordberg.model.Album;
import java.net.URI;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class Main {

  public static void main(String[] args) throws ExecutionException, InterruptedException {
    var client = new Client(URI.create("http://localhost:2342"));
    client.authenticate("admin", "insecure").thenRun(() -> {
      CompletableFuture<List<Album>> albums = client.getAlbums();
      try {
        System.out.println("albums = " + albums.get());
      } catch (InterruptedException e) {
        throw new RuntimeException(e);
      } catch (ExecutionException e) {
        throw new RuntimeException(e);
      }
    }).get();
  }
}
