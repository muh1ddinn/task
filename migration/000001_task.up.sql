CREATE TABLE categories (
  id UUID PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at INTEGER DEFAULT 0
);

CREATE TABLE contact (
  id UUID PRIMARY KEY,
  phone VARCHAR(20)  NOT NULL,
  name VARCHAR(50) NOT NULL,
  email VARCHAR(50) UNIQUE NOT NULL,
  address VARCHAR(50),
  category VARCHAR(50)  NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at INTEGER DEFAULT 0
);


CREATE TABLE contact_history (
    id UUID PRIMARY KEY,
    contact_id UUID NOT NULL,
    phone VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    address VARCHAR(50),
    category VARCHAR(50),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    change_type VARCHAR(10) NOT NULL,
    FOREIGN KEY (contact_id) REFERENCES contact(id)
);
