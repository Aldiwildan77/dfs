-- author: aldiwild77@gmail.com
CREATE TABLE IF NOT EXISTS file_metadata (
  id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'Unique identifier',
  guild_id VARCHAR(50) NOT NULL COMMENT 'Discord guild ID',
  channel_id VARCHAR(50) NOT NULL COMMENT 'Discord channel ID',
  message_id VARCHAR(50) NOT NULL COMMENT 'Discord message ID',
  url TEXT NOT NULL COMMENT 'URL of the file',
  original_file_url TEXT GENERATED ALWAYS AS (SUBSTRING_INDEX(url, '?', 1)) VIRTUAL,
  filename VARCHAR(255) NOT NULL COMMENT 'Original file name',
  hash TEXT NOT NULL COMMENT 'Current hash file',
  expired_at TIMESTAMP NOT NULL COMMENT 'Timestamp of expiration',
  issued_at TIMESTAMP NOT NULL COMMENT 'Timestamp of issuance',
  last_rotated_at TIMESTAMP NOT NULL COMMENT 'Timestamp of last rotation',
  last_accessed TIMESTAMP NOT NULL COMMENT 'Timestamp of last access',
  access_count INT DEFAULT 0 COMMENT 'Number of times the file has been accessed',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp'
);

CREATE INDEX idx_message_id ON file_metadata (message_id);

CREATE INDEX idx_last_rotated_at ON file_metadata (last_rotated_at);

CREATE UNIQUE INDEX uk_comp_guild_channel_message ON file_metadata (guild_id, channel_id, message_id);