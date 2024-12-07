package com.rezaverse.backend.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import com.rezaverse.backend.service.Handler;
import com.rezaverse.backend.service.ApiService;


import java.util.List;

@RestController
@RequestMapping("/api")
public class ApiController {

    private final Handler handler;
    private final ApiService databaseService;

    @Autowired
    public ApiController(Handler handler, ApiService databaseService) {
        this.handler = handler;
        this.databaseService = databaseService;
    }

    // API Routes
    @GetMapping("/cors")
    public ResponseEntity<String> corsAnywhere() {
        return ResponseEntity.ok(handler.corsAnywhere());
    }

    @GetMapping("/avatar-blogs/{avatarKey}")
    public ResponseEntity<List<String>> blogsByAvatarKey(@PathVariable String avatarKey) {
        return ResponseEntity.ok(handler.blogsByAvatarKey(avatarKey));
    }

    @GetMapping("/most-loved-recent-blogs/")
    public ResponseEntity<List<String>> mostLovedRecentBlogs() {
        return ResponseEntity.ok(handler.mostLovedRecentBlogs());
    }

    @GetMapping("/frontpage-product-previews/")
    public ResponseEntity<List<String>> frontpageProductPreviews() {
        return ResponseEntity.ok(handler.frontpageProductPreviews());
    }

    @GetMapping("/discounted-products-frontpage/")
    public ResponseEntity<List<String>> discountedProductsFrontpages() {
        return ResponseEntity.ok(handler.discountedProductsFrontpages());
    }

    @GetMapping("/store-details/{storeID}")
    public ResponseEntity<String> storeDetails(@PathVariable String storeID) {
        return ResponseEntity.ok(handler.storeDetails(storeID));
    }

    @PostMapping("/get-avatar-product/")
    public ResponseEntity<String> getAvatarProduct(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.getAvatarProduct(requestBody));
    }

    @PostMapping("/get-product-inventory-items/")
    public ResponseEntity<String> getProductInventoryItems(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.getProductInventoryItems(requestBody));
    }

    @PostMapping("/gen-token/")
    public ResponseEntity<String> genToken(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.genToken(requestBody));
    }

    @PostMapping("/profile-picture/")
    public ResponseEntity<String> grabProfilePicture(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.grabProfilePicture(requestBody));
    }

    // SL Only Routes
    @PostMapping("/cancel-order-sl/")
    public ResponseEntity<String> cancelOrder(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.cancelOrder(requestBody));
    }

    @GetMapping("/deliver-dropbox/")
    public ResponseEntity<String> deliverDropbox() {
        return ResponseEntity.ok(handler.deliverDropbox());
    }

    @PostMapping("/complete-order/")
    public ResponseEntity<String> completeOrder(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.completeOrder(requestBody));
    }

    @PostMapping("/insert-dropbox-repo/")
    public ResponseEntity<String> insertDropboxRepo(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.insertDropboxRepo(requestBody));
    }

    @PostMapping("/insert-dropbox/")
    public ResponseEntity<String> insertDropbox(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.insertDropbox(requestBody));
    }

    @PostMapping("/update-dropbox-contents/")
    public ResponseEntity<String> updateDropboxContents(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.updateDropboxContents(requestBody));
    }

    // Authenticated API Routes
    @PostMapping("/upload-picture/{pictureType}")
    public ResponseEntity<String> uploadPicture(@PathVariable String pictureType, @RequestBody String requestBody) {
        return ResponseEntity.ok(handler.uploadPicture(pictureType, requestBody));
    }

    @DeleteMapping("/answer/{answerID}")
    public ResponseEntity<String> deleteAnswer(@PathVariable String answerID) {
        return ResponseEntity.ok(handler.deleteAnswer(answerID));
    }

    @PostMapping("/answer/")
    public ResponseEntity<String> insertAnswer(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.insertAnswer(requestBody));
    }

    @DeleteMapping("/question/{questionID}")
    public ResponseEntity<String> deleteQuestion(@PathVariable String questionID) {
        return ResponseEntity.ok(handler.deleteQuestion(questionID));
    }

    @PostMapping("/question/")
    public ResponseEntity<String> insertQuestion(@RequestBody String requestBody) {
        return ResponseEntity.ok(handler.insertQuestion(requestBody));
    }

    // Other routes...
}