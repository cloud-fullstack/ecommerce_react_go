package com.rezaverse.backend.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class Handler {

    private final ApiService apiService;

    @Autowired
    public Handler(ApiService apiService) {
        this.apiService = apiService;
    }

    public String corsAnywhere() {
        // Implement your logic here
        return "CORS enabled";
    }

    public List<String> blogsByAvatarKey(String avatarKey) {
        // Call the ApiService method to get blogs by avatar key
        return apiService.getBlogsByAvatarKey(avatarKey);
    }

    public List<String> mostLovedRecentBlogs() {
        // Call the ApiService method to get most loved recent blogs
        return apiService.getMostLovedRecentBlogs();
    }

    public List<String> frontpageProductPreviews() {
        // Call the ApiService method to get frontpage product previews
        return apiService.getFrontpageProductPreviews();
    }

    public List<String> discountedProductsFrontpages() {
        // Call the ApiService method to get discounted products for the front page
        return apiService.getDiscountedProductsFrontpage();
    }

    public String storeDetails(String storeID) {
        // Call the ApiService method to get store details
        return apiService.getStoreDetails(storeID);
    }

    public String getAvatarProduct(String requestBody) {
        // Call the ApiService method to get avatar product
        return apiService.getAvatarProduct(requestBody);
    }

    public String getProductInventoryItems(String requestBody) {
        // Call the ApiService method to get product inventory items
        return apiService.getProductInventoryItems(requestBody);
    }

    public String genToken(String requestBody) {
        // Call the ApiService method to generate a token
        return apiService.generateToken(requestBody);
    }

    public String grabProfilePicture(String requestBody) {
        // Call the ApiService method to grab profile picture
        return apiService.grabProfilePicture(requestBody);
    }

    public String cancelOrder(String requestBody) {
        // Call the ApiService method to cancel an order
        return apiService.cancelOrder(requestBody);
    }

    public String deliverDropbox() {
        // Call the ApiService method to deliver dropbox
        return apiService.deliverDropbox();
    }

    public String completeOrder(String requestBody) {
        // Call the ApiService method to complete an order
        return apiService.completeOrder(requestBody);
    }

    public String insertDropboxRepo(String requestBody) {
        // Call the ApiService method to insert dropbox repo
        return apiService.insertDropboxRepo(requestBody);
    }

    public String insertDropbox(String requestBody) {
        // Call the ApiService method to insert dropbox
        return apiService.insertDropbox(requestBody);
    }

    public String updateDropboxContents(String requestBody) {
        // Call the ApiService method to update dropbox contents
        return apiService.updateDropboxContents(requestBody);
    }

    public String uploadPicture(String pictureType, String requestBody) {
        // Call the ApiService method to upload a picture
        return apiService.uploadPicture(pictureType, requestBody);
    }

    public String deleteAnswer(String answerID) {
        // Call the ApiService method to delete an answer
        return apiService.deleteAnswer(answerID);
    }

    public String insertAnswer(String requestBody) {
        // Call the ApiService method to insert an answer
        return apiService.insertAnswer(requestBody);
    }

    public String deleteQuestion(String questionID) {
        // Call the ApiService method to delete a question
        return apiService.deleteQuestion(questionID);
    }

    public String insertQuestion(String requestBody) {
        // Call the ApiService method to insert a question
        return apiService.insertQuestion(requestBody);
    }

    // Other handler methods can be added here...
}