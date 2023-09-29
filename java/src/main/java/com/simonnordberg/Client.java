package com.simonnordberg;


import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.annotations.SerializedName;
import com.simonnordberg.api.ApiError;
import com.simonnordberg.api.ApiResponse;

import java.net.http.HttpResponse.BodyHandlers;

import org.apache.http.client.utils.URIBuilder;

import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionException;

class NullSession extends Session {

  NullSession() {
    super(null);
  }
}

class Session {

  private final String id;

  Session(String id) {
    this.id = id;
  }

  public String getId() {
    return id;
  }

  public boolean isActive() {
    return id != null;
  }
}

record SessionCreateRequest(String username, String password) {

}

record SessionCreateResponse(String id) {

}


class RequestError extends CompletionException {

  ErrorResponse response;

  public RequestError(ErrorResponse response) {
    this.response = response;
  }

  @Override
  public String getMessage() {
    return this.response.error();
  }
}

record ErrorResponse(String error) {

}

public class Client {

  protected static final String SessionAPI = "/api/v1/session";
  protected static final String AlbumAPI = "/api/v1/albums";
  protected static final String FileAPI = "/api/v1/files";

  protected final Gson serializer;
  protected URI baseURI;
  protected Session session;

  public Client(URI baseURI) {
    this.baseURI = baseURI;
    this.serializer = new GsonBuilder().create();
    this.session = new NullSession();
  }

  public boolean isAuthenticated() {
    return this.session.isActive();
  }

  public CompletableFuture<Void> authenticate(String username, String password) {
    return asyncPostRequest(baseURI.resolve(SessionAPI),
        new SessionCreateRequest(username, password), SessionCreateResponse.class).thenApply(
        response -> {
          session = response.hasError() ? new NullSession() : new Session(response.getData().id());
          return null;
        }
    );
  }

  public CompletableFuture<Boolean> authenticate2(String username, String password) {
    return asyncPostRequest(baseURI.resolve(SessionAPI),
        new SessionCreateRequest(username, password), SessionCreateResponse.class).thenApply(
        response -> {
          if (response.hasError()) {
            session = new NullSession();
            return false;
          } else {
            session = new Session(response.getData().id());
            return true;
          }
        });
  }

  public CompletableFuture<List<com.simonnordberg.model.Album>> getAlbums() {
    try {
      URI uri = new URIBuilder(baseURI.resolve(AlbumAPI)).addParameter("count", "10").build();

      return asyncGetRequest(uri, com.simonnordberg.model.Album[].class).thenApply(
          response -> Arrays.asList(response.getData()));

    } catch (URISyntaxException e) {
      return new CompletableFuture<>();
    }
  }

  private <T> CompletableFuture<ApiResponse<T>> asyncGetRequest(URI uri, Class<T> responseType) {
    return CompletableFuture.supplyAsync(() -> getRequest(uri, responseType));
  }

  private <T> CompletableFuture<ApiResponse<T>> asyncPostRequest(URI uri, Object data,
      Class<T> responseType) {
    return CompletableFuture.supplyAsync(() -> postRequest(uri, data, responseType));
  }

  private <T> ApiResponse<T> getRequest(URI uri, Class<T> responseType) {
    try {
      HttpClient client = HttpClient.newHttpClient();
      HttpRequest request = HttpRequest.newBuilder().uri(uri)
          .header("Content-Type", "application/json").header("X-Session-ID", session.getId()).GET()
          .build();
      HttpResponse<String> response = client.send(request, BodyHandlers.ofString());
      return switch (response.statusCode()) {
        case 200 -> new ApiResponse<>(serializer.fromJson(response.body(), responseType));
        default ->
            throw new RequestError(serializer.fromJson(response.body(), ErrorResponse.class));
      };
    } catch (Exception e) {
      return new ApiResponse<>(new ApiError(e.getMessage()));
    }
  }


  private <T> ApiResponse<T> postRequest(URI uri, Object data, Class<T> responseType) {
    HttpClient client = null;
    HttpRequest request = null;
    try {
      client = HttpClient.newHttpClient();
      request = HttpRequest.newBuilder().uri(uri).header("Content-Type", "application/json")
          .POST(HttpRequest.BodyPublishers.ofString(serializer.toJson(data))).build();
      HttpResponse<String> response = client.send(request, BodyHandlers.ofString());
      return switch (response.statusCode()) {
        case 200 -> new ApiResponse<>(serializer.fromJson(response.body(), responseType));
        default ->
            throw new RequestError(serializer.fromJson(response.body(), ErrorResponse.class));
      };
    } catch (Exception e) {
      return new ApiResponse<>(new ApiError(e.getMessage()));
    }
  }
}

