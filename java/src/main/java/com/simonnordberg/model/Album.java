package com.simonnordberg.model;


import com.google.gson.annotations.SerializedName;

public class Album {

  @SerializedName("Title")
  String title;

  @Override
  public String toString() {
    return "Album{" + "title='" + title + '\'' + '}';
  }
}

