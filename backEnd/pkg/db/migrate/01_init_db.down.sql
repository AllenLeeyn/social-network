-- 1. Drop Triggers
DROP TRIGGER IF EXISTS update_comment_count_after_insert;
DROP TRIGGER IF EXISTS update_like_count_after_insert_in_comment_feedback;
DROP TRIGGER IF EXISTS update_like_count_after_update_in_comment_feedback;
DROP TRIGGER IF EXISTS update_like_count_after_insert_in_post_feedback;
DROP TRIGGER IF EXISTS update_like_count_after_update_in_post_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_comment_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_comment_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_insert_in_post_feedback;
DROP TRIGGER IF EXISTS update_dislike_count_after_update_in_post_feedback;

-- 2. Drop View
DROP VIEW IF EXISTS v_posts;

-- 3. Drop Tables
DROP TABLE IF EXISTS post_categories;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS post_feedback;
DROP TABLE IF EXISTS comment_feedback;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS account_type;
DROP TABLE IF EXISTS categories;
