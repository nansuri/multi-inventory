# Multi Inventory Management System

## Description

Multi Inventory Management System is a mobile-web application that allows users to manage their inventory in a multi-user environment. It provides a user-friendly interface for users to add, update, and delete inventory items, as well as view the inventory history. We are planning to create the MVP version.

## Features

- User authentication and authorization
- Sidebar menu
    - Dashboard
    - Inventory management 
        - Scanning barcode / QR using integrated barcode scanner
            - Adding new inventory items (Price, Location, Halal / Non Halal Status, etc)
        - Updating existing inventory items
        - Deleting inventory items
    - Sales Management
        - Create new sales order
            - Scanning barcode / QR using integrated barcode scanner, search in the inventory whether the item exists or not
            - Inside Sales page, we can see quantity of each item and total price
            - Sales can check whether all order already fullfiled or not
        - Order Checker
            - Order checker can check whether each items of the order is fullfiled or not
        - Update existing sales order
        - Delete sales order
    - User menu
        - User profile
        - Logout

## Development Environment

- Vue JS / flutter choose the best options for fast development
- Golang using DDD
- Docker compose