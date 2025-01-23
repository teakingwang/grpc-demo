-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,  -- 用户ID,自增主键
    username VARCHAR(50) NOT NULL,         -- 用户名,不允许为空
    email VARCHAR(100) NOT NULL,           -- 电子邮箱,不允许为空
    password VARCHAR(100) NOT NULL,        -- 密码哈希,不允许为空
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间,默认当前时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- 更新时间,自动更新
);

-- 创建用户名唯一索引
CREATE UNIQUE INDEX idx_username ON users(username);
-- 创建邮箱唯一索引
CREATE UNIQUE INDEX idx_email ON users(email);