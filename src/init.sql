-- Create article table if it does not exist
CREATE TABLE IF NOT EXISTS articles (
  id INT AUTO_INCREMENT PRIMARY KEY,
  link VARCHAR(255) UNIQUE,
  title VARCHAR(255),
  processed BOOLEAN DEFAULT FALSE
);

-- Create relations table if it does not exist
CREATE TABLE IF NOT EXISTS relations (
  parent INT,
  child INT,
  UNIQUE KEY uniq_relation (parent, child),
  FOREIGN KEY (parent) REFERENCES articles(id),
  FOREIGN KEY (child) REFERENCES articles(id)
);

-- Create index on url column if it does not exist
ALTER TABLE articles ADD UNIQUE INDEX IF NOT EXISTS link_index (link);

-- Insert the origin article (assuming DUAL table doesn't exist)
INSERT INTO articles (link, title)
  SELECT 'Marseille', 'Marseille'
  WHERE NOT EXISTS (SELECT 1 FROM articles WHERE link = 'Marseille');
