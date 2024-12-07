package com.rezaverse.backend.repository;

import com.rezaverse.backend.model.LogEntry;
import org.springframework.data.jpa.repository.JpaRepository;

public interface LogRepository extends JpaRepository<LogEntry, Long> {
    // Custom query methods can be defined here
}