package main;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

public class Person {
    public String userName;
    public long favoriteNumber;
    public String[] interests;
    @Override
    @JsonIgnore
    public String toString() {
        return "userName: " + this.userName +
                ", favoriteNumber: " + this.favoriteNumber +
                ", interests: " +  this.interests[0] + "/"+  this.interests[1];
    }
    /* Publicにプロパティを晒さなくても下記のようにアノテーション付ければ利用可能
    @JsonProperty
    private String userName;
    @JsonProperty
    private int favoriteNumber;
    @JsonProperty
    private String[] interests;

    // プロパティがPrivateの場合Jacksonが利用するためGetter/Setter必須

    @JsonIgnore
    public String getUserName() {
        return userName;
    }


    @JsonIgnore
    public void setUserName(String userName) {
        this.userName = userName;
    }


    @JsonIgnore
    public int getFavoriteNumber() {
        return favoriteNumber;
    }


    @JsonIgnore
    public void setFavoriteNumber(int favoriteNumber) {
        this.favoriteNumber = favoriteNumber;
    }

    @JsonIgnore
    public String[] getInterests() {
        return interests;
    }

    @JsonIgnore
    public void setInterests(String[] interests) {
        this.interests = interests;
    }
    @Override
    @JsonIgnore
    public String toString() {
        return "userName: " + this.userName +
                ", favoriteNumber: " + this.favoriteNumber +
                ", interests: " +  this.interests[0] + "/"+  this.interests[1];
    }
    */
    /*
    public Person() {}

    public Person(String userName, int favoriteNumber, String[] interests) {
        this.userName = userName;
        this.favoriteNumber = favoriteNumber;
        this.interests = interests;
    }
     */

}