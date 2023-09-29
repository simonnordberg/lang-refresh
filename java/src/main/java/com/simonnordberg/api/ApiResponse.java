package com.simonnordberg.api;

public class ApiResponse<T> {
    private T data;
    private ApiError error;

    public ApiResponse(T data) {
        this.data = data;
    }

    public ApiResponse(ApiError error) {
        this.error = error;
    }

    public T getData() {
        return data;
    }

    public ApiError getError() {
        return error;
    }

    public boolean hasError() {
        return error != null;
    }

    @Override
    public String toString() {
        return "ApiResponse{" +
                "data=" + data +
                ", error=" + error +
                '}';
    }
}
