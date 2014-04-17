#!/bin/bash

#create database yunyun
#  CREATE DATABASE yunyun;
#  USE yunyun;

#setup user table
#user | password | id(k) 
#  CREATE TABLE users (
#    id INT NOT NULL AUTO_INCREMENT, 
#    user VARCHAR(100) NOT NULL, 
#    pw_hash VARCHAR(128) NOT NULL,
#    PRIMARY KEY(id)
#  );

# add user: INSERT INTO user(user, pw_hash) VALUES("user", "pwhash");

#setup words table
#id | userId | kotoba | imi | level | nextReview
#  CREATE TABLE kotoba (
#    id INT NOT NULL AUTO_INCREMENT, 
#    user_id INT NOT NULL, 
#    kotoba VARCHAR(128) NOT NULL,
#    imi VARCHAR(256) NOT NULL, 
#    level INT NOT NULL, 
#    next_review TIMESTAMP NOT NULL,
#    PRIMARY KEY(id)
#  );
#
# add kotoba: INSERT INTO kotoba(user_id, kotoba, imi, level, next_review) 
#             VALUES("user_id", "kotoba", "imi", 3, "date");
