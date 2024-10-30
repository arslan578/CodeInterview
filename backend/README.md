# Project: Asset Management System

This document provides detailed information about the performance issues identified during the development of the application and the solutions implemented to address them. The application is designed to manage assets, including associated IP addresses and ports, with features like pagination, search, and data sorting.

## Table of Contents

1. [Overview](#overview)
2. [Performance Problems and Solutions](#performance-problems-and-solutions)
    - [Problem 1: Inefficient Fetching and Data Handling](#problem-1-inefficient-fetching-and-data-handling)
    - [Problem 2: Slow Data Sorting](#problem-2-slow-data-sorting)
    - [Problem 3: Excessive Duplicate Data Rows](#problem-3-excessive-duplicate-data-rows)
    - [Problem 4: Unnecessary Data Fetching During Search](#problem-4-unnecessary-data-fetching-during-search)
    - [Problem 5: Lack of Data Caching](#problem-5-lack-of-data-caching)
    - [Problem 6: Frontend Rendering of Large Data](#problem-6-frontend-rendering-of-large-data)
    - [Problem 7: Search Performance](#problem-7-search-performance)
3. [Summary of Changes](#summary-of-changes)
4. [Getting Started](#getting-started)
5. [Running the Application](#running-the-application)

---

## Overview

This application allows users to manage assets, each associated with IPs and ports. Initially, some performance issues were identified, which affected both backend and frontend efficiency. This document outlines those issues and the solutions implemented to improve the overall performance.

---

## Performance Problems and Solutions

### Problem 1: Inefficient Fetching and Data Handling

- **Issue**: The application was initially fetching all records at once, causing high memory usage and slow load times, especially with large datasets.

- **Solution**: Pagination was implemented, limiting the number of rows fetched from the backend per request. This significantly reduced the amount of data being transferred, lowering memory usage and improving load times.

- **Changes Made**:
    - Backend: SQL queries were updated to include `LIMIT` and `OFFSET` to fetch data in chunks.
    - Frontend: Pagination logic was introduced to allow users to navigate between pages.

### Problem 2: Slow Data Sorting

- **Issue**: Sorting was initially being performed on the frontend, which caused delays, especially for large datasets. Moreover, an artificial delay was added to simulate sorting, further slowing down the process.

- **Solution**: Sorting was moved to the backend, allowing the database to handle it efficiently. The artificial delay was removed to streamline the user experience.

- **Changes Made**:
    - SQL query updated with `ORDER BY host` to handle sorting directly in the database.
    - Removed the artificial sorting delay from the frontend.

### Problem 3: Excessive Duplicate Data Rows

- **Issue**: When fetching data with associated IPs and ports, duplicate rows were being returned due to the `LEFT JOIN` operation on multiple tables. This resulted in redundant data, slowing down the application.

- **Solution**: Subqueries using `GROUP_CONCAT` were introduced to group the associated IPs and ports into single rows, reducing redundancy.

- **Changes Made**:
    - Updated SQL query to concatenate IPs and ports into single rows using `GROUP_CONCAT`.

### Problem 4: Unnecessary Data Fetching During Search

- **Issue**: Each keystroke in the search input triggered a new data fetch request, leading to multiple unnecessary network calls and slowing down the system.

- **Solution**: Debouncing was implemented to delay the data fetching until the user stopped typing for a short period, reducing unnecessary fetch requests.

- **Changes Made**:
    - Debounce functionality was added to the search input to reduce the number of requests made during typing.

### Problem 5: Lack of Data Caching

- **Issue**: Every time a user navigated between pages, the application fetched the data again from the backend, even if it had been fetched before.

- **Solution**: A simple client-side caching mechanism was introduced to store fetched data, reducing redundant network requests.

- **Changes Made**:
    - Introduced caching for each page of data on the frontend.

### Problem 6: Frontend Rendering of Large Data

- **Issue**: The frontend attempted to render all the data at once, which caused slow rendering and freezing of the user interface for large datasets.

- **Solution**: Virtual scrolling was introduced to render only the visible rows, improving performance.

- **Changes Made**:
    - Implemented virtual scrolling using libraries like `react-window` or similar.

### Problem 7: Search Performance

- **Issue**: The search functionality was handled on the frontend, resulting in performance issues when searching through large datasets.

- **Solution**: Search queries were moved to the backend, allowing the database to handle search filtering more efficiently.

- **Changes Made**:
    - SQL `LIKE` queries were added to the backend, allowing the database to filter search results before sending them to the frontend.

---

## Summary of Changes

1. **Pagination**: Reduced load by fetching only a limited number of rows at a time.
2. **Backend Sorting**: Moved sorting logic to the backend for better performance.
3. **Reduced Duplicate Rows**: Used `GROUP_CONCAT` to eliminate redundant rows.
4. **Debounced Search**: Optimized search functionality to reduce unnecessary requests.
5. **Client-Side Caching**: Stored fetched data in memory to avoid redundant requests.
6. **Virtual Scrolling**: Improved rendering performance for large datasets.
7. **Backend Search**: Handled search filtering in the backend for better efficiency.

---
